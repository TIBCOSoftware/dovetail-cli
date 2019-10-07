/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/TIBCOSoftware/dovetail-cli/files"
	"github.com/TIBCOSoftware/dovetail-cli/languages"
	"github.com/TIBCOSoftware/dovetail-cli/model"
	wgutil "github.com/TIBCOSoftware/dovetail-cli/util"
	"github.com/project-flogo/core/app"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

type Generator struct {
	Opts *Options
}

type Options struct {
	CorDAppModelFile  string
	CorDAppVersion    string
	TargetDir         string
	CordAppNamespace  string
	ContractModelFile string
	DependencyFile    string
}

type InitiatorFlowConfig struct {
	Attrs []model.ResourceAttribute
}

type AllAssets struct {
	NS     string
	Assets []string
}

type DataState struct {
	NS             string
	App            string
	InitiatorFlows map[string]InitiatorFlowConfig
	Assets         AllAssets
}

var models map[string]*model.ResourceMetadataModel

// NewGenerator is the generator constructor
func NewGenerator(opts *Options) *Generator {
	return &Generator{Opts: opts}
}

// NewOptions is the options constructor
func NewOptions(cordappModel, version, target, ns, contractModelFile, depFile string) *Options {
	return &Options{CorDAppModelFile: cordappModel, CorDAppVersion: version, TargetDir: target, CordAppNamespace: ns, ContractModelFile: contractModelFile, DependencyFile: depFile}
}

