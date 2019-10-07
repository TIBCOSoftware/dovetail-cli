/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package contract

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"

	"github.com/TIBCOSoftware/dovetail-cli/files"
	"github.com/TIBCOSoftware/dovetail-cli/languages"
	"github.com/TIBCOSoftware/dovetail-cli/model"
	"github.com/TIBCOSoftware/dovetail-cli/pkg/contract"
	wgutil "github.com/TIBCOSoftware/dovetail-cli/util"
	"github.com/project-flogo/core/app"
	"github.com/project-flogo/core/trigger"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

type Generator struct {
	Opts *Options
}

type Options struct {
	ModelFile string
	Version   string
	State     string
	Commands  []string
	TargetDir string
	Namespace string
	Pom       string
}

type DataState struct {
	NS                string
	Class             string
	ContractClass     string
	CordaClass        string
	Attributes        []model.ResourceAttribute
	Parent            string
	Participants      []string
	ExitKeys          []string
	IsSchedulable     bool
	ScheduledActivity string
}

type ContractData struct {
	ContractClass string
	NS            string
	Flow          string
	Commands      []Command
	States        []DataState
	AppClass      string
}

type AppData struct {
	NS      string
	AppName string
}
type Command struct {
	Name       string
	NS         string
	Attributes []model.ResourceAttribute
}

type ModelResources struct {
	AppName      string
	Assets       []string
	Transactions []string
	Schemas      map[string]string
}

var models map[string]*model.ResourceMetadataModel

// NewGenerator is the generator constructor
func NewGenerator(opts *Options) contract.Generator {
	return &Generator{Opts: opts}
}

// NewOptions is the options constructor
func NewOptions(flowModel string, version string, state string, commands []string, target, ns string, pom string) *Options {
	return &Options{ModelFile: flowModel, Version: version, State: state, Commands: commands, TargetDir: target, Namespace: ns, Pom: pom}
}

// Generate generates a smart contract for the given options
func (g *Generator) Generate() error {
	logger.Println("Generating artifacts for Corda...")
	appCfg, err := model.ParseApp(g.Opts.ModelFile)
	if err != nil {
		return err
	}

	javaProject := languages.NewJava(g.Opts.TargetDir, appCfg.Name)

	err = javaProject.Init()
	if err != nil {
		return err
	}

	defer javaProject.Cleanup()

	resources, schedulables, isComposer, err := parseFlowApp(appCfg)
	if err != nil {
		return fmt.Errorf("error parsing flow app json file, err %v", err)
	}

	appdata := AppData{AppName: appCfg.Name}
	if g.Opts.Namespace != "" {
		appdata.NS = g.Opts.Namespace
	} else {
		appdata.NS, _ = splitNamespace(resources[0].Transactions[0], "")
	}

	//create directories
	resourcedir := wgutil.CreateDirIfNotExist(javaProject.GetAppDir(), "src/main/resources", strings.Replace(appdata.NS, ".", "/", -1))
	kotlindir := wgutil.CreateDirIfNotExist(javaProject.GetAppDir(), "src/main/kotlin")
	wgutil.CreateDirIfNotExist(javaProject.GetAppDir(), "target/kotlin/classes")

	//create ContractImpl (flows)
	err = createKotlinFile(kotlindir, appdata.NS, appdata, "kotlin.contractimpl.template", fmt.Sprintf("%s%s", appdata.AppName, "Impl.kt"))
	if err != nil {
		return fmt.Errorf("createContractImplFile kotlin.contractimpl.template err %v", err)
	}

	//create Resource
	err = createResourceFiles(resourcedir, g.Opts)
	if err != nil {
		return fmt.Errorf("createResourceFiles err %v", err)
	}

	if isComposer {
		err = createComposerContractFiles(kotlindir, g.Opts, resources[0], schedulables, appdata.NS+"."+appdata.AppName)
		if err != nil {
			return err
		}
	} else {
		err = createDovetailContractFiles(kotlindir, g.Opts, resources, schedulables, appdata.NS+"."+appdata.AppName)
		if err != nil {
			return err
		}
	}

	err = compileAndJar(javaProject.GetAppDir(), appdata.NS, appdata.AppName, g.Opts.Version, "kotlin.pom.xml", g.Opts.Pom)
	if err != nil {
		return fmt.Errorf("compileAndJar kotlin.pom.xml err %v", err)
	}

	//Cleanup
	err = os.RemoveAll(path.Join(javaProject.GetAppDir(), "generated-sources"))
	if err != nil {
		return err
	}
	err = os.RemoveAll(path.Join(javaProject.GetAppDir(), "maven-archiver"))
	if err != nil {
		return err
	}
	err = os.RemoveAll(path.Join(javaProject.GetAppDir(), "maven-status"))
	if err != nil {
		return err
	}
	err = os.RemoveAll(path.Join(javaProject.GetAppDir(), "target"))
	if err != nil {
		return err
	}

	// If it is file compress
	if javaProject.IsFile() {
		logger.Println("Compressing files...")
		err = files.ZipFolder(javaProject.GetInputTargetDir(), javaProject.GetTargetDir())
		if err != nil {
			return err
		}
	}

	logger.Printf("Finished generating artifacts for Corda")
	return nil
}

