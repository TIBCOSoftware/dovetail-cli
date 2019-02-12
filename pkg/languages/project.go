/*
* Copyright © 2019. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package project

// Project is an interface that represents a project for each of the languages supported
type Project interface {
	// Init creates initial structure
	Init() error
	// Cleanup removes all temp directories and files
	Cleanup() error
	// GetTarget returns the target folder
	GetTargetDir() string
	// GetInputTarget returns the target folder entered by the user
	GetInputTargetDir() string
	// GetAppDir returns the directory of the App
	GetAppDir() string
	// IsFile returns true if the destination target is a File
	IsFile() bool
}
