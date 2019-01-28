/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package fabadmin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackendConfig(t *testing.T) {
	bc := adminContext.config
	require.NotEmpty(t, bc.Version, "config version should not be empty")
	assert.Equal(t, "org1", bc.Client.Organization, "client org should be 'org1'")
	assert.Equal(t, 2, len(bc.EntityMatchers["certificateauthority"]), "entity matcher for CA should contain 2 patterns")
	assert.Equal(t, 3, len(bc.Organizations), "configuration should contain 3 organizations")
	assert.True(t, bc.GetOrganizationConfig("ordererorg").IsOrdererOrg, "ordererorg should have IsOrdererOrg=true")
	assert.Equal(t, "orderer.example.com", bc.GetOrdererName(), "first non-default orderer should be 'orderer.example.com'")
}
