#!/bin/bash
# unit test for hyperledger fabric admin

SDIR=$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

cd ${SDIR}/../../hyperledger-fabric/fabadmin
go test -v