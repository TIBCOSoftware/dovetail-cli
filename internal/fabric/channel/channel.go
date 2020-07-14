/*
* Copyright © 2020. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

package channel

import (
	//"encoding/json"
	//"fmt"
	"log"
)

//Channel is the channel
type Channel struct {
}

//CreateRequest is the request input data
type CreateRequest struct {
	Payload *CreatePayload `json:"payload"`
}

//CreatePayload is the payload of the request input data
type CreatePayload struct {
	Name string `json:"name"`
}

//CreateResponse is the response data
type CreateResponse struct {
	Status string `json:"status"`
}

//Create creates the channel
func (c *Channel) Create(request *CreateRequest) (*CreateResponse, error) {
	log.Printf("Creating channel...")

	/*err := checkInstallPayload(request.Payload)
	if err != nil {
		log.Printf("%+v", err)
		return nil, err
	}

	settings := cli.New()

	actionConfig := new(action.Configuration)

	namespace := settings.Namespace()
	if request.Payload.Namespace != nil {
		namespace = *request.Payload.Namespace
	}

	// You can pass an empty string instead of settings.Namespace() to list
	// all namespaces
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Printf("%+v", err)
		return nil, err
	}

	client := action.NewInstall(actionConfig)

	client.ReleaseName = request.Payload.ID

	chart := filepath.Join(os.Getenv("CHARTS_PATH"), fmt.Sprintf("%s-%s.tgz", request.Payload.Name, request.Payload.Version))

	cp, err := client.ChartPathOptions.LocateChart(chart, settings)
	if err != nil {
		log.Printf("%+v", err)
		return nil, err
	}

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		log.Printf("%+v", err)
		return nil, err
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		log.Printf("%+v", err)
		return nil, err
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			log.Printf("%+v", err)
			return nil, err
		}
	}

	client.Namespace = settings.Namespace()

	release, err := client.Run(chartRequested, request.Payload.Values)
	if err != nil {
		log.Printf("%+v", err)
		return nil, err
	}

	response := newReleaseInstallResponse(release)
	return response, nil
	*/
	return &CreateResponse{Status: "success"}, nil
}
