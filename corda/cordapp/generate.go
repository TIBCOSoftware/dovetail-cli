/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package cordapp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/TIBCOSoftware/dovetail-cli/files"
	"github.com/TIBCOSoftware/dovetail-cli/languages"
	"github.com/TIBCOSoftware/dovetail-cli/model"
	wgutil "github.com/TIBCOSoftware/dovetail-cli/util"
	"github.com/TIBCOSoftware/flogo-lib/app"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

type Generator struct {
	Opts *Options
}

type Options struct {
	ModelFile      string
	Version        string
	TargetDir      string
	Namespace      string
	DependencyFile string
	APIOnly        bool
}

type InitiatorFlowConfig struct {
	Attrs                     []model.ResourceAttribute
	HasObservers              bool
	SendTxnToObserverManually bool
}
type DataState struct {
	NS                string
	App               string
	InitiatorFlows    map[string]InitiatorFlowConfig
	ResponderFlows    map[string]string
	ObserverFlows     map[string]string
	SchedulableFlows  map[string]SchedulableFlowConfig
	ConfidentialFlows map[string]bool
}

type SchedulableFlowConfig struct {
	NS               string
	App              string
	FlowName         string
	InitiatingFlowNS string
}

var models map[string]*model.ResourceMetadataModel

// NewGenerator is the generator constructor
func NewGenerator(opts *Options) *Generator {
	return &Generator{Opts: opts}
}

// NewOptions is the options constructor
func NewOptions(flowModel string, version string, target, ns string, dependencyFile string, apiOnly bool) *Options {
	return &Options{ModelFile: flowModel, Version: version, TargetDir: target, Namespace: ns, DependencyFile: dependencyFile, APIOnly: apiOnly}
}