// Generate generates a CordAppfor the given options
func (g *Generator) Generate() error {
	logger.Println("Generating artifacts for corda client...")
	data := DataState{}
	if g.Opts.CorDAppModelFile != "" {
		app, err := model.ParseApp(g.Opts.CorDAppModelFile)
		if err != nil {
			return fmt.Errorf("error parsing flow app json file, err %v", err)
		}
		data, err = prepareData(g.Opts, app)
		if err != nil {
			return fmt.Errorf("prepareData err %v", err)
		}

		assets, err := getAllAssets(g.Opts.ContractModelFile)
		if err != nil {
			return fmt.Errorf("getAllAssets err %v", err)
		}
		assets.NS = g.Opts.CordAppNamespace
		data.Assets = assets

	} else {
		data.App = "DovetailCordAppClient"
		data.NS = "com.tibco.dovetail.corda.client.webserver"
	}
	return g.GenerateApp(data)
}
func (g *Generator) GenerateApp(data DataState) error {

	javaProject := languages.NewJava(path.Join(g.Opts.TargetDir, data.App), "client")

	err := javaProject.Init()
	if err != nil {
		return err
	}

	defer javaProject.Cleanup()

	//create directories
	kotlindir := wgutil.CreateDirIfNotExist(javaProject.GetAppDir(), "src/main/kotlin")
	wgutil.CreateDirIfNotExist(javaProject.GetAppDir(), "target/kotlin/classes")
	webdir := wgutil.CreateDirIfNotExist(kotlindir, strings.Replace(g.Opts.CordAppNamespace, ".", "/", -1), "client/webserver")
	cdir := wgutil.CreateDirIfNotExist(webdir, "controller")
	servicedir := wgutil.CreateDirIfNotExist(webdir, "service")

	//create custom controller file
	if len(data.InitiatorFlows) > 0 {
		err = createKotlinFile(cdir, data, "CustomController.template", fmt.Sprintf("%s%s", data.App, ".kt"))
		if err != nil {
			return fmt.Errorf("createKotlinFile CustomController.template err %v", err)
		}
	}

	if g.Opts.ContractModelFile != "" {
		err = createKotlinFile(servicedir, data.Assets, "eftl.template", "eftl.kt")
		if err != nil {
			return fmt.Errorf("createKotlinFile eftl.template err %v", err)
		}
		err = createKotlinFile(servicedir, data.Assets, "StatesTracker.template", "StatesTracker.kt")
		if err != nil {
			return fmt.Errorf("createKotlinFile StatesTracker.template err %v", err)
		}
	}
	err = createKotlinFile(cdir, data, "ServerController.template", "ServerController.kt")
	if err != nil {
		return fmt.Errorf("createKotlinFile ServerController.template err %v", err)
	}

	err = createKotlinFile(cdir, data, "SecurityController.template", "SecurityController.kt")
	if err != nil {
		return fmt.Errorf("createKotlinFile SecurityController.template err %v", err)
	}

	err = createKotlinFile(cdir, data, "QueryController.template", "QueryController.kt")
	if err != nil {
		return fmt.Errorf("createKotlinFile QueryController.template err %v", err)
	}

	err = createKotlinFile(cdir, data, "CashController.template", "CashController.kt")
	if err != nil {
		return fmt.Errorf("createKotlinFile CashController.template err %v", err)
	}

	err = createKotlinFile(cdir, data, "Common.template", "Common.kt")
	if err != nil {
		return fmt.Errorf("createKotlinFile Common.template err %v", err)
	}

	err = createKotlinFile(cdir, data, "FilterCriteriaBuilder.template", "FilterCriteriaBuilder.kt")
	if err != nil {
		return fmt.Errorf("createKotlinFile FilterCriteriaBuilder.template err %v", err)
	}

	err = createKotlinFile(cdir, data, "AccessControl.template", "AccessControl.kt")
	if err != nil {
		return fmt.Errorf("createKotlinFile AccessControl.template err %v", err)
	}

	err = createKotlinFile(webdir, data, "Server.template", "Server.kt")
	if err != nil {
		return fmt.Errorf("createKotlinFile Server.template err %v", err)
	}

	err = createKotlinFile(webdir, data, "NodeRPCConnection.template", "NodeRPCConnection.kt")
	if err != nil {
		return fmt.Errorf("createKotlinFile NodeRPCConnection.template err %v", err)
	}

	err = createKotlinFile(webdir, data, "SwaggerConfig.template", "SwaggerConfig.kt")
	if err != nil {
		return fmt.Errorf("createKotlinFile SwaggerConfig.template err %v", err)
	}

	pom := "kotlin.pom.xml"
	if g.Opts.CorDAppModelFile == "" {
		pom = "kotlin.pom.generic.xml"
	}

	err = compileAndJar(javaProject.GetAppDir(), data.NS, data.App, g.Opts.CorDAppVersion, g.Opts.DependencyFile, pom)
	if err != nil {
		return err
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

	logger.Printf("Finished generating artifacts for corda client")
	return nil
}

func compileAndJar(targetdir, ns, clazz, version, deppom, pomf string) error {
	logger.Printf("Compile corda client artifacts")
	pom, err := Asset("resources/" + pomf)
	if err != nil {
		return err
	}

	err = wgutil.CreateNewPom(pom, targetdir, deppom, pomf)
	if err != nil {
		return err
	}

	err = wgutil.MvnPackage(ns, clazz, version, pomf, targetdir)
	if err != nil {
		return err
	}

	return nil
}

func createKotlinFile(dir string, data interface{}, templateFile string, fileName string) error {
	logger.Printf("Create kotlin file %s with template %s....", fileName, templateFile)

	f, error := os.Create(path.Join(dir, fileName))
	if error != nil {
		return error
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	tmpl, err := Asset("resources/" + templateFile)
	if err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"ToLower": ToLower,
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

func prepareData(opts *Options, app *app.Config) (data DataState, err error) {
	logger.Println("Prepare data ....")
	data = DataState{NS: opts.CordAppNamespace, App: app.Name, InitiatorFlows: make(map[string]InitiatorFlowConfig)}

	for _, trigger := range app.Triggers {
		flowType := trigger.Settings["flowType"]
		if flowType != nil {
			if flowType.(string) == "initiator" {
				for _, handler := range trigger.Handlers {
					flowName := model.GetFlowName(handler.Actions[0].Settings["flowURI"].(string))
					config := InitiatorFlowConfig{Attrs: make([]model.ResourceAttribute, 0)}

					input := handler.Schemas.Output["transactionInput"]
					if input != nil {
						metadata := input.(map[string]interface{})["value"].(string)
						if metadata != "" {
							desc := struct {
								Description string `json: "description"`
							}{}
							json.Unmarshal([]byte(metadata), &desc)
							config.Attrs = model.ParseResourceModel(desc.Description).Attributes
						}
					}
					data.InitiatorFlows[flowName] = config
				}
			}
		}
	}
	return data, nil
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func getAllAssets(contractmodel string) (AllAssets, error) {
	assetnames := AllAssets{Assets: make([]string, 0)}
	app, err := model.ParseApp(contractmodel)
	if err != nil {
		return assetnames, err
	}
	assets := make(map[string]bool)
	for _, t := range app.Triggers {
		for _, h := range t.Handlers {
			asset := h.Settings["assetname"].(string)
			assets[asset] = true
		}
	}
	for k := range assets {
		assetnames.Assets = append(assetnames.Assets, k)
	}
	return assetnames, nil
}
