/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package provider

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// Creates a new fabsdk for the given configPath
func NewSDK(configPath string) (*fabsdk.FabricSDK, error) {
	backend, err := config.FromFile(configPath)()
	if err != nil {
		return nil, err
	}
	configProvider := func() ([]core.ConfigBackend, error) {
		return backend, nil
	}
	// Instantiate the SDK
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		return nil, err
	}
	return sdk, nil
}

func NewClientProvider(sdk *fabsdk.FabricSDK, userName, orgName string) (context.ClientProvider, error) {
	return sdk.Context(fabsdk.WithUser(userName), fabsdk.WithOrg(orgName)), nil
}
