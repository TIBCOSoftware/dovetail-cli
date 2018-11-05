/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package contract

import (
	"fmt"

	"github.com/TIBCOSoftware/dovetail-cli/hyperledger-fabric/provider"
	"github.com/TIBCOSoftware/dovetail-cli/pkg/contract"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
)

// Deployer defines the deployer attributes
type Deployer struct {
	Opts *Options
}

// Options defines the deployer options
type Options struct {
	ChaincodeID       string
	ChaincodePath     string
	ChaincodeVersion  string
	NetworkConfigPath string
	UserName          string
	OrgName           string
}

// NewDeployer is the deployer constructor
func NewDeployer(opts *Options) contract.Deployer {
	return &Deployer{Opts: opts}
}

// NewOptions is the options constructor
func NewOptions(chaincodeID, chaincodePath, chaincodeVersion, networkConfigPath, userName, orgName string) *Options {
	return &Options{ChaincodeID: chaincodeID, ChaincodePath: chaincodePath, ChaincodeVersion: chaincodeVersion, NetworkConfigPath: networkConfigPath, UserName: userName, OrgName: orgName}
}

// Deploy deploys a smart contract for the given options
func (d *Deployer) Deploy() error {
	fmt.Println("Deploying contract to hyperledger fabric blockchain...")
	fmt.Printf("Using network config path '%s'\n", d.Opts.NetworkConfigPath)

	// Create SDK
	sdk, err := provider.NewSDK(d.Opts.NetworkConfigPath)
	if err != nil {
		return err
	}

	// Create Client Provider
	cp, err := provider.NewClientProvider(sdk, d.Opts.UserName, d.Opts.OrgName)
	if err != nil {
		return err
	}

	// Create Client
	c, err := resmgmt.New(cp)
	if err != nil {
		return err
	}

	ccPkg, err := packager.NewCCPackage(d.Opts.ChaincodePath, "")
	if err != nil {
		return err
	}

	// Create Install Request
	req := resmgmt.InstallCCRequest{Name: d.Opts.ChaincodeID, Version: d.Opts.ChaincodeVersion, Path: d.Opts.ChaincodePath, Package: ccPkg}

	// Get the targets
	targets, err := getTargets()
	if err != nil {
		return err
	}

	responses, err := c.InstallCC(req, targets...)
	if err != nil {
		fmt.Printf("failed to install chaincode: %v", err)
	}

	if len(responses) > 0 {
		for _, response := range responses {
			fmt.Printf("Chaincode install response from target: '%s', status: '%d', info: '%s' \n", response.Target, response.Status, response.Info)
		}
		return nil
	}
	return nil
}

// TODO Get the targets
func getTargets() ([]resmgmt.RequestOption, error) {
	/*
		// Get backendConfig
		configBackend, err := sdk.Config()
		if err != nil {
			return err
		}

		// Get endpoint config
		endpointConfig, err := fab.ConfigFromBackend(configBackend)
		if err != nil {
			return err
		}
			//peerConfig, err := endpointConfig.PeerConfig("peer0.org1.example.com")
			//if err != nil {
			//	return err
			//}

			//opt := peer.FromPeerConfig(peerConfig)


			var targetPeer *peer.Peer

			networkPeers, err := endpointConfig.NetworkPeers()
			// Iterate through network peers
			for _, networkPeer := range networkPeers {
				fmt.Printf("Network Peer URL : %s", networkPeer.URL)
				if networkPeer.URL == "peer0.org1.example.com:7051" || networkPeer.URL == "localhost:7051" {
					fmt.Printf("Creating peer")
					// Get option
					opt := peer.FromPeerConfig(&networkPeer)
					targetPeer, err = peer.New(endpointConfig, opt)
					if err != nil {
						return err
					}
					break
				}
			}
			_ = resmgmt.WithTargets(targetPeer)
			fmt.Printf("Peer found with url %s and MSID %s\n", targetPeer.URL(), targetPeer.MSPID())

	*/
	return []resmgmt.RequestOption{}, nil
}
