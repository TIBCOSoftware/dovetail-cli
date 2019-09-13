/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

// Package dapp is the one containing all the cli commands for dapp operations
package dapp

import (
	"github.com/spf13/cobra"
)

func init() {
	DAppCmd.PersistentFlags().StringP("version", "v", "v1.0.0", "dapp version")
}

// DAppCmd is the command for distributed app
var DAppCmd = &cobra.Command{
	Use:   "dapp",
	Short: "Commands for Distributed Apps",
	Long:  `Commands for Distributed Apps`,
}
