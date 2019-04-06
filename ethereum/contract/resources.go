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
`

// MainRs struct for main rs template
type MainRs struct {
	ModelPath string
}

// MainRsTemplate main.rs template
const MainRsTemplate = `
extern crate dovetail_macro_derive;

use dovetail_macro_derive::JsonToDovetail;

#[derive(JsonToDovetail)]
#[path = "{{.ModelPath}}"]
struct MyFlorustApp;

pub trait HelloMacro {
    fn hello_macro();
}

fn main() {
    MyFlorustApp::hello_macro();
}

#[cfg(test)]
mod tests {

    //use hello_macro::HelloMacro;
    use dovetail_macro_derive::JsonToDovetail;

    #[derive(JsonToDovetail)]
    #[path = "{{.ModelPath}}"]
    struct MyFlorustApp;

    pub trait HelloMacro {
        fn hello_macro();
    }

    #[test]
    fn it_works() {
        MyFlorustApp::hello_macro();
        assert_eq!(2 + 2, 4);
    }
}
`
