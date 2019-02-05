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

// JavaProject is an implementation of a Java Project
type JavaProject struct {
	inputTargetDir string
	isFile         bool
	tempTargetDir  string
	appName        string
	appDir         string
}

// NewJava returns a new Java Project
func NewJava(inputTargetDir, appName string) api_lang.Project {
	return &JavaProject{inputTargetDir: inputTargetDir, appName: appName}
}

// Init initializes a java project structure
func (j *JavaProject) Init() error {
	j.isFile = len(filepath.Ext(j.inputTargetDir)) > 0
	targetDir := j.inputTargetDir
	// If file is provided target is a temp location
	if j.IsFile() {
		dir, err := ioutil.TempDir("", "")
		if err != nil {
			return err
		}
		j.tempTargetDir = dir
		targetDir = j.tempTargetDir
	}
	j.appDir = wgutil.CreateTargetDirs(path.Join(targetDir, strings.ToLower(j.appName)))
	return nil
}

// Cleanup removes all temp files (if any) created during the initialization
func (j *JavaProject) Cleanup() error {
	if len(j.tempTargetDir) > 0 {
		err := os.RemoveAll(j.tempTargetDir)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetTargetDir returns the target folder, either the inputtarget or the temptarget if it is created
func (j *JavaProject) GetTargetDir() string {
	if len(j.tempTargetDir) > 0 {
		return j.tempTargetDir
	}
	return j.inputTargetDir
}

// GetInputTargetDir returns the target folder entered by the user
func (j *JavaProject) GetInputTargetDir() string {
	return j.inputTargetDir
}

// GetAppDir returns the directory of the App
func (j *JavaProject) GetAppDir() string {
	return j.appDir
}

// IsFile returns true if the destination target is a File
func (j *JavaProject) IsFile() bool {
	return j.isFile
}
