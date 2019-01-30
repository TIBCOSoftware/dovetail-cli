#!/bin/bash
# unit test for dovetail

SDIR=$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

cd ${SDIR}/../../corda
go test -v ./...