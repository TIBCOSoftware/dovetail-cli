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
}

type DataState struct {
	NS            string
	Class         string
	ContractClass string
	CordaClass    string
	Attributes    []model.ResourceAttribute
	Parent        string
	Participants  []string
}

type ContractData struct {
	ContractClass string
	NS            string
	Flow          string
	Commands      []Command
	States        []DataState
}

type Command struct {
	Name       string
	NS         string
	Attributes []model.ResourceAttribute
}

var models map[string]*model.ResourceMetadataModel

// NewGenerator is the generator constructor
func NewGenerator(opts *Options) contract.Generator {
	return &Generator{Opts: opts}
}

// NewOptions is the options constructor
func NewOptions(flowModel string, version string, state string, commands []string, target, ns string) *Options {
	return &Options{ModelFile: flowModel, Version: version, State: state, Commands: commands, TargetDir: target, Namespace: ns}
}

// Generate generates a smart contract for the given options
func (g *Generator) Generate() error {
	logger.Println("Generating artifacts for Corda...")

	flow, err := model.ParseFlowApp(g.Opts.ModelFile)
	if err != nil {
		return fmt.Errorf("error parsing flow app json file, err %v", err)
	}

	models = parseAllResources(flow.Schemas)
	data, concepts, err := prepareContractStateData(g.Opts, flow, models)
	if err != nil {
		return fmt.Errorf("prepareContractStateData err %v", err)
	}

	javaProject := languages.NewJava(g.Opts.TargetDir, flow.AppName)

	err = javaProject.Init()
	if err != nil {
		return err
	}

	defer javaProject.Cleanup()

	//create directories
	resourcedir := wgutil.CreateDirIfNotExist(javaProject.GetAppDir(), "src/main/resources", strings.Replace(g.Opts.Namespace, ".", "/", -1))
	kotlindir := wgutil.CreateDirIfNotExist(javaProject.GetAppDir(), "src/main/kotlin")
	wgutil.CreateDirIfNotExist(javaProject.GetAppDir(), "target/kotlin/classes")

	//create ContractState
	for _, s := range data.States {
		err = createKotlinFile(kotlindir, g.Opts.Namespace, s, "kotlin.state.template", fmt.Sprintf("%s%s", s.Class, ".kt"))
		if err != nil {
			return fmt.Errorf("createContractStateKotlinFile kotlin.state.template err %v", err)
		}
	}

	//create Contract
	err = createKotlinFile(kotlindir, g.Opts.Namespace, data, "kotlin.contract.template", fmt.Sprintf("%s%s", data.ContractClass, "Contract.kt"))
	if err != nil {
		return fmt.Errorf("createContractJavaFile kotlin.contract.template err %v", err)
	}

	//create Resource
	err = createResourceFiles(resourcedir, g.Opts, models)
	if err != nil {
		return fmt.Errorf("createResourceFiles err %v", err)
	}

	//create common Concept files
	err = createConceptKotlinFiles(kotlindir, "kotlin.concept.template", concepts)
	if err != nil {
		return fmt.Errorf("createConceptJavaFiles kotlin.concept.template err %v", err)
	}

	err = compileAndJar(javaProject.GetAppDir(), g.Opts.Namespace, flow.AppName, g.Opts.Version, "kotlin.pom.xml")
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

func compileAndJar(targetdir, ns, clazz, version string, pomf string) error {
	logger.Printf("Compile smart contract artifacts")
	pom, err := Asset("resources/" + pomf)
	if err != nil {
		return err
	}
	err = wgutil.CopyContent(pom, path.Join(targetdir, pomf))
	if err != nil {
		return err
	}
	err = wgutil.MvnPackage(ns, clazz, version, pomf, targetdir)
	/*args := []string{"install", "-f", path.Join(targetdir, pomf), "-DbaseDir=" + targetdir, "-Dversion=" + version, "-DgroupId=" + ns, "-DartifactId=" + clazz}
	cmd := exec.Command("mvn", args...)
	logger.Printf("mvn command %v\n", cmd.Args)
	out, err := cmd.Output()*/
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

func createResourceFiles(dir string, opts *Options, models map[string]*model.ResourceMetadataModel) error {
	logger.Println("Copy resource file - transactions.json ...")

	err := wgutil.CopyFile(opts.ModelFile, path.Join(dir, "transactions.json"))
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
func prepareContractStateData(opts *Options, flow *model.ModelResources, models map[string]*model.ResourceMetadataModel) (data ContractData, conceptdata []DataState, err error) {
	logger.Println("Prepare contract state data ....")
	assets := flow.Assets
	transactions := flow.Transactions
	data = ContractData{NS: opts.Namespace, Flow: opts.ModelFile}

	fmt.Printf("***** Contract name = %s.%s%s\n", opts.Namespace, flow.AppName, "Contract")
	data.ContractClass = fmt.Sprintf("%s%s", flow.AppName, "Contract")
	conceptdata = make([]DataState, 0)
	states := make([]DataState, 0)
	concepts := make(map[string]string)

	//assets
	if opts.State != "" {
		state, process, err := prepareState(opts.State, opts.Namespace, models, concepts)
		if err != nil {
			return data, conceptdata, err
		}
		state.ContractClass = fmt.Sprintf("%s.%s%s", opts.Namespace, flow.AppName, "Contract")
		if process {
			states = append(states, state)
		}
	} else {
		for _, asset := range assets {
			metadata, ok := models[asset]
			if ok && !metadata.Metadata.IsAbstract {
				state, process, err := prepareState(asset, opts.Namespace, models, concepts)
				if err != nil {
					return data, conceptdata, err
				}
				state.ContractClass = fmt.Sprintf("%s.%s%s", opts.Namespace, flow.AppName, "Contract")
				if process {
					states = append(states, state)
				}
			}
		}
	}
	data.States = states

	//transactions
	var commands []string
	if len(opts.Commands) > 0 {
		commands = opts.Commands
	} else {
		commands = transactions
	}
	commandData, err := prepareCommands(opts.Namespace, commands, models, concepts)
	if err != nil {
		return data, conceptdata, err
	}
	data.Commands = commandData

	//common concepts
	for _, nm := range concepts {
		concept := DataState{}
		ns, clazz := splitNamespace(nm, opts.Namespace)
		concept.NS = ns
		concept.Class = clazz
		concept.Parent = models[nm].Metadata.Parent
		concept.Attributes = make([]model.ResourceAttribute, len(models[nm].Attributes))
		copy(concept.Attributes, models[nm].Attributes)
		conceptdata = append(conceptdata, concept)
	}
	return data, conceptdata, nil
}
func prepareState(asset, targetns string, models map[string]*model.ResourceMetadataModel, concepts map[string]string) (state DataState, generate bool, err error) {
	state = DataState{}
	ns, class := splitNamespace(asset, targetns)
	if ns != targetns {
		return DataState{}, false, nil
	}

	state.NS = ns
	state.Class = class
	state.CordaClass = strings.TrimSpace(getParentCordaClass(asset, models))

	state.Attributes = make([]model.ResourceAttribute, len(models[asset].Attributes))
	metadata := models[asset]
	copy(state.Attributes, metadata.Attributes)

	//check CordaParticipants decorator
	participants := make([]string, 0)
	for _, decorator := range metadata.Metadata.Decorators {
		if decorator.Name == "CordaParticipants" {
			for _, arg := range decorator.Args {
				participants = append(participants, arg)
			}
			break
		}
	}

	for idx, p := range participants {
		if !strings.HasPrefix(p, "$tx.") {
			return state, false, fmt.Errorf("CordaParticipants %s should be in the format of $tx.path.to.party", p)
		}
		participants[idx] = strings.TrimPrefix(p, "$tx.")
	}

	if len(participants) == 0 {
		//default to all top level party objects
		for _, attr := range state.Attributes {
			if attr.Type == "com.tibco.dovetail.system.Party" && attr.IsRef {
				participants = append(participants, toParticipant(attr.Name, attr))
			}
		}
	}

	if len(participants) == 0 {
		return state, false, fmt.Errorf("There must be at least one Party objects defined either as top Asset attribute or through CordaParticipants decorator")
	}

	state.Participants = participants

	for idx, attr := range state.Attributes {
		m := models[attr.Type]
		if m != nil {
			if m.Metadata.CordaClass != "" {
				state.Attributes[idx].Type = strings.TrimSpace(m.Metadata.CordaClass)
			}
		}
	}

	addCommonConcept(asset, models, concepts)
	return state, true, nil
}
func prepareCommands(targetns string, commands []string, models map[string]*model.ResourceMetadataModel, concepts map[string]string) ([]Command, error) {
	data := make([]Command, 0)
	for _, cmd := range commands {
		logger.Printf("Procesing command %s\n", cmd)
		command := Command{}
		command.Attributes = make([]model.ResourceAttribute, len(models[cmd].Attributes))
		ns, nm := splitNamespace(cmd, targetns)
		if ns != targetns {
			continue
		}
		command.Name = nm
		command.NS = ns

		m := models[cmd]
		if m == nil {
			return nil, fmt.Errorf("%s is not found in model", cmd)
		}

		isQuery := false
		for _, decorator := range m.Metadata.Decorators {
			if decorator.Name == "Query" {
				isQuery = true
				break
			}
		}

		if isQuery {
			continue
		}
		copy(command.Attributes, models[cmd].Attributes)
		data = append(data, command)

		addCommonConcept(cmd, models, concepts)
	}
	return data, nil
}
func addCommonConcept(resource string, models map[string]*model.ResourceMetadataModel, concepts map[string]string) map[string]string {
	am, ok := models[resource]
	if ok == true {
		if am.Metadata.Parent != "" {
			concepts = addCommonConcept(am.Metadata.Parent, models, concepts)
		}

		for _, attr := range am.Attributes {
			m := models[attr.Type]
			if m != nil && strings.Compare(strings.ToUpper(m.Metadata.Type), "CONCEPT") == 0 && m.Metadata.CordaClass == "" {
				concepts[attr.Type] = attr.Type
				concepts = addCommonConcept(attr.Type, models, concepts)
			}
		}
	}
	return concepts
}

func getParentCordaClass(asset string, models map[string]*model.ResourceMetadataModel) string {
	model := models[asset]
	cordaClass := ""
	if model.Metadata.Parent != "" {
		parent := models[model.Metadata.Parent]
		if parent != nil {
			if parent.Metadata.CordaClass != "" {
				cordaClass = parent.Metadata.CordaClass
			} else {
				if parent.Metadata.Parent != "" {
					return getParentCordaClass(model.Metadata.Parent, models)
				}
			}
		}
	}
	return cordaClass
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
			datatype = "java.math.BigDecimal"
			break
		case "DateTime":
			datatype = "String"
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
			datatype = "net.corda.core.contracts.Amount<Currency>"
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
				participants = append(participants, p)
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
