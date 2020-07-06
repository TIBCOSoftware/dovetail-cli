/*
* Copyright © 2020. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

package node

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var (
	// StartCmd starts the client service
	StartCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts the client service",
		Run:   start,
	}
)

// start starts the client service
func start(cmd *cobra.Command, args []string) {
	r := mux.NewRouter()

	// GET all releases.
	r.HandleFunc("/test", TestHandler).Methods("GET")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":9999", r))
}

// TestHandler returns a list of all releases for all namespaces
func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}
