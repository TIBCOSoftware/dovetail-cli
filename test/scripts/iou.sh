#!/bin/bash
# test iou build and deployment

SDIR=$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

# sample fabric config
FABRIC_TEST_CONFIG="${GOPATH}/src/github.com/hyperledger/fabric-sdk-go/test/fixtures/config"
FABRIC_CONFIG="${FABRIC_TEST_CONFIG}/config_test.yaml"
FABRIC_OVERRIDE="${FABRIC_TEST_CONFIG}/overrides/local_entity_matchers.yaml"

# install dovetail if not found in path
if [ "$(type dovetail)" == "" ]; then
  echo "dovetail not found, run install ..."
  cd ${SDIR}/../..
  make install
  PATH=${GOPATH}/bin:${PATH}
fi

# generate chaincode from dovetail model
TEST_ROOT=$(dirname "${SDIR}")
cd ${TEST_ROOT}/models/iou
if [ -f "IOU.json" ]; then
  echo "generate chaincode for IOU.json ..."
  dovetail contract generate -m IOU.json -t .
else
  echo "could not find IOU.json"
  exit 1
fi

# install chaincode on org1
CC_PATH="${TEST_ROOT}/models/iou/hlf/src/iou"
echo "install chaincode from ${CC_PATH} ..."
dovetail contract deploy --config ${FABRIC_CONFIG} --override ${FABRIC_OVERRIDE} --path ${CC_PATH} --id iou

# instantiate chaincode on org1
CC_POLICY="AND ('Org1MSP.member','Org2MSP.member')"
CHANNEL="orgchannel"
echo "instantiate chaincode from ${CC_PATH} ..."
dovetail contract instantiate --config ${FABRIC_CONFIG} --override ${FABRIC_OVERRIDE} --path "iou" --id iou --policy "${CC_POLICY}" --channel ${CHANNEL}

echo "chaincode iou instantiated on channel ${CHANNEL}"