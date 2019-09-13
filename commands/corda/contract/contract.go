/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

// Package contract is the one containing all the cli commands for contract operations
package contract

import (
	"github.com/spf13/cobra"
)

func init() {
	ContractCmd.PersistentFlags().StringP("version", "v", "v1.0.0", "Contract version")
}

// ContractCmd is the command for smart contracts
var ContractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Commands for Smart Contracts",
	Long:  `Commands for Smart Contracts`,
}
