/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

// Package client is the one containing all the cli commands for dapp operations
package client

import (
	"github.com/spf13/cobra"
)

func init() {
	ClientCmd.PersistentFlags().StringP("version", "v", "v1.0.0", "version for generic client. for CorDApp client, the version should match the CorDapp version")
}

// ClientCmd is the command for client app
var ClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Commands for Client Apps",
	Long:  `Commands for Client Apps`,
}
