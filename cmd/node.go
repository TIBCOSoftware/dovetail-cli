/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

package cmd

import (
	"github.com/TIBCOSoftware/dovetail-cli/cmd/node"
	"github.com/TIBCOSoftware/dovetail-cli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	port string

	// nodeCmd starts the dovetail client service
	nodeCmd = &cobra.Command{
		Use:   "node",
		Short: "Start dovetail client service",
	}
)

func init() {
	nodeCmd.AddCommand(node.StartCmd)

	nodeCmd.PersistentFlags().StringVarP(&port, config.NODE_PORT_KEY, config.NODE_PORT_SHORT_KEY, config.NODE_PORT_DEFAULT, "The port that the node is listening to")
	viper.BindPFlag(config.NODE_PORT_KEY, nodeCmd.PersistentFlags().Lookup(config.NODE_PORT_KEY))
	viper.SetDefault(config.NODE_PORT_KEY, config.NODE_PORT_DEFAULT)
}
