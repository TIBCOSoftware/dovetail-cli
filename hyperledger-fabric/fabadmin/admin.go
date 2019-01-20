/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

// Package fabadmin handles administration commands for hyperledger fabric
package fabadmin

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
)

// ChaincodePackage contains chaincode and installation metadata
type ChaincodePackage struct {
	Name    string
	Path    string
	Version string
	Package *resource.CCPackage
}

// Organization provides the definition of an organization in the network
type Organization struct {
	Name                   string
	MSPID                  string
	CryptoPath             string
	AdminUser              string
	AdminContext           context.ClientProvider
	Peers                  []fab.Peer // discovered peers; unique id peer.URL()
	CertificateAuthorities []string
	IsOrdererOrg           bool
}

// AdminContext provides configuration of admin client
type AdminContext struct {
	sdk    *fabsdk.FabricSDK
	config *BackendConfig
}

var adminContext *AdminContext

// InitAdminContext initialize SDK based on specified config.yaml and entity-matcher file
func InitAdminContext(configPath, entityMatcherOverride string) error {
	sdk, err := fabsdk.New(fetchConfigBackend(configPath, entityMatcherOverride))
	if err != nil {
		return err
	}
	config, err := initBackendConfig(sdk)
	if err != nil {
		return err
	}

	adminContext = &AdminContext{
		sdk:    sdk,
		config: config,
	}
	return nil
}

// CloseConnection cleans up cache and closes all connections
func CloseConnection() {
	if adminContext != nil && adminContext.sdk != nil {
		adminContext.sdk.Close()
	}
}

// CreateCCPackage creates a new golang chaincode package that contains source code under ${ccRoot}/src/${ccPath}
func CreateCCPackage(ccPath, ccRoot, name, version string) (*ChaincodePackage, error) {
	pkg, err := gopackager.NewCCPackage(ccPath, Subst(ccRoot))
	if err != nil {
		return nil, err
	}
	return &ChaincodePackage{
		Name:    name,
		Path:    ccPath,
		Version: version,
		Package: pkg,
	}, nil
}

// GetClientOrgName returns the name of the configured client's org
func GetClientOrgName() string {
	return adminContext.config.GetClientOrgName()
}

// InitOrganization iniailizes admin client of a specified org
func InitOrganization(orgName, orgAdminUser string) *Organization {
	org := adminContext.config.GetOrganizationConfig(orgName)
	org.AdminUser = orgAdminUser
	org.AdminContext = adminContext.sdk.Context(fabsdk.WithUser(orgAdminUser), fabsdk.WithOrg(orgName))

	// discover local peers in the specified org context
	ctx, err := contextImpl.NewLocal(org.AdminContext)
	if err != nil {
		fmt.Printf("failed to create local context for discover peers of org %s: %+v\n", orgName, err)
		return org
	}

	if org.IsOrdererOrg {
		// skip peer discovery
		fmt.Printf("skip peer discovery for orderer org %s\n", org.Name)
		return org
	}

	// discover live peers of this organization
	discoveredPeers, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
		func() (interface{}, error) {
			peers, err := ctx.LocalDiscoveryService().GetPeers()
			if err != nil {
				return nil, errors.Wrapf(err, "error getting peers for MSP [%s]", ctx.Identifier().MSPID)
			}
			return peers, nil
		},
	)
	if err != nil {
		fmt.Printf("failed to discover peers: %+v\n", err)
		return org
	}
	org.Peers = discoveredPeers.([]fab.Peer)
	return org
}

