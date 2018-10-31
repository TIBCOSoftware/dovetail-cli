/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package contract

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/TIBCOSoftware/dovetail-cli/pkg/contract"
	wgutil "github.com/TIBCOSoftware/dovetail-cli/util"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-lib/app"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

// Generator defines the generator attributes
type Generator struct {
	Opts *GenOptions
}

// GenOptions defines the generator options
type GenOptions struct {
	TargetDir         string
	ModelFile         string
	Version           string
	EnableTxnSecurity bool
}

type ResourceRef struct {
	Name string
	Ref  string
	ID   string
}

type TemplateData struct {
	CCName            string
	ActivityRefs      map[string]ResourceRef
	TriggerRefs       map[string]ResourceRef
	Functions         map[string]string
	EnableTxnSecurity bool
}

// NewGenerator is the generator constructor
func NewGenerator(opts *GenOptions) contract.Generator {
	return &Generator{Opts: opts}
}

// NewGenOptions is the options constructor
func NewGenOptions(targetPath, modelFile, version string, enableSecurity bool) *GenOptions {

	return &GenOptions{TargetDir: targetPath, ModelFile: modelFile, Version: version, EnableTxnSecurity: enableSecurity}
}

// Generate generates a smart contract for the given options
func (d *Generator) Generate() error {
	logger.Println("Generating Hyperledger Fabric smart contract...")

	appConfig, err := parseApp(d.Opts.ModelFile)
	if err != nil {
		return err
	}

	target := wgutil.CreateTargetDirs(path.Join(d.Opts.TargetDir, "hlf", "src", strings.ToLower(appConfig.Name)))

	activities, triggers, err := getAppResources(appConfig)
	if err != nil {
		return err
	}

	if len(triggers) == 0 || len(triggers) > 1 {
		return fmt.Errorf("There must be one and only one trigger defined in smart contract application")
	}

	err = wgutil.CopyFile(d.Opts.ModelFile, target+"/"+appConfig.Name+".json")
	if err != nil {
		return err
	}

	err = createShimSupportFile(target, appConfig.Name, d.Opts, activities, triggers)
	if err != nil {
		return err
	}

	err = createShimFile(target, appConfig.Name)
	if err != nil {
		return err
	}

	err = createFunctionImportFile(appConfig.Name, target, d.Opts.ModelFile)
	if err != nil {
		return err
	}

	err = vendorFiles(path.Join(d.Opts.TargetDir, "hlf"), target)
	if err != nil {
		return err
	}

	err = createResourceBundle(appConfig.Name, target, activities, triggers)
	if err != nil {
		return err
	}

	logger.Println("Generating Hyperledger Fabric smart contract... Done")
	return nil
}

func vendorFiles(target, srcdir string) error {
	separator := ":"
	if runtime.GOOS == "windows" {
		separator = ";"
	}
	err := os.Setenv("GOPATH", os.Getenv("GOPATH")+separator+target)
	if err != nil {
		return fmt.Errorf("error set up GOPATH:%v", err)
	}

	wdir, _ := os.Getwd()
	err = os.Chdir(srcdir)
	if err != nil {
		return fmt.Errorf("vendorFiles command 'cd' err %v", err)
	}
	cdir, err := os.Getwd()
	fmt.Printf("current dir=%s, gopath=%s\n", cdir, os.Getenv("GOPATH"))

	args := []string{"init"}
	cmd := exec.Command("govendor", args...)
	logger.Printf("govendor command %v\n", cmd.Args)
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("vendorFiles command 'govendor init' err %v, output=%s\n", err, string(out))
	}

	args = []string{"add", "+external"}
	cmd = exec.Command("govendor", args...)
	logger.Printf("govendor command %v\n", cmd.Args)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("vendorFiles govendor add +external err %v", err)
	}

	os.Chdir(wdir)
	return nil
}

