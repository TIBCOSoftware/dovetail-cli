/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package contract

// CargoToml struct for cargo toml template
type CargoToml struct {
	Name              string
	Version           string
	DovetailMacroPath string
	GitDependencies   []GitDependency
}

// GitDependency struct for cargo toml template dependencies
type GitDependency struct {
	ID  string
	URL string
}

// CargoTomlTemplate cargo toml template
const CargoTomlTemplate = `
[package]
name = "{{.Name}}"
version = "{{.Version}}"
authors = ["user <myemail@email.com>"]
edition = "2018"

[lib]
name = "{{.Name}}"
path = "src/lib.rs"

[[bin]]
name = "{{.Name}}"
path = "src/main.rs"

[dependencies]
{{range $gitDependency := .GitDependencies}}{{$gitDependency.ID}} = { git = "{{$gitDependency.URL}}" }
{{end}}
`

// MainRs struct for main rs template
type MainRs struct {
	Crates  map[string]struct{}
	Uses    map[string]struct{}
	Calls   map[string]struct{}
	Derives map[string]struct{}
}

// MainRsTemplate main.rs template
const MainRsTemplate = `{{range $crate, $_ := .Crates}}extern crate {{$crate}};
{{end}}

{{range $use, $_ := .Uses}}use {{$use}};
{{end}}

{{range $derive, $_ := .Derives}}
#[{{$derive}}()];{{end}}
fn main(){
	// Calling each trigger start method
	{{range $call, $_ := .Calls}}{{$call}}();
	{{end}}
}
`

// LibRs struct for lib rs template
type LibRs struct {
	Crates  map[string]struct{}
	Uses    map[string]struct{}
	Calls   map[string]struct{}
	Derives map[string]struct{}
}

// LibRsTemplate lib.rs template
const LibRsTemplate = `{{range $crate, $_ := .Crates}}extern crate {{$crate}};
{{end}}

{{range $use, $_ := .Uses}}use {{$use}};
{{end}}
`
