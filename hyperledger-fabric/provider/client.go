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

/*type fixture struct {
	cryptoSuiteConfig core.CryptoSuiteConfig
	identityConfig    msp.IdentityConfig
}*/

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

/*func (f *fixture) setup(configPath string) (*fabsdk.FabricSDK, error) {

	// Instantiate the SDK
	sdk, err := NewSDK(configPath)
	if err != nil {
		panic(fmt.Sprintf("SDK init failed: %v", err))
	}

	configBackend, err := sdk.Config()
	if err != nil {
		panic(fmt.Sprintf("Failed to get config: %v", err))
	}

	f.cryptoSuiteConfig = cryptosuite.ConfigFromBackend(configBackend)
	f.identityConfig, err = mspImpl.ConfigFromBackend(configBackend)
	if err != nil {
		panic(fmt.Sprintf("Failed to get identity config: %v", err))
	}

	// TODO IMPORTANT!!!!!!!
	// Delete all private keys from the crypto suite store
	// and users from the user store
	//cleanup(f.cryptoSuiteConfig.KeyStorePath())
	//cleanup(f.identityConfig.CredentialStorePath())

	//ctxProvider := sdk.Context()
	//ctx, err := ctxProvider()
	//if err != nil {
	//	panic(fmt.Sprintf("Failed to init context: %v", err))
	//}

	// Start Http Server if it's not running
	//if !caServer.Running() {
	//	caServer.Start(lis, ctx.CryptoSuite())
	//}

	return nil, sdk
}*/
