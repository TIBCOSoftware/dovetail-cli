#!/bin/bash
# start fabric network if necessary

GOPATH="${GOPATH:-${HOME}/go}"
FABRIC_SDK_PATH="${GOPATH}/src/github.com/hyperledger/fabric-sdk-go"

# verify that 4 peers are running
peers=$(docker ps | grep "hyperledger/fabric-peer" | wc -l)
if [ "${peers}" -ne 4 ]; then
    echo "not found 4 running peers, so start the sdk test network ..."
    cd ${FABRIC_SDK_PATH}
    make dockerenv-stable-up &
    sleep 10
    peers=$(docker ps | grep "hyperledger/fabric-peer" | wc -l)
    while [ "${peers}" -ne 4 ]; do
        echo "waiting for fabric network up ..."
        sleep 10
        peers=$(docker ps | grep "hyperledger/fabric-peer" | wc -l)
    done
    echo "fabric network is up"
else
    echo "found 4 peers already running"
fi