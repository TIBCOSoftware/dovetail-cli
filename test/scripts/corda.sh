#!/bin/bash
# unit test for corda

echo "Running corda tests"

SDIR=$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

cd ${SDIR}/../../corda
go test -v ./...