func compileAndJar(targetdir, ns, clazz, version string, pomf string, dependencyFile string) error {
	logger.Printf("Compile smart contract artifacts")
	pom, err := Asset("resources/" + pomf)
	if err != nil {
		return err
	}

	dep := ""

	if dependencyFile != "" {
		deppom, err := ioutil.ReadFile(dependencyFile)
		if err != nil {
			return err
		}

		dep = string(deppom)
	}
	newpom := strings.Replace(string(pom), "%%external%%", dep, 1)

	err = wgutil.CopyContent([]byte(newpom), path.Join(targetdir, pomf))
	if err != nil {
		return err
	}

	err = wgutil.MvnPackage(ns, clazz, version, pomf, targetdir)
	if err != nil {
		return err
	}

	err = wgutil.MvnInstall(ns, clazz, version, fmt.Sprintf("%s/kotlin-%s-%s.jar", targetdir, clazz, version))
	if err != nil {
		return err
	}
	return nil
}
func createConceptKotlinFiles(dir, template string, concepts []DataState) error {
	for _, data := range concepts {
		err := createKotlinFile(dir, data.NS, data, template, data.Class+".kt")
		if err != nil {
			return fmt.Errorf("Error creating java file %s, error %v", data.Class, err)
		}
	}
	return nil
}

func createResourceFiles(resourcedir string, opts *Options) error {
	logger.Println("Copy resource file - transactions.json ...")
	err := wgutil.CopyFile(opts.ModelFile, path.Join(resourcedir, "transactions.json"))
	if err != nil {
		return fmt.Errorf("Error creating transaction.json file, error %v", err)
	}

	return nil
}

