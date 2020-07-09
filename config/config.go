/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

// Package config contains a set of variables and functions related to configuration
package config

import (
	"github.com/spf13/viper"
)

// Blockchains Returns a list of the supported blockchains
func Blockchains() []string {
	return []string{HyperledgerFabric, Corda}
}

// IsNodeVerbose returns whether the node log should be verbose
func IsNodeVerbose() bool {
	return viper.GetBool(NodeVerboseKey)
}
