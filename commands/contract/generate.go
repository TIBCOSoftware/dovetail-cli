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
	"strings"

	"github.com/TIBCOSoftware/dovetail-cli/config"
	corda "github.com/TIBCOSoftware/dovetail-cli/corda/contract"
	fabric "github.com/TIBCOSoftware/dovetail-cli/hyperledger-fabric/contract"
	"github.com/TIBCOSoftware/dovetail-cli/model"
	"github.com/TIBCOSoftware/dovetail-cli/pkg/contract"
	"github.com/spf13/cobra"
)

var cordaState string
var cordaCommands string
var cordaNS string
var target string
var blockchain string
var smversion string
var modelfile string
var enableTxnSecurity bool

func init() {
	ContractCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringP("target", "t", ".", "Destination path for generated artifacts")
	generateCmd.Flags().StringP("state", "", "", "Corda only, optional, specify asset name to generate contract state, default to all assets in the specified namespace")
	generateCmd.Flags().StringP("commands", "", "", "Corda only, optional, comma delimited list of transactions(commands) allowed for the selected state txn1,txn2,..., default to all transactions")
	generateCmd.Flags().StringP("namespace", "", "", "Corda only, required, composer model namespace")
	generateCmd.Flags().BoolP("enableTransactionSecurity", "", false, "true to enable transaction level security for the targetd blockchain if supported")

	generateCmd.MarkFlagRequired("target")
}

var generateCmd = &cobra.Command{
	Use:              "generate",
	Short:            "Commands for generating smart contract artifacts",
	Long:             `Commands for generating smart contract artifacts`,
	Run: func(cmd *cobra.Command, args []string) {
		blockchain, err := ContractCmd.PersistentFlags().GetString("blockchain")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		smversion, err = ContractCmd.PersistentFlags().GetString("version")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		modelfile, err = ContractCmd.PersistentFlags().GetString("modelfile")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

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
		if target == "" {
			target = "./dovetail_generated"
		}
		target, err = filepath.Abs(target)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		generator, err := GetGenerator(blockchain)
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

// GetGenerator chooses the right generator
func GetGenerator(blockchain string) (contract.Generator, error) {
	switch strings.ToUpper(blockchain) {
	case strings.ToUpper(config.HYPERLEDGER_FABRIC):
		return createFabricGenerator()
	case strings.ToUpper(config.CORDA):
		return createCordaGenerator()
	default:
		return nil, fmt.Errorf("Unsupported blockchain to deploy '%s'", blockchain)
	}
}

func createFabricGenerator() (contract.Generator, error) {
	options := fabric.NewGenOptions(target, modelfile, smversion, enableTxnSecurity)

	fabricGen := fabric.NewGenerator(options)
	return fabricGen, nil
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

	options := corda.NewOptions(modelfile, smversion, cordaState, cmds, target, cordaNS)
	cordaGen := corda.NewGenerator(options)
	return cordaGen, nil
}

func validateModelFile(modelfile string) error {
	appConfig, err := model.ParseApp(modelfile)
	if err != nil {
		return err
	}

	if len(appConfig.Triggers) == 0 || len(appConfig.Triggers) > 1 {
		return fmt.Errorf("There must be one and only one trigger defined in smart contract application")
	}

	return nil
}
