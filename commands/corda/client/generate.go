/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

// Package contract is the one containing all the cli commands for contract operations
package client

import (
	"fmt"
	"os"
	"path/filepath"

	cordac "github.com/TIBCOSoftware/dovetail-cli/corda/client"
	"github.com/TIBCOSoftware/dovetail-cli/model"
	"github.com/TIBCOSoftware/dovetail-cli/pkg/contract"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	namespace string
	target    string
	caversion string
	modelfile string
)

func init() {
	ClientCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringP("target", "t", ".", "Destination path for generated artifacts, if a filename is given (With extension) the generated artifacts will compressed as a zip file with the file name provided")
	generateCmd.Flags().StringP("namespace", "", "", "CorDapp namespace to generate generic client")
	generateCmd.Flags().StringVarP(&modelfile, "modelfile", "m", "", "DApp flow model file to generate generic client")

	generateCmd.MarkFlagRequired("target")
	generateCmd.MarkFlagRequired("namespace")
	generateCmd.MarkFlagRequired("modelfile")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Commands for generating dapp artifacts",
	Long:  `Commands for generating dapp artifacts`,
	Run: func(cmd *cobra.Command, args []string) {

		smversion, err := ClientCmd.PersistentFlags().GetString("version")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		caversion = smversion

		if modelfile != "" {
			err = validateModelFile(modelfile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			namespace, err = cmd.Flags().GetString("namespace")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		target, err = cmd.Flags().GetString("target")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if target == "" {
			target = "./target"
		}

		target, err = filepath.Abs(target)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		generator, err := createCordaClientGenerator()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := generator.Generate(); err != nil {
			fmt.Printf("Error generating the contract: '%s'", err)
			os.Exit(1)
		}
	},
}

func createCordaClientGenerator() (contract.Generator, error) {
	if modelfile != "" && namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	options := cordac.NewOptions(modelfile, caversion, target, namespace)
	cordaGen := cordac.NewGenerator(options)
	return cordaGen, nil
}

func validateModelFile(modelfile string) error {
	appConfig, err := model.ParseApp(modelfile)
	if err != nil {
		return errors.Wrapf(err, "Failed to parse model file %s", modelfile)
	}

	if len(appConfig.Triggers) == 0 {
		return fmt.Errorf("There must be at least one trigger defined in dapp application")
	}

	return nil
}
