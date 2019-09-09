/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

// Package fabric is the one containing all the cli commands for hyperledger fabric target
package fabric

import (
	"github.com/spf13/cobra"
)

func init() {
	FabricCmd.PersistentFlags().StringP("version", "v", "v1.0.0", "fabric Chaincode version")
}

// FabricCmd is the command for hyperledger fabric app
var FabricCmd = &cobra.Command{
	Use:   "fabric",
	Short: "Commands for Hyperledger Fabric Apps",
	Long:  `Commands for Hyperledger Fabric Apps`,
}
