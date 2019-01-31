#!/bin/bash
# unit test for hyperledger fabric admin

echo "Running hyperledger fabric admin tests"

SDIR=$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

cd ${SDIR}/../../hyperledger-fabric/fabadmin
go test -v