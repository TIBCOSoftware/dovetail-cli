/*
* Copyright © 2020. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

package node

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TIBCOSoftware/dovetail-cli/config"
	"github.com/TIBCOSoftware/dovetail-cli/internal/fabric"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

// Start starts the client service
func Start() {
	r := mux.NewRouter()

	// Add fabric handlers
	fabric.AddHandlers(r)

	// Use default options
	handler := cors.AllowAll().Handler(r)

	nodePort := viper.GetString(config.NodePortKey)
	if config.IsNodeVerbose() {
		log.Printf("Server listening to port :%s", nodePort)
	}

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", nodePort), handler))
}
