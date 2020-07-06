/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

// Package cmd is the one containing all the cli commands
package cmd

import (
	"fmt"
	"os"

	"github.com/TIBCOSoftware/dovetail-cli/cmd/corda"
	"github.com/TIBCOSoftware/dovetail-cli/cmd/fabric"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(fabricCmd)
	rootCmd.AddCommand(cordaCmd)
}

// RootCmd is the root command for dovetail cli
var rootCmd = &cobra.Command{
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
