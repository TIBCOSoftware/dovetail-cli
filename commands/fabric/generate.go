/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

package fabric

import (
	"fmt"
	"os"
	"path/filepath"

	fabric "github.com/TIBCOSoftware/dovetail-cli/hyperledger-fabric/contract"
	"github.com/TIBCOSoftware/dovetail-cli/model"
	"github.com/TIBCOSoftware/dovetail-cli/pkg/contract"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	target            string
	ccversion         string
	modelfile         string
	enableTxnSecurity bool
	dovetailMacroPath string
)

func init() {
	FabricCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringP("target", "t", ".", "Destination path for generated artifacts, if a filename is given (With extension) the generated artifacts will compressed as a zip file with the file name provided")
	generateCmd.Flags().StringVarP(&modelfile, "modelfile", "m", "", "Smart contract flow model file")
	generateCmd.Flags().BoolP("enableTransactionSecurity", "", false, "true to enable transaction level security (default false)")

	generateCmd.MarkFlagRequired("target")
	generateCmd.MarkFlagRequired("modelfile")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Commands for generating chaincode artifacts",
	Long:  `Commands for generating chaincode artifacts`,
	Run: func(cmd *cobra.Command, args []string) {

		smversion, err := FabricCmd.PersistentFlags().GetString("version")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ccversion = smversion

		err = validateModelFile(modelfile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		enableTxnSecurity, err = cmd.Flags().GetBool("enableTransactionSecurity")
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

		generator, err := createFabricGenerator()
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

func createFabricGenerator() (contract.Generator, error) {
	options := fabric.NewGenOptions(target, modelfile, ccversion, enableTxnSecurity)

	fabricGen := fabric.NewGenerator(options)
	return fabricGen, nil
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
