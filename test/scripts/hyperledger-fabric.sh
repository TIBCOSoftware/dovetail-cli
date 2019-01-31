#!/bin/bash
# unit test for hyperledger-fabric

echo "Running hyperledger-fabric tests"

SDIR=$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

cd ${SDIR}

./start-fab-network.sh

./fabadmin.sh

./iou.sh

./stop-fab-network.sh