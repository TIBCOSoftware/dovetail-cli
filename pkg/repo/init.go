/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package repo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/TIBCOSoftware/dovetail-cli/config"
	"github.com/spf13/viper"
)

// Init initializes the configuration
func Init(force bool) error {
	// Load current config
	cfgFile, err := config.GetConfigFile(viper.GetString("config"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// First read if current config exists
	exists, isDir, err := cfgFileInfo(cfgFile)
	if err != nil {
		return err
	}

	if exists {
		if !force {
			return fmt.Errorf("Configuration file in path '%s' already exist, please use the --force flag to override the current configuration", cfgFile)
		}
		if isDir {
			return fmt.Errorf("Configuration path should be a file path, not a directory")
		}
	}

	// Create path to file
	err = os.MkdirAll(filepath.Dir(cfgFile), os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error creating the default configuration file path '%s'", err)
	}

	// Create or ovewrite the config file
	cfg := config.New()
	cfgJSON, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("Error creating the default configuration '%s'", err)
	}
	err = ioutil.WriteFile(cfgFile, cfgJSON, 0644)
	if err != nil {
		return fmt.Errorf("Error writting the default configuration '%s'", err)
	}

	return nil
}

// Returns the config file information if it exists, is a directory, and error
func cfgFileInfo(cfgFile string) (bool, bool, error) {
	var fileInfo os.FileInfo
	fileInfo, err := os.Stat(cfgFile)
	if err == nil {
		// Does exist
		return true, fileInfo.IsDir(), nil
	}
	if os.IsNotExist(err) {
		// Does not exist
		return false, false, nil
	}
	// There has been an error checking the file
	return false, false, fmt.Errorf("Error checking the configuration file '%s'", err)
}
