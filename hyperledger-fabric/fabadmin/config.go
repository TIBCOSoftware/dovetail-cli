/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

// Package fabadmin handles administration commands for hyperledger fabric
package fabadmin

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/endpoint"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/lookup"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// BackendConfig provides the definition of backend configuration
type BackendConfig struct {
	Version                string
	Client                 ClientConfig
	Channels               map[string]fab.ChannelEndpointConfig
	Organizations          map[string]fab.OrganizationConfig
	Orderers               map[string]fab.OrdererConfig
	Peers                  map[string]fab.PeerConfig
	CertificateAuthorities map[string]CAConfig
	EntityMatchers         map[string][]fab.MatchConfig
}

// ClientConfig provides the definition of the client configuration
type ClientConfig struct {
	Organization string
	CryptoConfig PathConfig
	TLSCerts     fab.ClientTLSConfig
}

// PathConfig provide the definition with a path attribute
type PathConfig struct {
	Path string
}

// EnrollCredentials holds credentials used for enrollment
type EnrollCredentials struct {
	EnrollID     string
	EnrollSecret string
}

// CAConfig defines a CA configuration
type CAConfig struct {
	URL         string
	GRPCOptions map[string]interface{}
	Registrar   EnrollCredentials
	CAName      string
	TLSCACerts  TLSCAConfig
}

// TLSCAConfig contains the TLS CA configuration
type TLSCAConfig struct {
	Path   string
	Client endpoint.TLSKeyPair
}

func fetchConfigBackend(configPath string, entityMatcherOverride string) core.ConfigProvider {
	configProvider := config.FromFile(Subst(configPath))

	if entityMatcherOverride != "" {
		return func() ([]core.ConfigBackend, error) {
			matcherProvider := config.FromFile(Subst(entityMatcherOverride))
			matcherBackends, err := matcherProvider()
			if err != nil {
				fmt.Printf("failed to read entity matcher config %s: %+v\n", entityMatcherOverride, err)
				// return the original config provider defined by configPath
				return configProvider()
			}

			currentBackends, err := configProvider()
			if err != nil {
				fmt.Printf("failed to read config %s: %+v\n", configPath, err)
				// return the entity matcher provider defined by entityMatcherOverride
				return matcherBackends, nil
			}

			// return the combined config with matcher precedency
			return append(matcherBackends, currentBackends...), nil
		}
	}
	return configProvider
}

// initBackendConfig constructs BackendConfig by extracting data from configuration files
func initBackendConfig(sdk *fabsdk.FabricSDK) (*BackendConfig, error) {
	configBackend, err := sdk.Config()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch backend config: %s", err.Error())
	}
	configLookup, ok := configBackend.(*lookup.ConfigLookup)
	if !ok {
		return nil, errors.New("failed to cast ConfigBackend to ConfigLookup")
	}

	version := configLookup.GetString("version")
	config := BackendConfig{Version: version}

	client := ClientConfig{}
	if err := configLookup.UnmarshalKey("client", &client); err != nil {
		fmt.Printf("failed to lookup client config %+v\n", err)
	} else {
		config.Client = client
	}

	channels := make(map[string]fab.ChannelEndpointConfig)
	if err := configLookup.UnmarshalKey("channels", &channels); err != nil {
		fmt.Printf("failed to lookup channels config %+v\n", err)
	} else {
		config.Channels = channels
	}

	orgs := make(map[string]fab.OrganizationConfig)
	if err := configLookup.UnmarshalKey("organizations", &orgs); err != nil {
		fmt.Printf("failed to lookup organizations config %+v\n", err)
	} else {
		config.Organizations = orgs
	}

	orderers := make(map[string]fab.OrdererConfig)
	if err := configLookup.UnmarshalKey("orderers", &orderers); err != nil {
		fmt.Printf("failed to lookup orderers config %+v\n", err)
	} else {
		config.Orderers = orderers
	}

	peers := make(map[string]fab.PeerConfig)
	if err := configLookup.UnmarshalKey("peers", &peers); err != nil {
		fmt.Printf("failed to lookup peers config %+v\n", err)
	} else {
		config.Peers = peers
	}

	cas := make(map[string]CAConfig)
	if err := configLookup.UnmarshalKey("certificateAuthorities", &cas); err != nil {
		fmt.Printf("failed to lookup CAs config %+v\n", err)
	} else {
		config.CertificateAuthorities = cas
	}

	matchers := make(map[string][]fab.MatchConfig)
	if err := configLookup.UnmarshalKey("entityMatchers", &matchers); err != nil {
		fmt.Printf("failed to lookup EntityMatchers config %+v\n", err)
	} else {
		config.EntityMatchers = matchers
	}

	return &config, nil
}

// GetOrganizationConfig returns configured attributes of a specified organization
func (bc *BackendConfig) GetOrganizationConfig(name string) *Organization {
	orgConfig, ok := bc.Organizations[name]
	if !ok {
		return nil
	}

	return &Organization{
		Name:                   name,
		MSPID:                  orgConfig.MSPID,
		CryptoPath:             orgConfig.CryptoPath,
		CertificateAuthorities: orgConfig.CertificateAuthorities,
		IsOrdererOrg:           orgConfig.Peers == nil,
	}
}

// GetOrdererName returns the name of first non-default orderer
func (bc *BackendConfig) GetOrdererName() string {
	for k := range bc.Orderers {
		if k != "_default" {
			return k
		}
	}
	return ""
}

// GetClientOrgName returns the name of the client's org
func (bc *BackendConfig) GetClientOrgName() string {
	return bc.Client.Organization
}
