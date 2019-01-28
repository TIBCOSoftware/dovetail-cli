/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

// Package contract implements generate and deploy chaincode for hyperledger fabric
package contract

import (
	"fmt"

	"github.com/TIBCOSoftware/dovetail-cli/hyperledger-fabric/fabadmin"
	"github.com/pkg/errors"
)

// Deployer defines the deployer attributes
type Deployer struct {
}

// Options defines the deployer options
type Options struct {
	ChaincodeID         string
	ChaincodePath       string
	ChaincodeRoot       string
	ChaincodeVersion    string
	NetworkConfigPath   string
	NetworkOverridePath string
	UserName            string
	OrgName             string
	ChaincodePolicy     string
	ChaincodeInitArgs   string
	ChannelID           string
}

// DeployOption func for each Opts argument
type DeployOption func(d *Deployer, opts *Options) error

// WithOrgName sets org name of a fabric network member
func WithOrgName(org string) DeployOption {
	return func(d *Deployer, o *Options) error {
		o.OrgName = org
		return nil
	}
}

// WithUserName sets admin user name of the org
func WithUserName(user string) DeployOption {
	return func(d *Deployer, o *Options) error {
		o.UserName = user
		return nil
	}
}

// WithChaincodeID sets the ID of the chaincode to be deployed
func WithChaincodeID(ccID string) DeployOption {
	return func(d *Deployer, o *Options) error {
		o.ChaincodeID = ccID
		return nil
	}
}

// WithChaincodeVersion sets the version of the chaincode to be deployed
func WithChaincodeVersion(ccVersion string) DeployOption {
	return func(d *Deployer, o *Options) error {
		o.ChaincodeVersion = ccVersion
		return nil
	}
}

// WithChaincodePath sets the chaincode source code location, which should be the full path of the generated source code folder
func WithChaincodePath(ccPath string) DeployOption {
	return func(d *Deployer, o *Options) error {
		tokens := fabadmin.ParseChaincodePath(ccPath)
		if len(tokens) == 2 {
			o.ChaincodeRoot = tokens[0]
			o.ChaincodePath = tokens[1]
		} else {
			o.ChaincodePath = tokens[0]
		}
		return nil
	}
}

// WithFabricConfig sets the local fabric network config yaml and entity matcher override file if necessary
func WithFabricConfig(configFiles ...string) DeployOption {
	return func(d *Deployer, o *Options) error {
		if len(configFiles) > 0 {
			o.NetworkConfigPath = configFiles[0]
		}
		if len(configFiles) > 1 {
			o.NetworkOverridePath = configFiles[1]
		}
		return nil
	}
}

// WithChaincodePolicy sets the policy for instantiating the chaincode
func WithChaincodePolicy(ccPolicy string) DeployOption {
	return func(d *Deployer, o *Options) error {
		o.ChaincodePolicy = ccPolicy
		return nil
	}
}

// WithChaincodeInitArgs sets the init-args for instantiating the chaincode
func WithChaincodeInitArgs(ccArgs string) DeployOption {
	return func(d *Deployer, o *Options) error {
		o.ChaincodeInitArgs = ccArgs
		return nil
	}
}

// WithChannelID sets the Fabric channel for instantiating the chaincode
func WithChannelID(channelID string) DeployOption {
	return func(d *Deployer, o *Options) error {
		o.ChannelID = channelID
		return nil
	}
}

func (d *Deployer) prepareDeployOpts(options ...DeployOption) (Options, error) {
	opts := Options{}
	for _, option := range options {
		err := option(d, &opts)
		if err != nil {
			return opts, errors.WithMessage(err, "failed to read deploy opts")
		}
	}

	return opts, nil
}

