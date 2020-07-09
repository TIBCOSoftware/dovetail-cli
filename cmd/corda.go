/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

// Package cmd is the one containing all the cli commands for corda operations
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/TIBCOSoftware/dovetail-cli/cmd/corda/client"
	"github.com/TIBCOSoftware/dovetail-cli/cmd/corda/contract"
	"github.com/TIBCOSoftware/dovetail-cli/cmd/corda/dapp"
)

func init() {
	cordaCmd.AddCommand(client.ClientCmd)
	cordaCmd.AddCommand(contract.ContractCmd)
	cordaCmd.AddCommand(dapp.DAppCmd)
}

// CordaCmd is the command for smart contracts
var cordaCmd = &cobra.Command{
	Use:   "corda",
	Short: "Commands for Corda apps",
	Long:  `Commands for Corda apps`,
}
