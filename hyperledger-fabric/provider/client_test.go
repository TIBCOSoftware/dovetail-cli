package provider

import (
	"testing"

	"github.com/pkg/errors"
)

func TestSdkSetupOk(t *testing.T) {
	// TODO get config from config
	configPath := "/Users/mtorres/dev/TIBCOSoftware/honeycomb-cli-dev/src/github.com/TIBCOSoftware/honeycomb/waggle/fixtures/config/config_test.yaml"
	f := fixture{}
	sdk := f.setup(configPath)

	clientConfig, err := sdk.provider.IdentityConfig().Client()
	if err != nil {
		return nil, errors.WithMessage(err, "retrieving client configuration failed")
	}

	//assert.NotNil(t, cfg, "Context should not be nil")

	// Check identifier
	//assert.NotNil(t, ctx.Identifier(), "Identifier should not be nil")
}

/*
func TestClientProviderOk(t *testing.T) {
	// TODO get config from config
	ctxProvider := provider.NewClientProvider("/Users/mtorres/dev/TIBCOSoftware/honeycomb-cli-dev/src/github.com/TIBCOSoftware/honeycomb/waggle/fixtures/config/config_test.yaml")
	ctx, err := ctxProvider()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, ctx, "Context should not be nil")

	// Check identifier
	assert.NotNil(t, ctx.Identifier(), "Identifier should not be nil")
}
*/
