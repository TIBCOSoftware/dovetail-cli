/*
* Copyright © 2020. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

package fabric

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/TIBCOSoftware/dovetail-cli/internal/fabric/channel"
	"github.com/gorilla/mux"
)

// AddHandlers adds all the handlers for hyperledger fabric
func AddHandlers(r *mux.Router) {
	// GET hyperledger fabric status.
	r.HandleFunc("/fabric/status", StatusHandler).Methods(http.MethodGet, http.MethodOptions)

	// POST hyperledger fabric channel.
	r.HandleFunc("/fabric/channel", ChannelHandler).Methods(http.MethodPost, http.MethodOptions)

}

// StatusHandler returns the status of the hyperledger fabric network
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"status": "Connected"}`)
}

// ChannelHandler handles the channel operations on the hyperledger fabric network
func ChannelHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := fmt.Errorf("Unsupported method %s", r.Method)
		errorHandler(w, r, http.StatusInternalServerError, err)
		return
	case http.MethodPost:
		// Create a new channel.
		c := &channel.Channel{}
		decoder := json.NewDecoder(r.Body)
		var payload *channel.CreatePayload
		err := decoder.Decode(&payload)
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err)
			return
		}

		req := &channel.CreateRequest{Payload: payload}
		res, err := c.Create(req)
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err)
			return
		}
		b, err := json.Marshal(res)
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	case http.MethodPut:
		err := fmt.Errorf("Unsupported method %s", r.Method)
		errorHandler(w, r, http.StatusInternalServerError, err)
		return
	case http.MethodDelete:
		err := fmt.Errorf("Unsupported method %s", r.Method)
		errorHandler(w, r, http.StatusInternalServerError, err)
		return
	default:
		err := fmt.Errorf("Unsupported method %s", r.Method)
		errorHandler(w, r, http.StatusInternalServerError, err)
		return
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, err error) {
	w.WriteHeader(status)
	fmt.Fprint(w, err.Error())
}
