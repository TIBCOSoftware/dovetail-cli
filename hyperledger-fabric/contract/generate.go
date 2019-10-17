/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

// Package contract implements generate and deploy chaincode for hyperledger fabric
package contract

import (
	"bufio"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/TIBCOSoftware/dovetail-cli/model"
	"github.com/TIBCOSoftware/dovetail-cli/pkg/contract"
	"github.com/project-flogo/cli/api"
	_ "github.com/project-flogo/flow/util"
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

	// Remove main file
	err = os.Remove(filepath.Join(appProject.SrcDir(), "main.go"))
	if err != nil {
		return err
	}

	logger.Println("Creating the chaincode shim file...")
	err = createChaincodeShimFile(appProject.SrcDir())
	if err != nil {
		return err
	}

	logger.Println("Generating Hyperledger Fabric smart contract... Done")
	logger.Printf("Generated artifacts: '%s'\n", appProject.Dir())
	return nil
}

func createChaincodeShimFile(targetdir string) error {
	gofile := "main.go"
	logger.Printf("Create shim %s file ....\n", gofile)

	f, error := os.Create(path.Join(targetdir, gofile))
	if error != nil {
		return error
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	t, err := template.New("chaincode_shim").Parse(ChaincodeShimTemplate)
	if err != nil {
		return err
	}

	err = t.Execute(writer, nil)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}
