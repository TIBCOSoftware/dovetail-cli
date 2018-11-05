/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package contract

import (
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/dovetail-cli/cmd"
	"github.com/TIBCOSoftware/dovetail-cli/config"
	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(contractCmd)
	contractCmd.PersistentFlags().StringP("blockchain", "b", config.HYPERLEDGER_FABRIC, fmt.Sprintf("Target blockchain to deploy to (%s)", strings.Join(config.Blockchains(), "|")))
	contractCmd.PersistentFlags().StringP("modelfile", "m", "", "Smart contract flow model file")
	contractCmd.PersistentFlags().StringP("version", "v", "1.0", "Smart contract version")

	// Required flags
	contractCmd.MarkFlagRequired("modelfile")
}

var contractCmd = &cobra.Command{
	Use:              "contract",
	Short:            "Commands for Smart Contracts",
	Long:             `Commands for Smart Contracts`,
	PersistentPreRun: cmd.ConfigPreRun,
}
