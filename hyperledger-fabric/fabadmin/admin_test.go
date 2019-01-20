/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package fabadmin

import (
	"fmt"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	TestRunID = "0"
)

const (
	configPath         = "${GOPATH}/src/github.com/hyperledger/fabric-sdk-go/test/fixtures/config/config_test.yaml"
	entityMatcherLocal = "${GOPATH}/src/github.com/hyperledger/fabric-sdk-go/test//fixtures/config/overrides/local_entity_matchers.yaml"
	ccPath             = "github.com/example_cc"
	ccRoot             = "${GOPATH}/src/github.com/hyperledger/fabric-sdk-go/test/fixtures/testdata"
	ccPolicy           = "AND ('Org1MSP.member','Org2MSP.member')"
	ccArgs             = "{\"Args\": [\"init\", \"a\", \"100\", \"b\", \"200\"]}"
	channelID          = "orgchannel"
	channelConfigPath  = "${GOPATH}/src/github.com/hyperledger/fabric-sdk-go/test/fixtures/fabric/v1.4/channel"
)

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		panic(fmt.Sprintf("unable to setup %+v", err))
	}
	r := m.Run()
	teardown()
	os.Exit(r)
}

func setup() error {
	// Create SDK setup for the tests
	if err := InitAdminContext(configPath, entityMatcherLocal); err != nil {
		return errors.Wrap(err, "Failed to create new SDK")
	}
	return nil
}

func teardown() {
	CloseConnection()
}

func TestOrganization(t *testing.T) {
	org := InitOrganization("org1", "Admin")
	assert.Equal(t, "Org1MSP", org.MSPID, "org1 MSP should be 'Org1MSP'")
	assert.False(t, org.IsOrdererOrg, "org1 should not be an orderer org")
	assert.Equal(t, 2, len(org.Peers), "org1 should have 2 peers")

	org = InitOrganization("org2", "Admin")
	assert.Equal(t, "Org2MSP", org.MSPID, "org2 MSP should be 'Org2MSP'")
	assert.False(t, org.IsOrdererOrg, "org2 should not be an orderer org")
	assert.Equal(t, 2, len(org.Peers), "org2 should have 2 peers")
}

func TestChannel(t *testing.T) {
	// check if org1 already joined channel
	org1 := InitOrganization("org1", "Admin")
	ok, err := org1.IsJoinedChannel(channelID)
	require.NoError(t, err, "failed to check if org %s joined channel %s", org1.MSPID, channelID)
	if ok {
		fmt.Println("skip channel setup: org1 already joined channel")
		return
	}

	// create new channel
	org2 := InitOrganization("org2", "Admin")
	ordererorg := InitOrganization("ordererorg", "Admin")
	lastConfigBlock, err := CreateChannel(channelID, channelConfigPath, ordererorg, []*Organization{org1, org2})
	require.NoError(t, err, "failed to create channel %s", channelID)

	// update anchor peer for orgs
	lastConfigBlock, err = org1.UpdateChannelAnchor(channelID, channelConfigPath, lastConfigBlock)
	require.NoError(t, err, "failed to update anchor peer of channel %s for org %s", channelID, org1.MSPID)
	assert.NotZero(t, lastConfigBlock, "channel update block number for org1 should not be zero")
	lastConfigBlock, err = org2.UpdateChannelAnchor(channelID, channelConfigPath, lastConfigBlock)
	require.NoError(t, err, "failed to update anchor peer of channel %s for org %s", channelID, org2.MSPID)
	assert.NotZero(t, lastConfigBlock, "channel update block number for org2 should not be zero")

	// orgs to join channel
	err = org1.JoinChannel(channelID)
	require.NoError(t, err, "failed for org %s to join channel %s", org1.MSPID, channelID)
	err = org2.JoinChannel(channelID)
	require.NoError(t, err, "failed for org %s to join channel %s", org2.MSPID, channelID)
}

func TestChaincode(t *testing.T) {
	// package example chaincode
	ccName := "example_cc_" + TestRunID
	pkg, err := CreateCCPackage(ccPath, ccRoot, ccName, "1.0")
	require.NoError(t, err, "failed to create cc package")

	// check if cc already installed on org1
	org1 := InitOrganization("org1", "Admin")
	ok, err := org1.IsCCInstalled(pkg)
	require.NoError(t, err, "failed to check if cc package is installed on org %s", org1.MSPID)
	if ok {
		fmt.Printf("skip cc install: chain code [%s:%s] already installed on all peers of org %s\n", pkg.Name, pkg.Version, org1.MSPID)
		return
	}

	// install chaincode on both peer orgs
	err = org1.InstallChaincode(pkg)
	require.NoError(t, err, "failed to install chaincode %s on org %s", ccName, org1.MSPID)
	org2 := InitOrganization("org2", "Admin")
	err = org2.InstallChaincode(pkg)
	require.NoError(t, err, "failed to install chaincode %s on org %s", ccName, org2.MSPID)

	// instantiate chaincode on org1
	err = org1.InstantiateChaincode(pkg, channelID, ccPolicy, ccArgs)
	require.NoError(t, err, "failed to instantiate chaincode %s on org %s", ccName, org2.MSPID)
}