// Generate generates a CordAppfor the given options
func (g *Generator) Generate() error {
	logger.Println("Generating artifacts for CordApp...")

	app, err := model.ParseApp(g.Opts.ModelFile)
	if err != nil {
		return fmt.Errorf("error parsing flow app json file, err %v", err)
	}

	return g.GenerateApp(app)
}
func (g *Generator) GenerateApp(app *app.Config) error {
	data, err := prepareData(g.Opts, app)
	if err != nil {
		return fmt.Errorf("prepareData err %v", err)
	}
	fmt.Printf("template data=%v\n", data)
	javaProject := languages.NewJava(g.Opts.TargetDir, data.App)

	err = javaProject.Init()
	if err != nil {
		return err
	}

	defer javaProject.Cleanup()

	//create directories
	resourcedir := wgutil.CreateDirIfNotExist(javaProject.GetAppDir(), "src/main/resources", strings.Replace(g.Opts.Namespace, ".", "/", -1))
	kotlindir := wgutil.CreateDirIfNotExist(javaProject.GetAppDir(), "src/main/kotlin")
	wgutil.CreateDirIfNotExist(javaProject.GetAppDir(), "target/kotlin/classes")

	//create app file
	err = createKotlinFile(kotlindir, g.Opts.Namespace, data, "app.template", fmt.Sprintf("%s%s", data.App, ".kt"))
	if err != nil {
		return fmt.Errorf("createAppKotlinFile app.template err %v", err)
	}

	//create abstract flow files
	if len(data.InitiatorFlows) > 0 {
		err = createKotlinFile(kotlindir, g.Opts.Namespace, data, "abstractflow.template", "AbstractFlows.kt")
		if err != nil {
			return fmt.Errorf("createFlowKotlinFile abstractflow.template err %v", err)
		}
	}

	if !g.Opts.APIOnly {
		//create flow files
		err = createKotlinFile(kotlindir, g.Opts.Namespace, data, "flow.template", "Flows.kt")
		if err != nil {
			return fmt.Errorf("createFlowKotlinFile flow.template err %v", err)
		}

		//schedulable flow file
		if len(data.SchedulableFlows) > 0 {
			for _, s := range data.SchedulableFlows {
				err = createKotlinFile(kotlindir, s.NS, s, "schedulable.template", fmt.Sprintf("%s.kt", s.FlowName))
				if err != nil {
					return fmt.Errorf("createFlowKotlinFile schedulable.template err %v", err)
				}
			}
		}
	}

	//create Resource
	err = createResourceFiles(resourcedir, g.Opts, models)
	if err != nil {
		return fmt.Errorf("createResourceFiles err %v", err)
	}

	err = compileAndJar(javaProject.GetAppDir(), g.Opts.Namespace, data.App, g.Opts.Version, "kotlin.pom.xml", g.Opts.DependencyFile)
	if err != nil {
		return err
	}

	if len(data.InitiatorFlows) > 0 {
		err = wgutil.MvnInstall(g.Opts.Namespace, data.App+"-api", g.Opts.Version, fmt.Sprintf("%s/kotlin-%s-%s-api.jar", javaProject.GetAppDir(), data.App, g.Opts.Version))
		if err != nil {
			return err
		}
	}

	if !g.Opts.APIOnly {
		err = wgutil.MvnInstall(g.Opts.Namespace, data.App, g.Opts.Version, fmt.Sprintf("%s/kotlin-%s-%s.jar", javaProject.GetAppDir(), data.App, g.Opts.Version))
		if err != nil {
			return err
		}
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

	logger.Printf("Finished generating artifacts for CordApp")
	return nil
}

func compileAndJar(targetdir, ns, clazz, version string, pomf string, dependencyFile string) error {
	logger.Printf("Compile cordapp artifacts")
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

	return nil
}

func createResourceFiles(dir string, opts *Options, models map[string]*model.ResourceMetadataModel) error {
	logger.Println("Copy resource file - app.json ...")

	err := wgutil.CopyFile(opts.ModelFile, path.Join(dir, "app.json"))
	if err != nil {
		return fmt.Errorf("Error creating app.json file, error %v", err)
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

	tmpl, err := Asset("resources/" + templateFile)
	if err != nil {
		return err
	}

	t, err := template.New(templateFile).Parse(string(tmpl))
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
	data = DataState{NS: opts.Namespace, App: app.Name, InitiatorFlows: make(map[string]InitiatorFlowConfig), ResponderFlows: make(map[string]string), SchedulableFlows: make(map[string]SchedulableFlowConfig), ConfidentialFlows: make(map[string]bool), ObserverFlows: make(map[string]string)}

	for _, trigger := range app.Triggers {
		flowType := trigger.Settings["flowType"]
		if flowType != nil {
			if flowType.(string) == "initiator" {
				for _, handler := range trigger.Handlers {
					flowName := model.GetFlowName(handler.Action.Data)
					data.ConfidentialFlows[flowName] = handler.Settings["useAnonymousIdentity"].(bool)
					config := InitiatorFlowConfig{Attrs: make([]model.ResourceAttribute, 0)}
					config.HasObservers = handler.Settings["hasObservers"].(bool)
					config.SendTxnToObserverManually = handler.Settings["observerManual"].(bool)

					input := handler.Outputs["transactionInput"]
					if input != nil {
						metadata := input.(map[string]interface{})["metadata"].(string)
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
			} else if flowType.(string) == "receiver" {
				for _, handler := range trigger.Handlers {
					flowName := model.GetFlowName(handler.Action.Data)
					data.ResponderFlows[flowName] = handler.Settings["initiatorFlow"].(string)
					data.ConfidentialFlows[flowName] = handler.Settings["useAnonymousIdentity"].(bool)
				}
			} else if flowType.(string) == "observer" {
				for _, handler := range trigger.Handlers {
					flowName := model.GetFlowName(handler.Action.Data)
					data.ObserverFlows[flowName] = handler.Settings["initiatorFlow"].(string)
				}
			} else if flowType.(string) == "schedulable" {
				for _, handler := range trigger.Handlers {
					flowName := model.GetFlowName(handler.Action.Data)
					asset := handler.Settings["asset"].(string)
					ns := asset[:strings.LastIndex(asset, ".")]
					data.SchedulableFlows[flowName] = SchedulableFlowConfig{NS: ns, App: app.Name, FlowName: flowName, InitiatingFlowNS: opts.Namespace}
				}
			}
		}
	}
	return data, nil
}
