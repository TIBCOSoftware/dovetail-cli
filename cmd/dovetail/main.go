/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package main

import (
	"github.com/TIBCOSoftware/dovetail-cli/cmd/dovetail/commands"
)

func init() {
	commands.RootCmd.AddCommand(commands.VersionCmd)
}

func main() {
	commands.Execute()
}