// CreateChannel creates a specified channel. This require multiple org admins' signature;
// this proecess should be separated into multiple steps, if nobody knows signing keys for all peer orgs
func CreateChannel(channel, configRoot string, ordererOrg *Organization, peerOrgs []*Organization) (uint64, error) {
	// assuming the first non-default orderer listed in backend config belongs to ordererOrg (always true if only one orderer org)
	ordererName := adminContext.config.GetOrdererName()
	ordererClient, err := resmgmt.New(ordererOrg.AdminContext)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to create resource management client for orderer org %s", ordererOrg.MSPID)
	}

	// collect all peers' signing identies
	peerAdmins := []msp.SigningIdentity{}
	for _, org := range peerOrgs {
		mspClient, err := mspclient.New(adminContext.sdk.Context(), mspclient.WithOrg(org.Name))
		if err != nil {
			return 0, errors.Wrapf(err, "failed to create msp client for org %s", org.MSPID)
		}
		adminUser, err := mspClient.GetSigningIdentity(org.AdminUser)
		if err != nil {
			return 0, errors.Wrapf(err, "failed to get signing identity for user %s of org %s", org.AdminUser, org.MSPID)
		}
		peerAdmins = append(peerAdmins, adminUser)
	}

	// create a channel for channel tx file: ${configRoot}/${channel}.tx
	req := resmgmt.SaveChannelRequest{ChannelID: channel,
		ChannelConfigPath: path.Join(Subst(configRoot), channel+".tx"),
		SigningIdentities: peerAdmins}
	fmt.Printf("create channel %s using orderer org %s\n", channel, ordererOrg.MSPID)
	txID, err := ordererClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(ordererName))
	if err != nil {
		return 0, errors.Wrapf(err, "failed to create channel %s", channel)
	}
	fmt.Println("create channel returned tranaction id", txID)

	// wait for orderer config update (work-around fabric 1.0 bug so each block contains only 1 config update)
	queryClient, err := resmgmt.New(peerOrgs[0].AdminContext)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to create query client for org %s", peerOrgs[0].MSPID)
	}
	return waitForOrdererConfigUpdate(queryClient, channel, ordererName, true, 0)
}

// UpdateChannelAnchor updates channel config to set anchor peer for an org
func (org *Organization) UpdateChannelAnchor(channel, configRoot string, lastConfigBlock uint64) (uint64, error) {
	// assuming the first non-default orderer listed in backend config belongs to ordererOrg (always true if only one orderer org)
	ordererName := adminContext.config.GetOrdererName()

	mspClient, err := mspclient.New(adminContext.sdk.Context(), mspclient.WithOrg(org.Name))
	if err != nil {
		return 0, errors.Wrapf(err, "failed to create msp client for org %s", org.MSPID)
	}
	adminUser, err := mspClient.GetSigningIdentity(org.AdminUser)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get signing identity for user %s of org %s", org.AdminUser, org.MSPID)
	}

	orgClient, err := resmgmt.New(org.AdminContext)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to create client for org %s", org.MSPID)
	}

	// set channel anchor according to file: ${configRoot}/${channel}${orgMSPID}anchors.tx
	req := resmgmt.SaveChannelRequest{ChannelID: channel,
		ChannelConfigPath: path.Join(Subst(configRoot), channel+org.MSPID+"anchors.tx"),
		SigningIdentities: []msp.SigningIdentity{adminUser}}
	txID, err := orgClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(ordererName))
	if err != nil {
		return 0, errors.Wrapf(err, "failed to update anchor peer for channel %s", channel)
	}
	fmt.Println("update channel anchor returned tranaction id", txID)

	return waitForOrdererConfigUpdate(orgClient, channel, ordererName, false, lastConfigBlock)
}

// waitForOrdererConfigUpdate waits until the config block update has been committed.
// In Fabric 1.0 there is a bug that panics the orderer if more than one config update is added to the same block.
// This function may be invoked after each config update as a workaround.
func waitForOrdererConfigUpdate(client *resmgmt.Client, channelID, ordererName string, genesis bool, lastConfigBlock uint64) (uint64, error) {
	blockNum, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
		func() (interface{}, error) {
			chConfig, err := client.QueryConfigFromOrderer(channelID, resmgmt.WithOrdererEndpoint(ordererName))
			if err != nil {
				return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), err.Error(), nil)
			}

			currentBlock := chConfig.BlockNumber()
			if currentBlock <= lastConfigBlock && !genesis {
				return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("Block number was not incremented [%d, %d]", currentBlock, lastConfigBlock), nil)
			}
			return &currentBlock, nil
		},
	)
	return *blockNum.(*uint64), err
}