func createKotlinFile(dir, ns string, data interface{}, templateFile string, fileName string) error {
	logger.Printf("Create kotlin file %s with template %s....", fileName, templateFile)

	javadir := wgutil.CreateDirIfNotExist(dir, strings.Replace(ns, ".", "/", -1))
	f, error := os.Create(path.Join(javadir, fileName))
	if error != nil {
		return error
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	funcMap := template.FuncMap{
		"GetKotlinType":        GetKotlinType,
		"GetKotlinTypeNoArray": GetKotlinTypeNoArray,
		"GetFuncName":          GetFuncName,
		"ToKotlinString":       ToKotlinString,
		"GetParticipants":      GetParticipants,
		"GetFlowSha256":        GetFlowSha256,
	}

	tmpl, err := Asset("resources/" + templateFile)
	if err != nil {
		return err
	}

	t, err := template.New(templateFile).Funcs(funcMap).Parse(string(tmpl))
	if err != nil {
		return fmt.Errorf("Template %s err %v", templateFile, err)
	}

	err = t.Execute(writer, data)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func parseAllResources(schemas map[string]string) map[string]*model.ResourceMetadataModel {
	models := make(map[string]*model.ResourceMetadataModel)
	for r, s := range schemas {
		//get the description field for metadata
		desc := struct {
			Description string `json: "description"`
		}{}
		json.Unmarshal([]byte(s), &desc)
		models[r] = model.ParseResourceModel(desc.Description)
	}

	return models
}

func splitNamespace(assetName, targetns string) (ns, clazz string) {
	idx := strings.LastIndex(assetName, ".")
	ns = ""
	clazz = assetName

	if idx > 0 {
		ns = assetName[0:idx]
		clazz = assetName[idx+1:]
	} else {
		ns = targetns
		clazz = assetName
	}
	return
}

func ToKotlinString(attr model.ResourceAttribute) string {
	code := `"\"` + attr.Name + `\":" + `
	if attr.IsArray {
		code = code + `"[" + `
	}

	datatype := GetKotlinType(attr)
	switch datatype {
	case "net.corda.core.identity.Party":
		code = code + `"\"" +` + attr.Name + `!!.toString() + "\""`
		break
	case "List<net.corda.core.identity.Party>":
		code = code + attr.Name + `!!.stream().map(p -> p.toString()).collect(java.util.stream.Collectors.joining(","))`
		break
	case "net.corda.core.contracts.Amount<Currency>":
		code = code + `"{\"quantity\":" + ` + attr.Name + `!!.quantity + ", \"currency\":\"" + ` + attr.Name + `!!.token.currencyCode + "\"}"`
		break
	case "net.corda.core.contracts.UniqueIdentifier":
		code = code + `"\"" + linearId.toString() + "\""`
	case "String":
		code = code + attr.Name
		break
	case "List<String>":
		code = code + attr.Name + `!!.stream().collect(java.util.stream.Collectors.joining(","))`
		break
	case "Boolean":
		code = code + attr.Name
		break
	case "List<Boolean>":
		code = code + attr.Name + `!!.stream().map(v -> v.booleanValue()).collect(java.util.stream.Collectors.joining(","))`
		break
	case "java.math.BigDecimal":
		code = code + `"\"" +` + attr.Name + `!!.toPlainString()+ "\""`
		break
	case "List<java.math.BigDecimal>":
		code = code + attr.Name + `!!.stream().map(v ->v.toPlainString()).collect(java.util.stream.Collectors.joining(","))`
		break
	case "Int":
		code = code + attr.Name
		break
	case "java.util.List<Int>":
		code = code + attr.Name + `!!.stream().map(v -> v.intValue()).collect(java.util.stream.Collectors.joining(","))`
		break
	case "Long":
		code = code + attr.Name
		break
	case "List<Long>":
		code = code + attr.Name + `!!.stream().map(v -> v.longValue()).collect(java.util.stream.Collectors.joining(","))`
		break
	default:
		if attr.IsArray {
			code = code + attr.Name + `!!.stream().map{v -> v.toString()}.collect(java.util.stream.Collectors.joining(","))`
		} else {
			code = code + attr.Name + `!!.toString()`
		}
	}

	if attr.IsArray {
		code = code + ` + "]"`
	}
	return code
}

func GetKotlinType(attr model.ResourceAttribute) string {
	datatype := attr.Type

	if strings.Compare(attr.Name, "linearId") == 0 {
		datatype = "net.corda.core.contracts.UniqueIdentifier"
	} else {
		datatype = GetKotlinTypeNoArray(attr)
		if attr.IsArray {
			datatype = "List<" + datatype + ">"
		}
	}

	return datatype
}

func GetKotlinTypeNoArray(attr model.ResourceAttribute) string {
	datatype := attr.Type

	if strings.Compare(attr.Name, "linearId") == 0 {
		datatype = "net.corda.core.contracts.UniqueIdentifier"
	} else {
		switch datatype {
		case "Integer":
			datatype = "Int"
			break
		case "Double":
		case "Decimal":
			datatype = "java.math.BigDecimal"
			break
		case "DateTime":
		case "com.tibco.dovetail.system.Instant":
			datatype = "java.time.Instant"
			break
		case "com.tibco.dovetail.system.LocalDate":
			datatype = "java.time.LocalDate"
			break
		case "com.tibco.dovetail.system.Party":
			datatype = "AbstractParty"
			break
		case "org.hyperledger.composer.system.Participant":
			datatype = "AbstractParty"
			break
		case "com.tibco.dovetail.system.Cash":
			datatype = "net.corda.finance.contracts.asset.Cash.State"
			break
		case "com.tibco.dovetail.system.Amount":
		case "com.tibco.dovetail.system.Amount<Currency>":
			datatype = "net.corda.core.contracts.Amount<Currency>"
			break
		case "com.tibco.dovetail.system.UniqueIdentifier":
			datatype = "net.corda.core.contracts.UniqueIdentifier"
			break
		}
	}

	return datatype
}
func GetFuncName(name string) string {
	runes := []rune(name)
	return strings.ToUpper(string(runes[0])) + string(runes[1:])
}

func GetParticipants(state DataState) string {
	participants := make([]string, 0)
	if len(state.Participants) > 0 {
		//declarative participants $tx.path.to.Party
		for _, p := range state.Participants {
			if strings.HasPrefix(p, "$tx.") {
				tokens := strings.Split(p, ".")
				attr, ok := getAttr(tokens[0], state.Attributes)
				if !ok {
					panic(fmt.Sprintf("error processing participants %s.%s, %s is not found\n", "$tx.", p, tokens[0]))
				}
				if len(tokens) == 1 {
					participants = append(participants, toParticipant(attr.Name, attr))
				} else {
					rangevar := attr.Name
					for idx := 1; idx < len(tokens); idx++ {
						if attr.IsOptional {
							participants = append(participants, "	if ("+rangevar+" != null ) {")
						}
						if attr.IsArray {
							participants = append(participants, "	for ( v"+strconv.Itoa(idx)+" in "+rangevar+"!!) {")
							rangevar = "v" + strconv.Itoa(idx) + "." + tokens[idx]
						} else {
							rangevar = rangevar + "." + tokens[idx]
						}

						attr, ok = getAttr(tokens[idx], models[attr.Type].Attributes)
						if !ok {
							panic(fmt.Sprintf("error processing participants %s.%s, %s is not found\n", "$tx.", p, tokens[idx]))
						}
					}

					participants = append(participants, "	"+toParticipant(rangevar, attr))
					attr, _ := getAttr(tokens[0], state.Attributes)
					for idx := 1; idx < len(tokens); idx++ {
						if attr.IsArray {
							participants = append(participants, "	}")
						}
						if attr.IsOptional {
							participants = append(participants, "	}")
						}
						attr, _ = getAttr(tokens[idx], models[attr.Type].Attributes)
					}
				}
			} else {
				participants = append(participants, "	participants.add("+p+")")
			}
		}
	}
	return strings.Join(participants, "\n")
}

func toParticipant(varname string, attr model.ResourceAttribute) string {
	datatype := GetKotlinType(attr)
	if strings.Compare(datatype, "AbstractParty") == 0 {
		return "	participants.add(" + varname + ")"
	} else if strings.Compare(datatype, "List<AbstractParty>") == 0 {
		return "	participants.addAll(" + varname + ")"
	} else {
		panic(fmt.Sprintf("attribute %s's data type is %s, must be type of participant", varname, datatype))
	}
}
func getAttr(name string, attrs []model.ResourceAttribute) (model.ResourceAttribute, bool) {
	for _, attr := range attrs {
		if attr.Name == name {
			return attr, true
		}
	}
	return model.ResourceAttribute{}, false
}
func GetFlowSha256(modelfile string) string {
	content, _ := ioutil.ReadFile(modelfile)
	sum := sha256.Sum256(content)
	return fmt.Sprintf("%x", sum)
}

func parseFlowApp(appCfg *app.Config) ([]ModelResources, map[string]string, bool, error) {
	triggerByname := make(map[string]*trigger.Config)
	isComposer := false
	isDovetail := false
	for _, t := range appCfg.Triggers {
		if t.Ref == "#transaction" {
			isComposer = true
		} else if t.Ref == "#action" {
			isDovetail = true
		}
		triggerByname[t.Ref] = t
	}

	if isComposer && isDovetail {
		return nil, nil, isComposer, fmt.Errorf("Cann't mix smart contract trigger types, only one is supported")
	}

	if !isComposer && !isDovetail {
		return nil, nil, isComposer, fmt.Errorf("There is no smart contract transactions defined")
	}

	resources := make([]ModelResources, 0)
	if isComposer {
		txnTrigger := triggerByname["#transaction"]
		model, err := processComposerTrigger(appCfg.Name, txnTrigger)
		if err != nil {
			return nil, nil, isComposer, err
		}
		resources = append(resources, model)
	} else if isDovetail {
		txnTrigger := triggerByname["#action"]
		rs, err := processDovetailTrigger(txnTrigger)
		if err != nil {
			return nil, nil, isComposer, err
		}
		resources = rs
	}

	schedulables := make(map[string]string)
	schedulerTrigger, ok := triggerByname["CordaSmartContractEventScheduler"]
	if ok {
		for _, h := range schedulerTrigger.Handlers {
			schedulableState, ok := h.Settings["assetName"]
			if !ok {
				return nil, nil, isComposer, fmt.Errorf("Schedulable asset is not defined")
			}
			schedulables[schedulableState.(string)] = model.GetFlowName(h.Actions[0].Settings["flowURI"].(string))
		}
	}
	return resources, schedulables, isComposer, nil
}
