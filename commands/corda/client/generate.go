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
	namespace         string
	target            string
	caversion         string
	cordappmodelfile  string
	contractmodelfile string
	dependencypom     string
	streaming         string
)

func init() {
	ClientCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringP("target", "t", ".", "Destination path for generated artifacts, if a filename is given (With extension) the generated artifacts will compressed as a zip file with the file name provided")
	generateCmd.Flags().StringVarP(&namespace, "cordapp-ns", "", "", "CorDapp namespace, not required to generate generic client")
	generateCmd.Flags().StringVarP(&cordappmodelfile, "cordapp-json", "", "", "CorDApp flow json file, not required to generate generic client")
	generateCmd.Flags().StringVarP(&contractmodelfile, "smartcontract-json", "", "", "Smart Contract flow json file, not required to generate generic client")
	generateCmd.Flags().StringVarP(&dependencypom, "dependency-file", "", "", "pom snippet to include smart contract dependency")
	generateCmd.Flags().StringVarP(&streaming, "streaming", "", "none", "publish vault changes to messaging service, none or eftl")

	generateCmd.MarkFlagRequired("target")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Commands for generating client webserver artifacts",
	Long:  `Commands for generating client webserver artifacts`,
	Run: func(cmd *cobra.Command, args []string) {

		smversion, err := ClientCmd.PersistentFlags().GetString("version")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		caversion = smversion

		if cordappmodelfile != "" {
			err = validateModelFile(cordappmodelfile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = validateModelFile(contractmodelfile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if namespace == "" {
				fmt.Println("namespace is required")
				os.Exit(1)
			}

			if dependencypom == "" {
				fmt.Println("Must specify the dependency-file")
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
	if cordappmodelfile != "" && namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}

	options := cordac.NewOptions(cordappmodelfile, caversion, target, namespace, contractmodelfile, dependencypom, streaming)
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
