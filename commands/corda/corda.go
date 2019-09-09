/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

// Package corda is the one containing all the cli commands for corda operations
package corda

import (
	"github.com/spf13/cobra"
)

func init() {
	CordaCmd.PersistentFlags().StringP("version", "v", "v1.0.0", "Contract version")
}

// CordaCmd is the command for smart contracts
var CordaCmd = &cobra.Command{
	Use:   "corda",
	Short: "Commands for Corda apps",
	Long:  `Commands for Corda apps`,
}
