#!/bin/bash

echo "Installing dovetail-cli..."
go get -u -d github.com/TIBCOSoftware/dovetail-cli
cd ${GOPATH}/src/github.com/TIBCOSoftware/dovetail-cli
make install