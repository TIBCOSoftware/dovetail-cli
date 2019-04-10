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

[dependencies]
dovetail_macro_derive = { path = "{{.DovetailMacroPath}}" }
serde_json = "1.0"
serde_derive = "1.0"
serde = "1.0"
{{range $gitDependency := .GitDependencies}}{{$gitDependency.ID}} = { git = "{{$gitDependency.URL}}", branch = "{{$gitDependency.Branch}}" } {{end}}
`

// MainRs struct for main rs template
type MainRs struct {
	ModelPath              string
	GitTriggerDependencies []GitDependency
}

// MainRsTemplate main.rs template
const MainRsTemplate = `
{{range $gitTriggerDependency := .GitTriggerDependencies}}extern crate {{$gitTriggerDependency.ID}};{{end}}
`
