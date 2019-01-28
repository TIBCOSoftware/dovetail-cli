/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */

// Package fabadmin handles administration commands for hyperledger fabric
package fabadmin

import (
	"bytes"
	"os"
	"strings"
)

// Subst replaces instances of '${VARNAME}' (eg ${GOPATH}) with the variable.
// Variables names that are not set by the SDK are replaced with the environment variable.
func Subst(path string) string {
	const (
		sepPrefix = "${"
		sepSuffix = "}"
	)

	splits := strings.Split(path, sepPrefix)

	var buffer bytes.Buffer

	// first split precedes the first sepPrefix so should always be written
	buffer.WriteString(splits[0]) // nolint: gas

	for _, s := range splits[1:] {
		subst, rest := substVar(s, sepPrefix, sepSuffix)
		buffer.WriteString(subst) // nolint: gas
		buffer.WriteString(rest)  // nolint: gas
	}

	return buffer.String()
}

// substVar searches for an instance of a variables name and replaces them with their value.
// The first return value is substituted portion of the string or noMatch if no replacement occurred.
// The second return value is the unconsumed portion of s.
func substVar(s string, noMatch string, sep string) (string, string) {
	endPos := strings.Index(s, sep)
	if endPos == -1 {
		return noMatch, s
	}

	v, ok := os.LookupEnv(s[:endPos])
	if !ok {
		return noMatch, s
	}

	return v, s[endPos+1:]
}

// ParseChaincodePath separates path of format xxx/src/yyy, and returns [xxx, yyy] as array of 2 elements,
// or the original string with env replaced by Subst()
func ParseChaincodePath(path string) []string {
	exPath := Subst(path)
	srcPos := strings.LastIndex(exPath, "/src/")
	if srcPos < 0 {
		if strings.HasPrefix(exPath, "src/") {
			return []string{exPath[4:]}
		}
		// no src/ in the path
		return []string{exPath}
	}
	return []string{exPath[:srcPos], exPath[srcPos+5:]}
}
