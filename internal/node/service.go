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
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// Start starts the client service
func Start() {
	r := mux.NewRouter()

	// GET all releases.
	r.HandleFunc("/test", TestHandler).Methods("GET")

	nodePort := viper.GetString(config.NodePortKey)
	if config.IsNodeVerbose() {
		log.Printf("Server listening to port :%s", nodePort)
	}

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", nodePort), r))
}

// TestHandler returns a list of all releases for all namespaces
func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}
