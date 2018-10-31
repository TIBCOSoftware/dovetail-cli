/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package cmd

import (
	"fmt"
	"os"

	"github.com/TIBCOSoftware/dovetail-cli/pkg/repo"
	"github.com/spf13/cobra"
)

var force bool

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Initialize the application configuration even if it already exists (default \"false\")")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the application configuration",
	Long:  `Initialize the application configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := repo.Init(force); err != nil {
			fmt.Printf("Error initializing the configuration: '%s'", err)
			os.Exit(1)
		}
	},
}
