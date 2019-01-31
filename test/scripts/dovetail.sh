#!/bin/bash
# unit test for dovetail

echo "Running dovetail tests"

SDIR=$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

cd ${SDIR}/../../model
go test -v