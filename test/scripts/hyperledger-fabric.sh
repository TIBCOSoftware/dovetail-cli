#!/bin/bash
# unit test for hyperledger-fabric

echo "Running hyperledger-fabric tests"

SDIR=$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

cd ${SDIR}

./start-fab-network.sh

./fabadmin.sh
status=$?

if [ $status -ne 0 ]; then
  echo "failed fabadmin test"
else
  ./iou.sh
  status=$?
fi

./stop-fab-network.sh
exit $status
