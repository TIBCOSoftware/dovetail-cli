/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package contract

// cargo toml template
const cargoToml = `
[package]
name = "{{.Name}}"
version = "{{.Version}}"
authors = ["user <myemail@email.com>"]
edition = "2018"
`
