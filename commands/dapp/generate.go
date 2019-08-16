/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

// Package contract is the one containing all the cli commands for contract operations
package dapp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TIBCOSoftware/dovetail-cli/config"
	cordac "github.com/TIBCOSoftware/dovetail-cli/corda/client"
	corda "github.com/TIBCOSoftware/dovetail-cli/corda/cordapp"
	"github.com/TIBCOSoftware/dovetail-cli/model"
	"github.com/TIBCOSoftware/dovetail-cli/pkg/contract"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	namespace  string
	target     string
	blockchain string
	smversion  string
	modelfile  string
	pom        string
	apiOnly    bool
)

func init() {
	DAppCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringP("target", "t", ".", "Destination path for generated artifacts, if a filename is given (With extension) the generated artifacts will compressed as a zip file with the file name provided")
	generateCmd.Flags().StringP("namespace", "", "", "Corda only, required")
	generateCmd.Flags().StringVarP(&modelfile, "model-file", "m", "", "DApp flow model file")
	generateCmd.Flags().StringVarP(&pom, "dependency-file", "", "", "dependency xml file")
	generateCmd.Flags().BoolVarP(&apiOnly, "api", "", false, "Corda only, generate API artifacts only")

	generateCmd.MarkFlagRequired("target")
	generateCmd.MarkFlagRequired("modelfile")
	generateCmd.MarkFlagRequired("namespace")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Commands for generating dapp artifacts",
	Long:  `Commands for generating dapp artifacts`,
	Run: func(cmd *cobra.Command, args []string) {
		blockchain, err := DAppCmd.PersistentFlags().GetString("blockchain")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		smversion, err = DAppCmd.PersistentFlags().GetString("version")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

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

		target, err = cmd.Flags().GetString("target")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if target == "" {
			target = "./dovetail_generated"
		}
		target, err = filepath.Abs(target)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		generators, err := GetGenerators(blockchain)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, generator := range generators {
			if err := generator.Generate(); err != nil {
				fmt.Printf("Error generating the contract: '%s'", err)
				os.Exit(1)
			}
		}
	},
}

// GetGenerator chooses the right generator
func GetGenerators(blockchain string) ([]contract.Generator, error) {
	generators := make([]contract.Generator, 0)
	switch strings.ToUpper(blockchain) {
	case strings.ToUpper(config.CORDA):
		g, err := createCordAppGenerator()
		if err != nil {
			return nil, err
		}

		generators = append(generators, g)

		gc, err := createCordaClientGenerator()
		if err != nil {
			return nil, err
		}
		generators = append(generators, gc)
		return generators, nil
	default:
		return nil, fmt.Errorf("Unsupported blockchain to create dapp '%s'", blockchain)
	}
}

func createCordAppGenerator() (contract.Generator, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}

	options := corda.NewOptions(modelfile, smversion, target, namespace, pom, apiOnly)
	cordaGen := corda.NewGenerator(options)
	return cordaGen, nil
}

func createCordaClientGenerator() (contract.Generator, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}

	options := cordac.NewOptions(modelfile, smversion, target, namespace)
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
