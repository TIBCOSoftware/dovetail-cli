/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

// Package contract implements generate smart contract for ethereum
package contract

import (
	"bufio"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/TIBCOSoftware/dovetail-cli/languages"
	"github.com/TIBCOSoftware/dovetail-cli/model"
	"github.com/TIBCOSoftware/dovetail-cli/pkg/contract"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

// Generator defines the generator attributes
type Generator struct {
	Opts *GenOptions
}

// GenOptions defines the generator options
type GenOptions struct {
	TargetDir string
	ModelFile string
}

// NewGenerator is the generator constructor
func NewGenerator(opts *GenOptions) contract.Generator {
	return &Generator{Opts: opts}
}

// NewGenOptions is the options constructor
func NewGenOptions(targetPath, modelFile string) *GenOptions {
	return &GenOptions{TargetDir: targetPath, ModelFile: modelFile}
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

	err = createCargoTomlFile(rustProject.GetTargetDir(), appConfig.Name)
	if err != nil {
		return err
	}

	logger.Println("Generating Ethereum smart contract... Done")
	return nil
}

func createCargoTomlFile(targetdir, appName string) error {
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
	data := CargoToml{Name: appName, Version: "0.0.1"}

	err = t.Execute(writer, data)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}
