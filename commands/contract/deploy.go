/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

// Package contract is the one containing all the cli commands for contract operations
package contract

import (
	"fmt"
	"os"
	"strings"

	"github.com/TIBCOSoftware/dovetail-cli/config"
	fabc "github.com/TIBCOSoftware/dovetail-cli/hyperledger-fabric/contract"
	"github.com/spf13/cobra"
)

var (
	id              string
	path            string
	policy          string
	initArgs        string
	channel         string
	userName        string
	orgName         string
	networkConfig   string
	networkOverride string
)

func init() {
	ContractCmd.AddCommand(deployCmd)
	ContractCmd.AddCommand(instantiateCmd)

	deployCmd.Flags().StringVar(&id, "id", "", "Id of the Smart Contract")
	deployCmd.Flags().StringVar(&path, "path", "", "Source folder of the generated Smart Contract, e.g., /path/to/hlf/src/myapp")
	deployCmd.Flags().StringVar(&userName, "user", "Admin", "Admin user Name of the org")
	deployCmd.Flags().StringVar(&orgName, "org", "", "Organization Name, defaut to client org in network config")
	deployCmd.Flags().StringVar(&networkConfig, "config", "", "Path of base Fabric network config file, e.g., /path/to/config.yaml")
	deployCmd.Flags().StringVar(&networkOverride, "override", "", "Path of Fabric network override config file")

	// Required flags
	deployCmd.MarkFlagRequired("id")
	deployCmd.MarkFlagRequired("path")
	deployCmd.MarkFlagRequired("config")

	instantiateCmd.Flags().StringVar(&id, "id", "", "Id of the Smart Contract")
	instantiateCmd.Flags().StringVar(&path, "path", "", "Path of the deployed Smart Contract, e.g., myapp")
	instantiateCmd.Flags().StringVar(&policy, "policy", "", "Endorsement policy to instantiate engines for the Smart Contract")
	instantiateCmd.Flags().StringVar(&initArgs, "init", "{\"Args\": [\"init\"]}", "init args to instantiate engines for the Smart Contract")
	instantiateCmd.Flags().StringVar(&channel, "channel", "", "channel ID to instantiate engines for the Smart Contract")
	instantiateCmd.Flags().StringVar(&userName, "user", "Admin", "Admin user Name of the org")
	instantiateCmd.Flags().StringVar(&orgName, "org", "", "Organization Name, defaut to client org in network config")
	instantiateCmd.Flags().StringVar(&networkConfig, "config", "", "Path of base Fabric network config file, e.g., /path/to/config.yaml")
	instantiateCmd.Flags().StringVar(&networkOverride, "override", "", "Path of Fabric network override config file")

	// Required flags
	instantiateCmd.MarkFlagRequired("id")
	instantiateCmd.MarkFlagRequired("path")
	instantiateCmd.MarkFlagRequired("policy")
	instantiateCmd.MarkFlagRequired("channel")
	instantiateCmd.MarkFlagRequired("config")
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the smart contract to the chosen blockchain",
	Long:  `Deploy the smart contract to the chosen blockchain`,
	Run: func(cmd *cobra.Command, args []string) {
		blockchain, err := ContractCmd.PersistentFlags().GetString("blockchain")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		version, err := ContractCmd.PersistentFlags().GetString("version")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		switch strings.ToUpper(blockchain) {
		case strings.ToUpper(config.HYPERLEDGER_FABRIC):
			deployer := &fabc.Deployer{}
			err := deployer.Deploy(
				fabc.WithChaincodeID(id),
				fabc.WithChaincodePath(path),
				fabc.WithChaincodeVersion(version),
				fabc.WithUserName(userName),
				fabc.WithOrgName(orgName),
				fabc.WithFabricConfig(networkConfig, networkOverride),
			)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		default:
			fmt.Println("Unsupported blockchain for deploy:", blockchain)
		}
	},
}

var instantiateCmd = &cobra.Command{
	Use:   "instantiate",
	Short: "Start instances for the smart contract on the chosen blockchain",
	Long:  `Start instances for the smart contract to the chosen blockchain`,
	Run: func(cmd *cobra.Command, args []string) {
		blockchain, err := ContractCmd.PersistentFlags().GetString("blockchain")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		version, err := ContractCmd.PersistentFlags().GetString("version")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		switch strings.ToUpper(blockchain) {
		case strings.ToUpper(config.HYPERLEDGER_FABRIC):
			deployer := &fabc.Deployer{}
			err := deployer.Instantiate(
				fabc.WithChaincodeID(id),
				fabc.WithChaincodeVersion(version),
				fabc.WithChaincodePath(path),
				fabc.WithChaincodeInitArgs(initArgs),
				fabc.WithChannelID(channel),
				fabc.WithChaincodePolicy(policy),
				fabc.WithUserName(userName),
				fabc.WithOrgName(orgName),
				fabc.WithFabricConfig(networkConfig, networkOverride),
			)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		default:
			fmt.Println("Unsupported blockchain for instantiate:", blockchain)
		}
	},
}