// Deploy deploys a smart contract for the given options
func (d *Deployer) Deploy(options ...DeployOption) error {
	fmt.Println("Deploying contract to hyperledger fabric blockchain...")
	opts, err := d.prepareDeployOpts(options...)
	if err != nil {
		return err
	}
	fmt.Printf("Using network config path '%s'\n", opts.NetworkConfigPath)

	// create cc package
	pkg, err := fabadmin.CreateCCPackage(opts.ChaincodePath, opts.ChaincodeRoot, opts.ChaincodeID, opts.ChaincodeVersion)
	if err != nil {
		return errors.Wrapf(err, "Failed to create cc package for source at %s/src/%s", opts.ChaincodeRoot, opts.ChaincodePath)
	}

	// Create SDK
	if err := fabadmin.InitAdminContext(opts.NetworkConfigPath, opts.NetworkOverridePath); err != nil {
		return errors.Wrap(err, "Failed to create new SDK")
	}

	// check if chaincode of the same ID and version already installed on the org
	orgName := opts.OrgName
	if orgName == "" {
		orgName = fabadmin.GetClientOrgName()
	}
	org := fabadmin.InitOrganization(orgName, opts.UserName)
	done, err := org.IsCCInstalled(pkg)
	if err != nil {
		fabadmin.CloseConnection()
		return errors.Wrapf(err, "Failed to check if cc already installed on org %s with user %s", orgName, opts.UserName)
	}
	if !done {
		fmt.Printf("Installing chaincode [%s:%s] on all peers of org %s ...\n", opts.ChaincodeID, opts.ChaincodeVersion, orgName)
		if err := org.InstallChaincode(pkg); err != nil {
			fabadmin.CloseConnection()
			return errors.Wrapf(err, "Failed to install cc [%s:%s] for org %s", opts.ChaincodeID, opts.ChaincodeVersion, orgName)
		}
	} else {
		fmt.Printf("Chaincode [%s:%s] already installed on peers of org %s", opts.ChaincodeID, opts.ChaincodeVersion, orgName)
	}
	fabadmin.CloseConnection()
	return nil
}

// Instantiate starts instances for a smart contract on specified channel and org
func (d *Deployer) Instantiate(options ...DeployOption) error {
	fmt.Println("Deploying contract to hyperledger fabric blockchain...")
	opts, err := d.prepareDeployOpts(options...)
	if err != nil {
		return err
	}
	fmt.Printf("Using network config path '%s'\n", opts.NetworkConfigPath)

	// create cc package
	pkg := &fabadmin.ChaincodePackage{
		Name:    opts.ChaincodeID,
		Path:    opts.ChaincodePath,
		Version: opts.ChaincodeVersion,
	}

	// Create SDK
	if err := fabadmin.InitAdminContext(opts.NetworkConfigPath, opts.NetworkOverridePath); err != nil {
		return errors.Wrap(err, "Failed to create new SDK")
	}

	// check if org has joined channel
	orgName := opts.OrgName
	if orgName == "" {
		orgName = fabadmin.GetClientOrgName()
	}
	org := fabadmin.InitOrganization(orgName, opts.UserName)
	ok, err := org.IsJoinedChannel(opts.ChannelID)
	if err != nil {
		fabadmin.CloseConnection()
		return errors.Wrapf(err, "Failed to check if org %s joined channel %s", orgName, opts.ChannelID)
	}
	if !ok {
		fabadmin.CloseConnection()
		return errors.Errorf("org %s has not joined channel %s", orgName, opts.ChannelID)
	}

	// check if chaincode of the same ID and version is installed on the org
	ok, err = org.IsCCInstalled(pkg)
	if err != nil {
		fabadmin.CloseConnection()
		return errors.Wrapf(err, "Failed to check if cc is installed on org %s with user %s", orgName, opts.UserName)
	}
	if ok {
		fmt.Printf("Instantiating chaincode [%s:%s] on all peers of org %s ...\n", opts.ChaincodeID, opts.ChaincodeVersion, orgName)
		if err := org.InstantiateChaincode(pkg, opts.ChannelID, opts.ChaincodePolicy, opts.ChaincodeInitArgs); err != nil {
			fabadmin.CloseConnection()
			return errors.Wrapf(err, "Failed to instantiate cc [%s:%s] for org %s", opts.ChaincodeID, opts.ChaincodeVersion, orgName)
		}
	} else {
		fabadmin.CloseConnection()
		return errors.Errorf("Chaincode [%s:%s] is not installed on peers of org %s", opts.ChaincodeID, opts.ChaincodeVersion, orgName)
	}
	fabadmin.CloseConnection()
	return nil
}
