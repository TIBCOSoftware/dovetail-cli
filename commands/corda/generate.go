/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

package corda

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	corda "github.com/TIBCOSoftware/dovetail-cli/corda/contract"
	"github.com/TIBCOSoftware/dovetail-cli/model"
	"github.com/TIBCOSoftware/dovetail-cli/pkg/contract"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	cordaState    string
	cordaCommands string
	cordaNS       string
	target        string
	cversion      string
	modelfile     string
)

func init() {
	CordaCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringP("target", "t", ".", "Destination path for generated artifacts, if a filename is given (With extension) the generated artifacts will compressed as a zip file with the file name provided")
	generateCmd.Flags().StringP("state", "", "", "Optional, specify asset name to generate contract state, default to all assets in the specified namespace")
	generateCmd.Flags().StringP("commands", "", "", "Optional, comma delimited list of transactions(commands) allowed for the selected state txn1,txn2,..., default to all transactions")
	generateCmd.Flags().StringP("namespace", "", "", "Required, composer model namespace")
	generateCmd.Flags().StringVarP(&modelfile, "modelfile", "m", "", "Smart contract flow model file")

	generateCmd.MarkFlagRequired("target")
	generateCmd.MarkFlagRequired("modelfile")
	generateCmd.MarkFlagRequired("namespace")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Commands for generating contract artifacts",
	Long:  `Commands for generating contract artifacts`,
	Run: func(cmd *cobra.Command, args []string) {

		smversion, err := CordaCmd.PersistentFlags().GetString("version")
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

		cordaState, err = cmd.Flags().GetString("state")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cordaCommands, err = cmd.Flags().GetString("commands")
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
	cmds := make([]string, 0)
	if cordaCommands != "" {
		cmds = strings.Split(cordaCommands, ",")
		for i, v := range cmds {
			cmds[i] = strings.TrimSpace(v)
		}
	}

	options := corda.NewOptions(modelfile, cversion, cordaState, cmds, target, cordaNS)
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
