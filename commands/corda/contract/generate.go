/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

package contract

import (
	"fmt"
	"os"
	"path/filepath"

	corda "github.com/TIBCOSoftware/dovetail-cli/corda/contract"
	"github.com/TIBCOSoftware/dovetail-cli/model"
	"github.com/TIBCOSoftware/dovetail-cli/pkg/contract"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	cordaNS   string
	target    string
	modelfile string
	pom       string
	cversion  string
)

func init() {
	ContractCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringP("target", "t", ".", "Destination path for generated artifacts, if a filename is given (With extension) the generated artifacts will compressed as a zip file with the file name provided")
	generateCmd.Flags().StringP("namespace", "", "", "Required, composer model namespace")
	generateCmd.Flags().StringVarP(&modelfile, "modelfile", "m", "", "Smart contract flow model file")
	generateCmd.Flags().StringVarP(&pom, "dependency-file", "", "", "dependency xml file")

	generateCmd.MarkFlagRequired("target")
	generateCmd.MarkFlagRequired("modelfile")
	generateCmd.MarkFlagRequired("namespace")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Commands for generating contract artifacts",
	Long:  `Commands for generating contract artifacts`,
	Run: func(cmd *cobra.Command, args []string) {

		smversion, err := ContractCmd.PersistentFlags().GetString("version")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cversion = smversion

		err = validateModelFile(modelfile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cordaNS, err = cmd.Flags().GetString("namespace")
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
			target = "./target"
		}

		target, err = filepath.Abs(target)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		generator, err := createCordaGenerator()
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

func createCordaGenerator() (contract.Generator, error) {
	if cordaNS == "" {
		return nil, fmt.Errorf("namespace is required")
	}

	options := corda.NewOptions(modelfile, cversion, target, cordaNS, pom)
	cordaGen := corda.NewGenerator(options)
	return cordaGen, nil
}

func validateModelFile(modelfile string) error {
	appConfig, err := model.ParseApp(modelfile)
	if err != nil {
		return errors.Wrapf(err, "Failed to parse model file %s", modelfile)
	}

	if len(appConfig.Triggers) == 0 {
		return fmt.Errorf("There must be at least one trigger defined in smart contract application")
	}

	return nil
}
