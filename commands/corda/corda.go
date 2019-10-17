/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

// Package corda is the one containing all the cli commands for corda operations
package corda

import (
	"github.com/spf13/cobra"

	"github.com/TIBCOSoftware/dovetail-cli/commands/corda/client"
	"github.com/TIBCOSoftware/dovetail-cli/commands/corda/contract"
	"github.com/TIBCOSoftware/dovetail-cli/commands/corda/dapp"
)

func init() {
	CordaCmd.AddCommand(client.ClientCmd)
	CordaCmd.AddCommand(contract.ContractCmd)
	CordaCmd.AddCommand(dapp.DAppCmd)
}

// CordaCmd is the command for smart contracts
var CordaCmd = &cobra.Command{
	Use:   "corda",
	Short: "Commands for Corda apps",
	Long:  `Commands for Corda apps`,
}
