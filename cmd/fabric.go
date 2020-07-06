/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/TIBCOSoftware/dovetail-cli/cmd/fabric/contract"
)

func init() {
	fabricCmd.AddCommand(contract.ContractCmd)
}

// FabricCmd is the command for hyperledger fabric app
var fabricCmd = &cobra.Command{
	Use:   "fabric",
	Short: "Commands for Hyperledger Fabric Apps",
	Long:  `Commands for Hyperledger Fabric Apps`,
}
