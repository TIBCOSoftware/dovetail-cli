/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

// Package contract is the one containing all the cli commands for contract operations
package contract

import (
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/dovetail-cli/config"
	"github.com/spf13/cobra"
)

func init() {
	ContractCmd.PersistentFlags().StringP("blockchain", "b", config.HYPERLEDGER_FABRIC, fmt.Sprintf("Target blockchain to deploy to (%s)", strings.Join(config.Blockchains(), "|")))
	ContractCmd.PersistentFlags().StringP("modelfile", "m", "", "Smart contract flow model file")
	ContractCmd.PersistentFlags().StringP("version", "v", "1.0", "Smart contract version")

	// Required flags
	ContractCmd.MarkPersistentFlagRequired("modelfile")
}


// ContractCmd is the command for smart contracts
var ContractCmd = &cobra.Command{
	Use:              "contract",
	Short:            "Commands for Smart Contracts",
	Long:             `Commands for Smart Contracts`,
}
