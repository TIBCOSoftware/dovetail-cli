/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

// Package dapp is the one containing all the cli commands for contract operations
package dapp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TIBCOSoftware/dovetail-cli/config"
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
	caversion  string
	modelfile  string
	pom        string
	apiOnly    bool
)

func init() {
	DAppCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringP("target", "t", ".", "Destination path for generated artifacts, if a filename is given (With extension) the generated artifacts will compressed as a zip file with the file name provided")
	generateCmd.Flags().StringP("namespace", "", "", "Required, namespace of the dapp")
	generateCmd.Flags().StringVarP(&modelfile, "modelfile", "m", "", "DApp flow model file")
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

		smversion, err := DAppCmd.PersistentFlags().GetString("version")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		caversion = smversion

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
			target = "./target"
		}
		target, err = filepath.Abs(target)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		g, err := createCordAppGenerator()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := g.Generate(); err != nil {
			fmt.Printf("Error generating the contract: '%s'", err)
			os.Exit(1)
		}

	},
}

//GetGenerators chooses the right generator
func GetGenerators(blockchain string) ([]contract.Generator, error) {
	generators := make([]contract.Generator, 0)
	switch strings.ToUpper(blockchain) {
	case strings.ToUpper(config.Corda):
		g, err := createCordAppGenerator()
		if err != nil {
			return nil, err
		}

		generators = append(generators, g)
		/*
			if !apiOnly {
				gc, err := createCordaClientGenerator()
				if err != nil {
					return nil, err
				}
				generators = append(generators, gc)
			}
		*/
		return generators, nil
	default:
		return nil, fmt.Errorf("Unsupported blockchain to create dapp '%s'", blockchain)
	}
}

func createCordAppGenerator() (contract.Generator, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}

	options := corda.NewOptions(modelfile, caversion, target, namespace, pom, apiOnly)
	cordaGen := corda.NewGenerator(options)
	return cordaGen, nil
}

/*
func createCordaClientGenerator() (contract.Generator, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}

	options := cordac.NewOptions(modelfile, caversion, target, namespace)
	cordaGen := cordac.NewGenerator(options)
	return cordaGen, nil
}
*/
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
