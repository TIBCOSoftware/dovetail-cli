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
If you see test errors, please refer the Troubleshooting or Support section below.

For step by step instructions on how to setup Project Dovetail™ environment, please go to the installation section on the [documentation page](https://tibcosoftware.github.io/dovetail/getting-started/getting-started-cli/)

### Note on third party dependencies

Once you install the cli, [these](./go.sum) third party dependencies will be downloaded to your machine. Please note that these third party dependencies are subject to their own license terms.

### Contributing

New contributions are welcome. If you would like to submit one, follow the instructions in the contributions section on the [documentation page](https://tibcosoftware.github.io/dovetail/contributing/contributing/)

## License
dovetail-cli is licensed under a BSD-type license. See [LICENSE](https://github.com/TIBCOSoftware/dovetail-cli/blob/master/LICENSE) for license text.

### Support
For Q&A you can contact us at tibcolabs@tibco.com.

## Troubleshooting

### Fabric admin test fails on Ubuntu

The current version of Fabric SDK supports Go 1.11.0-1.11.4. Thus, if the installation failed to download Go dependencies for Fabric SDK, you will need to download Go 1.11.4 and change the `$GOROOT` and `$PATH` environment variables to point to this version.

If the `fabric admin tests` failed with the following error:
```
hyperledger/fabric/core/operations/system.go:227:23: not enough arguments in call to s.statsd.SendLoop
```
You may resolve dependency issues for the Fabric SDK as follows.
```
cd ${GOPATH}/src/github.com/hyperledger/fabric-sdk-go
make depend
```
Edit the `Makefile` to turn off `gometalinter` for the target `.PHONY: unit-test`, i.e., in the command under this target, update the variable to use `TEST_WITH_LINTER=false`, and then execute the unit tests
```
make unit-test
```
You may also run the integration tests of the Fabric SDK to make sure that all dependencies are updated correctly, i.e.,
```
make integration-tests-stable-local
```
If the Fabric SDK tests complete successfully, you can clean up the docker containers from the tests as follows:
```
docker kill $(docker ps | egrep "fabsdkgo|hyperledger" | awk '{print $1}')
docker rm $(docker ps -a | egrep "fabsdkgo|hyperledger" | awk '{print $1}')
docker rmi $(docker images | grep fabsdkgo | awk '{print $3}')
```
You can then try to build and test the dovetail-cli again, i.e.,
```
cd ${GOPATH}/src/github.com/TIBCOSoftware/dovetail-cli
make
```