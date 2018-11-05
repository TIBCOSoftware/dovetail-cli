/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package contract

import (
	"fmt"
	"os"
	"strings"

	"github.com/TIBCOSoftware/dovetail-cli/config"
	fabricc "github.com/TIBCOSoftware/dovetail-cli/hyperledger-fabric/contract"
	"github.com/TIBCOSoftware/dovetail-cli/pkg/contract"
	"github.com/spf13/cobra"
)

var (
	id       string
	path     string
	version  string
	userName string
	orgName  string
)

func init() {
	//TODO: to be supported
	//contractCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringVar(&id, "id", "", "Id of the Smart Contract")
	deployCmd.Flags().StringVar(&path, "path", "", "Path of the Smart Contract")
	deployCmd.Flags().StringVarP(&version, "version", "v", "", "Version of the Smart Contract")
	deployCmd.Flags().StringVar(&userName, "user", "Admin", "User Name")
	deployCmd.Flags().StringVar(&orgName, "org", "org1", "Organization Name")

	// Required flags
	deployCmd.MarkFlagRequired("id")
	deployCmd.MarkFlagRequired("path")
	deployCmd.MarkFlagRequired("version")
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the smart contract to the chosen blockchain",
	Long:  `Deploy the smart contract to the chosen blockchain`,
	Run: func(cmd *cobra.Command, args []string) {
		blockchain, err := contractCmd.PersistentFlags().GetString("blockchain")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		modelFile, err := contractCmd.PersistentFlags().GetString("modelfile")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		deployer, err := GetDeployer(blockchain, modelFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := deployer.Deploy(); err != nil {
			fmt.Printf("Error deploying the contract: '%s'", err)
			os.Exit(1)
		}
	},
}

// GetDeployer chooses the right deployer
func GetDeployer(blockchain, modelFile string) (contract.Deployer, error) {
	switch strings.ToUpper(blockchain) {
	case strings.ToUpper(config.HYPERLEDGER_FABRIC):
		fabDep, err := createFabricDeployer(modelFile)
		if err != nil {
			return nil, err
		}
		return fabDep, nil
	default:
		return nil, fmt.Errorf("Unsupported blockchain to deploy '%s'", blockchain)
	}
}

func createFabricDeployer(modelFile string) (contract.Deployer, error) {
	options := fabricc.NewOptions(id, path, version, modelFile, userName, orgName)
	fabDep := fabricc.NewDeployer(options)
	return fabDep, nil
}
