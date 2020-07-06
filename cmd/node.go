/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

package cmd

import (
	"github.com/TIBCOSoftware/dovetail-cli/cmd/node"
	"github.com/spf13/cobra"
)

func init() {
	nodeCmd.AddCommand(node.StartCmd)
}

var (
	// nodeCmd starts the dovetail client service
	nodeCmd = &cobra.Command{
		Use:   "node",
		Short: "Start dovetail client service",
	}
)
