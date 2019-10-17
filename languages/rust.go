/*
* Copyright © 2019. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package languages

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	api_lang "github.com/TIBCOSoftware/dovetail-cli/pkg/languages"
	wgutil "github.com/TIBCOSoftware/dovetail-cli/util"
)

// RustProject is an implementation of a Rust Project
type RustProject struct {
	inputTargetDir string
	isFile         bool
	tempTargetDir  string
	appName        string
	appDir         string
}

// NewRust returns a new Rust Project
func NewRust(inputTargetDir, appName string) api_lang.Project {
	return &RustProject{inputTargetDir: inputTargetDir, appName: appName}
}

// Init initializes a rust project structure
func (g *RustProject) Init() error {
	g.isFile = len(filepath.Ext(g.inputTargetDir)) > 0
	targetDir := g.inputTargetDir
	// If file is provided target is a temp location
	if g.IsFile() {
		dir, err := ioutil.TempDir("", "")
		if err != nil {
			return err
		}
		g.tempTargetDir = dir
		targetDir = g.tempTargetDir
	}
	g.appDir = wgutil.CreateTargetDirs(path.Join(targetDir, strings.ToLower(g.appName), "src"))
	return nil
}

// Cleanup removes all temp files (if any) created during the initialization
func (g *RustProject) Cleanup() error {
	if len(g.tempTargetDir) > 0 {
		err := os.RemoveAll(g.tempTargetDir)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetTargetDir returns the target folder, either the inputtarget or the temptarget if it is created
func (g *RustProject) GetTargetDir() string {
	if len(g.tempTargetDir) > 0 {
		return g.tempTargetDir
	}
	return g.inputTargetDir
}

// GetInputTargetDir returns the target folder entered by the user
func (g *RustProject) GetInputTargetDir() string {
	return g.inputTargetDir
}

// GetAppDir returns the directory of the App
func (g *RustProject) GetAppDir() string {
	return g.appDir
}

// IsFile returns true if the destination target is a File
func (g *RustProject) IsFile() bool {
	return g.isFile
}
