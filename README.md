# Dovetail cli

[![Build Status](https://travis-ci.org/TIBCOSoftware/dovetail-cli.svg?branch=master)](https://travis-ci.org/TIBCOSoftware/dovetail-cli)

## What is dovetail cli?

Dovetail cli is a command line tool that enables smart contract generation for different blockchain technologies from the same model.

This allows definition of your smart contracts on a model driven approach and abstraction of your smart contract logic from the low level code, improving visibility, auditability and reduce errors.

## Installation

Install [Go version 1.11.x](https://golang.org/doc/install) and [set GOPATH environment variable](https://golang.org/doc/code.html#GOPATH).  Clone this project, then install and test it as follows:
```
export PATH=${GOPATH}/bin:${PATH}
go get -u -d github.com/TIBCOSoftware/dovetail-cli
cd ${GOPATH}/src/github.com/TIBCOSoftware/dovetail-cli
make
```

For step by step instructions on how to setup Project Dovetail™ environment, please go to the installation section on the [documentation page](https://tibcosoftware.github.io/dovetail/getting-started/getting-started-cli/)

### Note on third party dependencies

Once you install the cli, [these](./go.sum) third party dependencies will be downloaded to your machine. Please note that these third party dependencies are subject to their own license terms.

### Contributing

New contributions are welcome. If you would like to submit one, follow the instructions in the contributions section on the [documentation page](https://tibcosoftware.github.io/dovetail/contributing/contributing/)

## License
dovetail-cli is licensed under a BSD-type license. See [LICENSE](https://github.com/TIBCOSoftware/dovetail-cli/blob/master/LICENSE) for license text.

### Support
For Q&A you can contact us at tibcolabs@tibco.com.