// JoinChannel makes the org join a specified channel
func (org *Organization) JoinChannel(channel string) error {
	// assuming the first non-default orderer listed in backend config belongs to ordererOrg (always true if only one orderer org)
	ordererName := adminContext.config.GetOrdererName()
	orgClient, err := resmgmt.New(org.AdminContext)
	if err != nil {
		return errors.Wrapf(err, "failed to create client for org %s", org.MSPID)
	}
	fmt.Printf("org %s joins channel %s\n", org.MSPID, channel)
	return orgClient.JoinChannel(channel, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(ordererName))
}

// IsJoinedChannel returns true if the org joined a specified channel on specified peers.  check all peers if peer is not specified.
func (org *Organization) IsJoinedChannel(channel string, peers ...fab.Peer) (bool, error) {
	targetPeers := peers
	if peers == nil {
		// check all discovered peers of the org
		targetPeers = org.Peers
	}

	orgResmg, err := resmgmt.New(org.AdminContext)
	if err != nil {
		return false, errors.Wrapf(err, "failed to create resource management client for org %s", org.MSPID)
	}

	for _, p := range targetPeers {
		resp, err := orgResmg.QueryChannels(resmgmt.WithTargets(p))
		if err != nil {
			return false, errors.Wrapf(err, "failed to query channel for org %s on peer %s", org.MSPID, p.URL())
		}

		// verify the specified channel is one of the joined channels of p
		var joined = false
		for _, chInfo := range resp.Channels {
			if chInfo.ChannelId == channel {
				joined = true
			}
		}
		if !joined {
			fmt.Printf("peer %s of org %s has not joined channel %s\n", p.URL(), org.MSPID, channel)
			return false, nil
		}
	}
	return true, nil
}

// InstallChaincode installs a chaincode package on specified peers of an org.  Install it on all peers of the org if peer URLs are not specified.
func (org *Organization) InstallChaincode(pkg *ChaincodePackage, peerURLs ...string) error {
	orgResmg, err := resmgmt.New(org.AdminContext)
	if err != nil {
		return errors.Wrapf(err, "failed to create resource management client for org %s", org.MSPID)
	}

	// pickup target peers
	var peers []fab.Peer
	var targetURLs string
	if peerURLs != nil {
		// collect peers that match the specified URLs
		targetURLs = "|" + strings.Join(peerURLs, "|") + "|"
		for _, p := range org.Peers {
			if strings.Contains(targetURLs, "|"+p.URL()+"|") {
				peers = append(peers, p)
			}
		}
	} else {
		// install on all peers discovered in the org
		peers = org.Peers
	}

	// ignore if no peer matching specified URLs
	if len(peers) == 0 {
		fmt.Printf("no peer found in org %s matching URL %s\n", org.MSPID, targetURLs)
		return nil
	}

	for _, p := range peers {
		fmt.Printf("install chaincode [%s:%s] on peer %s\n", pkg.Name, pkg.Version, p.URL())
	}

	// install chaincode on specified peers
	installCCReq := resmgmt.InstallCCRequest{Name: pkg.Name, Path: pkg.Path, Version: pkg.Version, Package: pkg.Package}
	_, err = orgResmg.InstallCC(installCCReq,
		resmgmt.WithTargets(peers...),
		resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return errors.Wrapf(err, "failed to install chaincode %s for org %s", pkg.Name, org.MSPID)
	}
	return nil
}

