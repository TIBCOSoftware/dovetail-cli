package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Testsimple(t *testing.T) {
	assert.True(t, true)
}

/*
const (
	configPath = "${GOPATH}/src/github.com/hyperledger/fabric-sdk-go/test/fixtures/config/config_test.yaml"
)

func TestNewSdkOk(t *testing.T) {

	sdk, err := NewSDK(configPath)

	assert.Nil(t, err, "Error should be nil for NewSDK configPath :'%s'", configPath)
	assert.NotNil(t, sdk, "SDK should not be nil for NewSDK configPath :'%s'", configPath)
}
*/
