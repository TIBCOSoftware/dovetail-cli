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

	"github.com/TIBCOSoftware/dovetail-cli/commands/corda"
	//"github.com/TIBCOSoftware/dovetail-cli/commands/fabric"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
	//RootCmd.AddCommand(fabric.FabricCmd)
	RootCmd.AddCommand(corda.CordaCmd)
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
