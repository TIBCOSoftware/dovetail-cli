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

// GoProject is an implementation of a Go Project
type GoProject struct {
	inputTargetDir string
	isFile         bool
	tempTargetDir  string
	appName        string
	appDir         string
}

// NewGo returns a new Golang Project
func NewGo(inputTargetDir, appName string) api_lang.Project {
	return &GoProject{inputTargetDir: inputTargetDir, appName: appName}
}

// Init initializes a golang project structure
func (g *GoProject) Init() error {
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
	g.appDir = wgutil.CreateTargetDirs(path.Join(targetDir, strings.ToLower(g.appName), "src", strings.ToLower(g.appName)))
	return nil
}

// Cleanup removes all temp files (if any) created during the initialization
func (g *GoProject) Cleanup() error {
	if len(g.tempTargetDir) > 0 {
		err := os.RemoveAll(g.tempTargetDir)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetTargetDir returns the target folder, either the inputtarget or the temptarget if it is created
func (g *GoProject) GetTargetDir() string {
	if len(g.tempTargetDir) > 0 {
		return g.tempTargetDir
	}
	return g.inputTargetDir
}

// GetInputTargetDir returns the target folder entered by the user
func (g *GoProject) GetInputTargetDir() string {
	return g.inputTargetDir
}

// GetAppDir returns the directory of the App
func (g *GoProject) GetAppDir() string {
	return g.appDir
}

// IsFile returns true if the destination target is a File
func (g *GoProject) IsFile() bool {
	return g.isFile
}
