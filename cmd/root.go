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

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(fabricCmd)
	rootCmd.AddCommand(cordaCmd)
	rootCmd.AddCommand(nodeCmd)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

var (
	cfgFile string

	// Verbose argument for verbose output
	Verbose bool

	// rootCmd is the root command for dovetail cli
	rootCmd = &cobra.Command{
		Use:   "dovetail",
		Short: "dovetail is a flexible blockchain tool",
		Long:  `dovetail is a flexible blockchain tool`,
	}
)

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		// Search config in home directory with name "config" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// Execute is the starting point
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
