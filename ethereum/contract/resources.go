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
	ID     string
	URL    string
	Branch string
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
{{range $gitDependency := .GitDependencies}}{{$gitDependency.ID}} = { git = "{{$gitDependency.URL}}", branch = "{{$gitDependency.Branch}}" } {{end}}
`

// MainRs struct for main rs template
type MainRs struct {
	ModelPath              string
	GitTriggerDependencies []GitDependency
}

// MainRsTemplate main.rs template
const MainRsTemplate = `{{range $gitTriggerDependency := .GitTriggerDependencies}}extern crate {{$gitTriggerDependency.ID}};{{end}}

{{range $gitTriggerDependency := .GitTriggerDependencies}}use {{$gitTriggerDependency.ID}}::start_{{$gitTriggerDependency.ID}};{{end}}

fn main(){
	// Calling each trigger start method
	
	{{range $gitTriggerDependency := .GitTriggerDependencies}}start_{{$gitTriggerDependency.ID}}();{{end}}
}
`

// LibRs struct for lib rs template
type LibRs struct {
	ModelPath              string
	GitTriggerDependencies []GitDependency
}

// LibRsTemplate lib.rs template
const LibRsTemplate = `{{range $gitTriggerDependency := .GitTriggerDependencies}}extern crate {{$gitTriggerDependency.ID}};{{end}}

{{range $gitTriggerDependency := .GitTriggerDependencies}}use {{$gitTriggerDependency.ID}}::start_{{$gitTriggerDependency.ID}};{{end}}
`
