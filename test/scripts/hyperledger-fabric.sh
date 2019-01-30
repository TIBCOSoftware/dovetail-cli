#!/bin/bash
# unit test for hyperledger-fabric

echo "Running hyperledger-fabric tests"

SDIR=$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

cd ${SDIR}/../../hyperledger-fabric/contract
go test -v

cd ${SDIR}/../../hyperledger-fabric/provider
go test -v