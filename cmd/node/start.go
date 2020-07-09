/*
* Copyright © 2020. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

package node

import (
	int_node "github.com/TIBCOSoftware/dovetail-cli/internal/node"
	"github.com/spf13/cobra"
)

var (
	// StartCmd starts the client service
	StartCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts the client service",
		Run:   start,
	}
)

// start starts the client service
func start(cmd *cobra.Command, args []string) {
	int_node.Start()
}
