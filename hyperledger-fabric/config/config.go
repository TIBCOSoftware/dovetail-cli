/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package config

var (
	orgIdsDefault = []string{"org1", "org2"}
	//filePathDefault =
)

type Config struct {
	OrgIds   []string `json:"orgIds"`
	FilePath string   `json:"filePath"`
}

// New creates a new Hyperledger Config with all default values
func New() *Config {
	return &Config{OrgIds: orgIdsDefault}
}
