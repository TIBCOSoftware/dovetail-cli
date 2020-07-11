/*
* Copyright © 2020. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

package node

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/TIBCOSoftware/dovetail-cli/config"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

// Start starts the client service
func Start() {
	r := mux.NewRouter()

	// GET all releases.
	r.HandleFunc("/test", TestHandler).Methods(http.MethodGet, http.MethodOptions)

	// Use default options
	handler := cors.Default().Handler(r)

	nodePort := viper.GetString(config.NodePortKey)
	if config.IsNodeVerbose() {
		log.Printf("Server listening to port :%s", nodePort)
	}

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", nodePort), handler))
}

// TestHandler returns a list of all releases for all namespaces
func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	if r.Method == http.MethodOptions {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"status": "Connected"}`)
}