func createFunctionImportFile(ccName, targetdir, modelFile string) error {
	content, err := ioutil.ReadFile(modelFile)
	if err != nil {
		return err
	}

	strcontent := string(content)
	usedFuncs := TemplateData{CCName: ccName, Functions: make(map[string]string)}
	for nm, ref := range functions {
		if strings.Contains(strcontent, nm) {
			usedFuncs.Functions[nm] = ref
		}
	}

	if len(usedFuncs.Functions) == 0 {
		return nil
	}

	gofile := "import.go"
	logger.Printf("Create %s file ....\n", gofile)

	f, error := os.Create(path.Join(targetdir, gofile))
	if error != nil {
		return error
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	tmpl, err := Asset("resources/" + "import.template")
	if err != nil {
		return err
	}

	t, err := template.New("import").Parse(string(tmpl))
	if err != nil {
		return fmt.Errorf("error processing import.template file, err %v", err)
	}

	err = t.Execute(writer, usedFuncs)
	if err != nil {
		return err
	}
	writer.Flush()

	return nil
}
func createResourceBundle(ccName, targetdir string, activities, triggers map[string]ResourceRef) error {
	files := make([]string, 0)
	files = append(files, path.Join(targetdir, ccName+".json"))

	for _, ref := range activities {
		files = append(files, path.Join(targetdir, "vendor", ref.Ref, "activity.json"))
	}
	for _, ref := range triggers {
		files = append(files, path.Join(targetdir, "vendor", ref.Ref, "trigger.json"))
	}
	args := []string{"-prefix", targetdir, "-pkg", "main", "-o", path.Join(targetdir, "resources.go")}
	args = append(args, files...)
	cmd := exec.Command("go-bindata", args...)
	logger.Printf("go-bindata command %v\n", cmd.Args)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("createResourceBundle err %v", err)
	}
	return nil
}
func createShimSupportFile(targetdir, ccName string, opts *GenOptions, activities, triggers map[string]ResourceRef) error {
	gofile := "shim_support.go"
	logger.Printf("Create shim support %s file ....\n", gofile)

	f, error := os.Create(path.Join(targetdir, gofile))
	if error != nil {
		return error
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	tmpl, err := Asset("resources/" + "shim_support.template")
	if err != nil {
		return err
	}

	t, err := template.New("shim_support").Parse(string(tmpl))
	if err != nil {
		return fmt.Errorf("Error processing shim_support.template file, err %v", err)
	}

	data := TemplateData{CCName: ccName, ActivityRefs: activities, TriggerRefs: triggers, EnableTxnSecurity: opts.EnableTxnSecurity}
	err = t.Execute(writer, data)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func createShimFile(targetdir, ccName string) error {
	gofile := "shim.go"
	logger.Printf("Create shim %s file ....\n", gofile)

	f, error := os.Create(path.Join(targetdir, gofile))
	if error != nil {
		return error
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	tmpl, err := Asset("resources/" + "shim.template")
	if err != nil {
		return err
	}

	t, err := template.New("shim").Parse(string(tmpl))
	if err != nil {
		return fmt.Errorf("Error processing shim.template file, err %v", err)
	}

	data := TemplateData{CCName: ccName}
	err = t.Execute(writer, data)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func copyResourceFiles(rslibpath, target string, activities, triggers map[string]ResourceRef) error {
	acdir := wgutil.CreateDirIfNotExist(target, "activity")
	for _, ref := range activities {
		rsdir := wgutil.CreateDirIfNotExist(acdir, ref.Name)
		err := wgutil.CopyFile(rslibpath+"/"+ref.Ref+"/activity.json", rsdir+"/activity.json")
		if err != nil {
			return err
		}
	}

	trigdir := wgutil.CreateDirIfNotExist(target, "trigger")
	for _, ref := range triggers {
		rsdir := wgutil.CreateDirIfNotExist(trigdir, ref.Name)
		err := wgutil.CopyFile(rslibpath+"/"+ref.Ref+"/trigger.json", rsdir+"/trigger.json")
		if err != nil {
			return err
		}
	}
	return nil
}

func parseApp(modelfile string) (*app.Config, error) {
	appCfg := &app.Config{}

	flowjson, err := ioutil.ReadFile(modelfile)
	if err != nil {
		return appCfg, err
	}

	jsonParser := json.NewDecoder(bytes.NewReader(flowjson))
	err = jsonParser.Decode(&appCfg)
	if err != nil {
		return nil, err
	}

	return appCfg, nil
}

func getAppResources(appConfig *app.Config) (activities, triggers map[string]ResourceRef, err error) {
	activities = make(map[string]ResourceRef)
	triggers = make(map[string]ResourceRef)

	for _, tConfig := range appConfig.Triggers {
		tokens := strings.Split(tConfig.Ref, "/")
		triggers[tConfig.Ref] = ResourceRef{Name: tokens[len(tokens)-1], Ref: tConfig.Ref, ID: tConfig.Id}
	}

	rConfigs := appConfig.Resources
	for _, rConfig := range rConfigs {
		var defRep *definition.DefinitionRep
		err := json.Unmarshal(rConfig.Data, &defRep)
		if err != nil {
			return triggers, nil, err
		}

		for _, t := range defRep.Tasks {
			tokens := strings.Split(t.ActivityCfgRep.Ref, "/")
			activities[t.ActivityCfgRep.Ref] = ResourceRef{Name: tokens[len(tokens)-1], Ref: t.ActivityCfgRep.Ref}
		}
	}
	return activities, triggers, nil
}

var functions map[string]string

func init() {
	//TODO: how to support built-in and user defined functions
	/*
		functions = make(map[string]string)
		functions["array.contains"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/array/contains"
		functions["array.create"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/array/create"
		functions["array.forEach"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/array/forEach"
		functions["boolean.false"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/boolean/false"
		functions["boolean.not"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/boolean/not"
		functions["boolean.true"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/boolean/true"
		functions["datetime.currentDate"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/datetime/currentDate"
		functions["datetime.currentDatetime"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/datetime/currentDatetime"
		functions["datetime.currentTime"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/datetime/currentTime"
		functions["datetime.formatDate"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/datetime/formatDate"
		functions["datetime.formatDatetime"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/datetime/formatDatetime"
		functions["datetime.formatTime"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/datetime/formatTime"
		functions["datetime.now"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/datetime/now"
		functions["number.int64"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/number/int64"
		functions["number.len"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/number/len"
		functions["string.base64ToString"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/base64ToString"
		functions["string.concat"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/concat"
		functions["string.contains"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/contains"
		functions["string.dateFormat"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/dateFormat"
		functions["string.datetimeFormat"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/datetimeFormat"
		functions["string.endsWith"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/endsWith"
		functions["string.equals"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/equals"
		functions["string.equalsignorecase"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/equalsignorecase"
		functions["string.index"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/index"
		functions["string.lastIndex"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/lastIndex"
		functions["string.length"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/length"
		functions["string.lowerCase"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/lowerCase"
		functions["string.regex"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/regex"
		functions["string.split"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/split"
		functions["string.startsWith"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/startsWith"
		functions["string.stringToBase64"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/stringToBase64"
		functions["string.substring"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/substring"
		functions["string.substringAfter"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/substringAfter"
		functions["string.substringBefore"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/substringBefore"
		functions["string.timeFormat"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/timeFormat"
		functions["string.tostring"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/tostring"
		functions["string.trim"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/trim"
		functions["string.upperCase"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/upperCase"
		functions["string.renderJSON"] = "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/utility/renderJSON"
	*/
}
