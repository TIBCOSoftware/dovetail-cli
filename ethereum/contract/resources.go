/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package contract

// CargoToml struct for cargo toml template
type CargoToml struct {
	Name    string
	Version string
}

// CargoTomlTemplate cargo toml template
const CargoTomlTemplate = `
[package]
name = "{{.Name}}"
version = "{{.Version}}"
authors = ["user <myemail@email.com>"]
edition = "2018"
`