// convertCCArgs converts JSON {"Args": ["value1", "value2" ...]} to [][]byte
func convertCCArgs(args string) ([][]byte, error) {
	var f interface{}
	if err := json.Unmarshal([]byte(args), &f); err != nil {
		return nil, errors.Wrapf(err, "error parsing JSON string: %s", args)
	}
	m := f.(map[string]interface{})
	for k, v := range m {
		if strings.ToLower(k) == "args" {
			list, ok := v.([]interface{})
			if !ok {
				return nil, errors.Errorf("args value is not an array: %v", v)
			}
			var result [][]byte
			for _, s := range list {
				result = append(result, []byte(s.(string)))
			}
			return result, nil
		}
	}
	return nil, errors.Errorf("JSON string does not contain key name 'Args': %s", args)
}

// InstantiateChaincode instantiates a chaincode on all peers of the specified org
func (org *Organization) InstantiateChaincode(pkg *ChaincodePackage, channelID string, policy string, initArgs string) error {
	ccPolicy, err := cauthdsl.FromString(policy)
	if err != nil {
		return errors.Wrapf(err, "error creating CC policy [%s]", policy)
	}
	args, err := convertCCArgs(initArgs)
	if err != nil {
		return errors.Wrapf(err, "error converting cc init-args [%s]", initArgs)
	}

	orgResmg, err := resmgmt.New(org.AdminContext)
	if err != nil {
		return errors.Wrapf(err, "failed to create resource management client for org %s", org.MSPID)
	}

	fmt.Printf("instantiate cc %s on channel %s for org %s ...\n", pkg.Name, channelID, org.MSPID)
	_, err = orgResmg.InstantiateCC(
		channelID,
		resmgmt.InstantiateCCRequest{
			Name:       pkg.Name,
			Path:       pkg.Path,
			Version:    pkg.Version,
			Args:       args,
			Policy:     ccPolicy,
			CollConfig: nil,
		},
		resmgmt.WithTargets(org.Peers...),
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	)
	return err
}

// UpgradeChaincode instantiates a chaincode on all peers of the specified org
func (org *Organization) UpgradeChaincode(pkg *ChaincodePackage, channelID string, policy string, initArgs string) error {
	ccPolicy, err := cauthdsl.FromString(policy)
	if err != nil {
		return errors.Wrapf(err, "error creating CC policy [%s]", policy)
	}
	args, err := convertCCArgs(initArgs)
	if err != nil {
		return errors.Wrapf(err, "error converting cc init-args [%s]", initArgs)
	}

	orgResmg, err := resmgmt.New(org.AdminContext)
	if err != nil {
		return errors.Wrapf(err, "failed to create resource management client for org %s", org.MSPID)
	}

	fmt.Printf("Upgrade cc %s on channel %s for org %s\n", pkg.Name, channelID, org.MSPID)
	_, err = orgResmg.UpgradeCC(
		channelID,
		resmgmt.UpgradeCCRequest{
			Name:       pkg.Name,
			Path:       pkg.Path,
			Version:    pkg.Version,
			Args:       args,
			Policy:     ccPolicy,
			CollConfig: nil,
		},
		resmgmt.WithTargets(org.Peers...),
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	)
	return err
}

// IsCCInstalled returns true if the specified CC is already installed on all specified peers of an org. check all peers if peers are not specified
func (org *Organization) IsCCInstalled(pkg *ChaincodePackage, peers ...fab.Peer) (bool, error) {
	targetPeers := peers
	if targetPeers == nil {
		targetPeers = org.Peers
	}

	orgResmg, err := resmgmt.New(org.AdminContext)
	if err != nil {
		return false, errors.Wrapf(err, "failed to create resource management client for org %s", org.MSPID)
	}

	for _, p := range targetPeers {
		resp, err := orgResmg.QueryInstalledChaincodes(resmgmt.WithTargets(p))
		if err != nil {
			return false, errors.Wrapf(err, "failed to query peer %s of org %s", p.URL(), org.MSPID)
		}
		found := false
		for _, ccInfo := range resp.Chaincodes {
			if ccInfo.Name == pkg.Name && ccInfo.Version == pkg.Version {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("chaincode [%s:%s] is not installed on peer %s of org %s\n", pkg.Name, pkg.Version, p.URL(), org.MSPID)
			return false, nil
		}
	}
	return true, nil
}
