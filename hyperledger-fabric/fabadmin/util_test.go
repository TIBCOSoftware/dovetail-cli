/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package fabadmin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubst(t *testing.T) {
	os.Setenv("FOO", "fvar")
	os.Setenv("BAR", "bvar")
	path := "${FOO}/path/${BAR}/sep/${BAZ}/suff"
	assert.Equal(t, "fvar/path/bvar/sep/${BAZ}/suff", Subst(path), "Subst does not work")
}

func TestParseChaincodePath(t *testing.T) {
	os.Setenv("FOO", "fvar")
	os.Setenv("BAR", "bvar")
	tokens := ParseChaincodePath("xxx/${FOO}/src/${BAR}/yyy")
	assert.Equal(t, 2, len(tokens), "should return 2 array elements")
	assert.Equal(t, "xxx/fvar", tokens[0], "first element value does not match")
	assert.Equal(t, "bvar/yyy", tokens[1], "second element value does not match")

	tokens = ParseChaincodePath("src/${BAR}/yyy")
	assert.Equal(t, 1, len(tokens), "should return 1 array element")
	assert.Equal(t, "bvar/yyy", tokens[0], "first element value does not match")

	tokens = ParseChaincodePath("${FOO}/${BAR}/yyy")
	assert.Equal(t, 1, len(tokens), "should return 1 array elements")
	assert.Equal(t, "fvar/bvar/yyy", tokens[0], "first element value does not match")

	tokens = ParseChaincodePath("xxx/${FOO}/src/")
	assert.Equal(t, 2, len(tokens), "should return 2 array elements")
	assert.Empty(t, tokens[1], "second element should be empty")
}
