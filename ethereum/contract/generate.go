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
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

const (
	dovetailDeriveID     = "dovetail_derive"
	dovetailDeriveURL    = "https://github.com/torresashjian/dovetail-rust-lib/dovetail_derive"
	dovetailDeriveBranch = "issue-1/first-app"
)

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

	modelFileName := "app.json"

	err = wgutil.CopyFile(d.Opts.ModelFile, filepath.Join(rustProject.GetTargetDir(), appConfig.Name, modelFileName))
	if err != nil {
		return err
	}

	err = createCargoTomlFile(rustProject.GetTargetDir(), d.Opts.DovetailMacroPath, appConfig)
	if err != nil {
		return err
	}

	err = createMainFile(rustProject.GetAppDir(), modelFileName, appConfig)
	if err != nil {
		return err
	}

	err = createLibFile(rustProject.GetAppDir(), modelFileName, appConfig)
	if err != nil {
		return err
	}

	logger.Println("Generating Ethereum smart contract... Done")
	return nil
}

func createCargoTomlFile(targetdir, dovetailMacroPath string, appConfig *app.Config) error {
	tomlFileName := "Cargo.toml"
	logger.Printf("Create cargo %s file ....\n", tomlFileName)

	f, error := os.Create(path.Join(targetdir, appConfig.Name, tomlFileName))
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

	data := CargoToml{Name: appConfig.Name, Version: "0.0.1", DovetailMacroPath: dovetailMacroPath, GitDependencies: dependencies}

	err = t.Execute(writer, data)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func getGitDependencies(appConfig *app.Config) ([]GitDependency, error) {
	dependencies := []GitDependency{}

	// Get internal dovetail dependencies
	dovedependencies, err := getDovetailGitDependencies(appConfig)
	if err != nil {
		return nil, err
	}
	// Get Dovetail dependencies
	dependencies = append(dependencies, dovedependencies...)

	tdependencies, err := getTriggerGitDependencies(appConfig.Triggers)
	if err != nil {
		return nil, err
	}
	// Get trigger dependencies
	dependencies = append(dependencies, tdependencies...)

	// TODO get all other dependencies (activities, etc...)

	return dependencies, nil
}

func getDovetailGitDependencies(appConfig *app.Config) ([]GitDependency, error) {
	dependencies := []GitDependency{}

	/*// Get dovetail derive
	dovetailDerive := GitDependency{ID: dovetailDeriveID, URL: dovetailDeriveURL, Branch: dovetailDeriveBranch}

	dependencies = append(dependencies, dovetailDerive)
	*/
	return dependencies, nil
}

func getTriggerGitDependencies(triggers []*trigger.Config) ([]GitDependency, error) {
	dependencies := []GitDependency{}

	// Get trigger dependencies
	for _, trigger := range triggers {
		url, err := getDependencyURL(trigger.Ref)
		if err != nil {
			return nil, err
		}
		id := getDependencyID(trigger.Ref)
		// TODO pass this as parameter before release
		branch := "issue-41/sawtooth-contrib"
		dependencies = append(dependencies, GitDependency{ID: id, URL: url, Branch: branch})
	}

	return dependencies, nil
}

func getDependencyURL(ref string) (string, error) {
	// Remove initial https://
	ref = strings.TrimPrefix(ref, "https://")
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

func createMainFile(appDir, modelFileName string, appConfig *app.Config) error {
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

	doveInternalMainDependencies, err := getDoveInternalMainDependencies()
	if err != nil {
		return err
	}

	triggerMainDependencies, err := getTriggerMainDependencies(appConfig)
	if err != nil {
		return err
	}

	data := mergeMainRs(doveInternalMainDependencies, triggerMainDependencies)

	err = t.Execute(writer, data)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func createLibFile(appDir, modelFileName string, appConfig *app.Config) error {
	libFileName := "lib.rs"
	logger.Printf("Create lib %s file ....\n", libFileName)

	f, error := os.Create(path.Join(appDir, libFileName))
	if error != nil {
		return error
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	t, err := template.New("lib_rs").Parse(LibRsTemplate)
	if err != nil {
		return err
	}

	doveInternalLibDependencies, err := getDoveInternalLibDependencies()
	if err != nil {
		return err
	}

	triggerLibDependencies, err := getTriggerLibDependencies(appConfig)
	if err != nil {
		return err
	}

	data := mergeLibRs(doveInternalLibDependencies, triggerLibDependencies)

	err = t.Execute(writer, data)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func getDoveInternalMainDependencies() ([]MainRs, error) {
	rs := []MainRs{}
	/*// Add dovetailDerive crate
	dovetailDeriveCrates := map[string]struct{}{
		dovetailDeriveID: {},
	}
	// Add dovetailDerive::app use
	dovetailDeriveUses := map[string]struct{}{
		fmt.Sprintf("%s::%s", dovetailDeriveID, "app"): {},
	}
	// Add dovetailDerive::app derive
	dovetailDeriveDerives := map[string]struct{}{
		"app": {},
	}
	dovetailDerive := MainRs{Crates: dovetailDeriveCrates, Uses: dovetailDeriveUses, Derives: dovetailDeriveDerives}
	rs = append(rs, dovetailDerive)*/
	return rs, nil
}

func getTriggerMainDependencies(appConfig *app.Config) ([]MainRs, error) {
	rs := []MainRs{}

	// Get trigger dependencies
	for _, trigger := range appConfig.Triggers {
		// Add crates
		crates := map[string]struct{}{
			getDependencyID(trigger.Ref): {},
		}
		startFn := fmt.Sprintf("%s%s", "start_", getDependencyID(trigger.Ref))
		// Add uses
		uses := map[string]struct{}{
			// Start use
			fmt.Sprintf("%s::%s", getDependencyID(trigger.Ref), startFn): {},
		}
		// Add calls
		calls := map[string]struct{}{
			startFn: {},
		}
		// Add derives
		derives := map[string]struct{}{}
		newrs := MainRs{Crates: crates, Uses: uses, Calls: calls, Derives: derives}
		rs = append(rs, newrs)
	}
	return rs, nil
}

func mergeMaps(maps ...map[string]struct{}) map[string]struct{} {
	result := make(map[string]struct{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func mergeMainRs(rsAs, rsBs []MainRs) MainRs {
	crates := make(map[string]struct{})
	uses := make(map[string]struct{})
	calls := make(map[string]struct{})
	derives := make(map[string]struct{})
	for _, rsA := range rsAs {
		crates = mergeMaps(crates, rsA.Crates)
		uses = mergeMaps(uses, rsA.Uses)
		calls = mergeMaps(calls, rsA.Calls)
		derives = mergeMaps(derives, rsA.Derives)
	}
	for _, rsB := range rsBs {
		crates = mergeMaps(crates, rsB.Crates)
		uses = mergeMaps(uses, rsB.Uses)
		calls = mergeMaps(calls, rsB.Calls)
		derives = mergeMaps(derives, rsB.Derives)
	}
	return MainRs{
		Crates:  crates,
		Uses:    uses,
		Calls:   calls,
		Derives: derives,
	}
}

func mergeLibRs(rsAs, rsBs []LibRs) LibRs {
	crates := make(map[string]struct{})
	uses := make(map[string]struct{})
	calls := make(map[string]struct{})
	derives := make(map[string]struct{})
	for _, rsA := range rsAs {
		crates = mergeMaps(crates, rsA.Crates)
		uses = mergeMaps(uses, rsA.Uses)
		calls = mergeMaps(calls, rsA.Calls)
		derives = mergeMaps(derives, rsA.Derives)
	}
	for _, rsB := range rsBs {
		crates = mergeMaps(crates, rsB.Crates)
		uses = mergeMaps(uses, rsB.Uses)
		calls = mergeMaps(calls, rsB.Calls)
		derives = mergeMaps(derives, rsB.Derives)
	}
	return LibRs{
		Crates:  crates,
		Uses:    uses,
		Calls:   calls,
		Derives: derives,
	}
}

func getDoveInternalLibDependencies() ([]LibRs, error) {
	rs := []LibRs{}
	/*// Add dovetailDerive crate
	dovetailDeriveCrates := map[string]struct{}{
		dovetailDeriveID: {},
	}
	// Add dovetailDerive::app use
	dovetailDeriveUses := map[string]struct{}{
		fmt.Sprintf("%s::%s", dovetailDeriveID, "app"): {},
	}
	dovetailDerive := LibRs{Crates: dovetailDeriveCrates, Uses: dovetailDeriveUses}
	rs = append(rs, dovetailDerive)*/
	return rs, nil
}

func getTriggerLibDependencies(appConfig *app.Config) ([]LibRs, error) {
	rs := []LibRs{}

	// Get trigger dependencies
	for _, trigger := range appConfig.Triggers {
		// Add crates
		crates := map[string]struct{}{
			getDependencyID(trigger.Ref): {},
		}
		startFn := fmt.Sprintf("%s%s", "start_", getDependencyID(trigger.Ref))
		// Add uses
		uses := map[string]struct{}{
			// Start use
			fmt.Sprintf("%s::%s", getDependencyID(trigger.Ref), startFn): {},
		}
		// Add calls
		calls := map[string]struct{}{}
		// Add derives
		derives := map[string]struct{}{}
		newrs := LibRs{Crates: crates, Uses: uses, Calls: calls, Derives: derives}
		rs = append(rs, newrs)
	}
	return rs, nil
}
