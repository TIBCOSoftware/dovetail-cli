/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

// Package contract implements generate smart contract for ethereum
package contract

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/TIBCOSoftware/dovetail-cli/languages"
	"github.com/TIBCOSoftware/dovetail-cli/model"
	"github.com/TIBCOSoftware/dovetail-cli/pkg/contract"
	wgutil "github.com/TIBCOSoftware/dovetail-cli/util"
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
	DovetailMacroPath string
}

// NewGenerator is the generator constructor
func NewGenerator(opts *GenOptions) contract.Generator {
	return &Generator{Opts: opts}
}

// NewGenOptions is the options constructor
func NewGenOptions(targetPath, modelFile, dovetailMacroPath string) *GenOptions {
	return &GenOptions{TargetDir: targetPath, ModelFile: modelFile, DovetailMacroPath: dovetailMacroPath}
}

// Generate generates a smart contract for the given options
func (d *Generator) Generate() error {
	logger.Println("Generating Ethereum smart contract...")

	appConfig, err := model.ParseApp(d.Opts.ModelFile)
	if err != nil {
		return err
	}

	rustProject := languages.NewRust(d.Opts.TargetDir, appConfig.Name)

	err = rustProject.Init()
	if err != nil {
		return err
	}

	defer rustProject.Cleanup()

	modelFileName := fmt.Sprintf("%s.%s", appConfig.Name, "json")

	err = wgutil.CopyFile(d.Opts.ModelFile, filepath.Join(rustProject.GetAppDir(), modelFileName))
	if err != nil {
		return err
	}

	err = createCargoTomlFile(rustProject.GetTargetDir(), appConfig.Name, d.Opts.DovetailMacroPath, appConfig)
	if err != nil {
		return err
	}

	err = createMainFile(rustProject.GetAppDir(), appConfig.Name, modelFileName)
	if err != nil {
		return err
	}

	logger.Println("Generating Ethereum smart contract... Done")
	return nil
}

func createCargoTomlFile(targetdir, appName, dovetailMacroPath string, appConfig *app.Config) error {
	tomlFileName := "Cargo.toml"
	logger.Printf("Create cargo %s file ....\n", tomlFileName)

	f, error := os.Create(path.Join(targetdir, appName, tomlFileName))
	if error != nil {
		return error
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	t, err := template.New("cargo_toml").Parse(CargoTomlTemplate)
	if err != nil {
		return err
	}

	dependencies, err := getGitDependencies(appConfig)
	if err != nil {
		return err
	}

	data := CargoToml{Name: appName, Version: "0.0.1", DovetailMacroPath: dovetailMacroPath, GitDependencies: dependencies}

	err = t.Execute(writer, data)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func getGitDependencies(appConfig *app.Config) ([]GitDependency, error) {
	dependencies := []GitDependency{}
	// Get trigger dependencies
	for _, trigger := range appConfig.Triggers {
		url, err := getDependencyURL(trigger.Ref)
		if err != nil {
			return nil, err
		}
		id := getDependencyID(trigger.Ref)
		// TODO remove ethereum branch
		branch := "ethereum"
		dependencies = append(dependencies, GitDependency{ID: id, URL: url, Branch: branch})
	}
	return dependencies, nil
}

func getDependencyURL(ref string) (string, error) {
	seg := strings.Split(ref, "/")
	if len(seg) < 3 {
		return "", fmt.Errorf("Invalid dependency URL %s", ref)
	}
	dependencyURL := fmt.Sprintf("%s/%s/%s", seg[0], seg[1], seg[2])

	return fmt.Sprintf("https://%s", dependencyURL), nil
}

func getDependencyID(ref string) string {
	// Get last element of ref
	return filepath.Base(ref)
}

func createMainFile(appDir, appName, modelFileName string) error {
	mainFileName := "main.rs"
	logger.Printf("Create main %s file ....\n", mainFileName)

	f, error := os.Create(path.Join(appDir, mainFileName))
	if error != nil {
		return error
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	t, err := template.New("main_rs").Parse(MainRsTemplate)
	if err != nil {
		return err
	}

	data := MainRs{ModelPath: path.Join("src", modelFileName)}

	err = t.Execute(writer, data)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}
