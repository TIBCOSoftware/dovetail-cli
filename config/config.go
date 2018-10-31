/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package config

import (
	"fmt"
	"path/filepath"

	hlconfig "github.com/TIBCOSoftware/dovetail-cli/hyperledger-fabric/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Config struct {
	Hyperledger *hlconfig.Config `json:"hyperledger"`
}

// New is a config constructor with default values
func New() *Config {
	return &Config{Hyperledger: hlconfig.New()}
}

// DefaultLocation returns the full location of the default configuration file
func DefaultLocation() (string, error) {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, DEFAULT_FOLDER_NAME, DEFAULT_FILE_NAME), nil
}

// GetConfigFile file returns the config file location
// If the config flag is set the default location is overriden
func GetConfigFile(cfgFileFlag string) (string, error) {
	if len(cfgFileFlag) == 0 {
		defaultCfgFile, err := DefaultLocation()
		if err != nil {
			return "", fmt.Errorf("Can't read config: '%s'", err)
		}
		return defaultCfgFile, nil
	}

	return cfgFileFlag, nil
}

// ReadConfigFile reads the config after setting the right file path to viper
// If the config flag is set the default location is overriden
func ReadConfigFile(cfgFileFlag string) error {
	cfgFile, err := GetConfigFile(cfgFileFlag)
	if err != nil {
		return fmt.Errorf("Can't read config: '%s'", err)
	}
	viper.SetConfigFile(cfgFile)
	if err = viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Can't read config: '%s'", err)
	}
	return nil
}

// Returns a list of the supported blockchains
func Blockchains() []string {
	return []string{HYPERLEDGER_FABRIC, CORDA}
}
