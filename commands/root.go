/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

// Package commands is the one containing all the cli commands
package commands

import (
	"fmt"
	"os"

	dc "github.com/TIBCOSoftware/dovetail-cli/commands/client"
	"github.com/TIBCOSoftware/dovetail-cli/commands/contract"
	"github.com/TIBCOSoftware/dovetail-cli/commands/dapp"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(contract.ContractCmd)
	RootCmd.AddCommand(dapp.DAppCmd)
	RootCmd.AddCommand(dc.ClientCmd)
}

// RootCmd is the root command for dovetail cli
var RootCmd = &cobra.Command{
	Use:   "dovetail",
	Short: "dovetail is a flexible blockchain tool",
	Long:  `dovetail is a flexible blockchain tool`,
}

// Execute is the starting point
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
