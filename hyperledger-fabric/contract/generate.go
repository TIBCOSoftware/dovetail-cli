/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

// Package contract implements generate and deploy chaincode for hyperledger fabric
package contract

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	//"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	//"path/filepath"
	"runtime"
	"strings"

	//"github.com/TIBCOSoftware/dovetail-cli/files"
	//"github.com/TIBCOSoftware/dovetail-cli/languages"
	"github.com/TIBCOSoftware/dovetail-cli/model"
	"github.com/TIBCOSoftware/dovetail-cli/pkg/contract"
	wgutil "github.com/TIBCOSoftware/dovetail-cli/util"
	"github.com/project-flogo/cli/api"
	"github.com/project-flogo/core/app"
	"github.com/project-flogo/flow/definition"
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

	appConfig, err := model.ParseApp(d.Opts.ModelFile)
	if err != nil {
		return err
	}

	// Create project
	appProject, err := api.CreateProject(d.Opts.TargetDir, appConfig.Name, d.Opts.ModelFile, "")
	if err != nil {
		return err
	}
	/*goProject := languages.NewGo(d.Opts.TargetDir, appConfig.Name)

	err = goProject.Init()
	if err != nil {
		return err
	}

	defer goProject.Cleanup()

	appDir := goProject.GetAppDir()

	activities, triggers, err := getAppResources(appConfig)
	if err != nil {
		return err
	}

	err = wgutil.CopyFile(d.Opts.ModelFile, filepath.Join(appDir, fmt.Sprintf("%s.%s", appConfig.Name, "json")))
	if err != nil {
		return err
	}

	err = createShimSupportFile(appDir, appConfig.Name, d.Opts, activities, triggers)
	if err != nil {
		return err
	}

	err = createShimFile(appDir, appConfig.Name)
	if err != nil {
		return err
	}

	err = createFunctionImportFile(appConfig.Name, appDir, d.Opts.ModelFile)
	if err != nil {
		return err
	}

	err = vendorFiles(path.Join(goProject.GetTargetDir(), strings.ToLower(appConfig.Name)), appDir)
	if err != nil {
		return err
	}

	err = createResourceBundle(appConfig.Name, appDir, activities, triggers)
	if err != nil {
		return err
	}

	// If it is file compress
	if goProject.IsFile() {
		logger.Println("Compressing files...")
		err = files.ZipFolder(goProject.GetInputTargetDir(), goProject.GetTargetDir())
		if err != nil {
			return err
		}
	}
	*/

	logger.Println("Generating Hyperledger Fabric smart contract... Done")
	logger.Printf("Generated artifacts: '%s'\n", appProject.Dir())
	return nil
}

func vendorFiles(target, srcdir string) error {
	separator := ":"
	if runtime.GOOS == "windows" {
		separator = ";"
	}
	err := os.Setenv("GOPATH", strings.Join([]string{os.Getenv("GOPATH"), target}, separator))
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

/*func createFunctionImportFile(ccName, targetdir, modelFile string) error {
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
}*/
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
