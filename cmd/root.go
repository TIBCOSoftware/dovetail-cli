/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package cmd

import (
	"fmt"
	"os"

	"github.com/TIBCOSoftware/dovetail-cli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFileFlag string

func init() {
	cfgFileDefault, err := config.DefaultLocation()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	RootCmd.PersistentFlags().StringVar(&cfgFileFlag, config.CONFIG_FLAG_KEY, cfgFileDefault, "Configuration file path")
	viper.BindPFlag(config.CONFIG_FLAG_KEY, RootCmd.PersistentFlags().Lookup(config.CONFIG_FLAG_KEY))
}

var RootCmd = &cobra.Command{
	Use:   "dovetail-cli",
	Short: "dovetail-cli is a flexible blockchain tool",
	Long:  `dovetail-cli is a flexible blockchain tool`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ConfigPreRun(cmd *cobra.Command, args []string) {
	err := config.ReadConfigFile(viper.GetString(config.CONFIG_FLAG_KEY))
	if err != nil {
		fmt.Println(err)
		fmt.Println("Make sure you have initialized the application configuration (Look at init command)")
		os.Exit(1)
	}
}
