package model

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/project-flogo/core/app"
)

type ResourceMetadataModel struct {
	Metadata struct {
		Type         string   `json:"type"`
		Parent       string   `json:"parent, omitempty"`
		IssueSigners []string `json:"issueSigners, omitempty"`
		Participants []string `json:"participants, omitempty"`
		ExitSigners  []string `json:"exitSigners, omitempty"`
		CordaClass   string   `json:"cordaClass, omitempty"`
		Name         []string `json:"name, omitempty"`
		Module       []string `json:"module, omitempty"`
		IdentifiedBy string   `json:"identifiedBy, omitempty"`
		IsAbstract   bool     `json:"isAbstract, omitempty"`
		Decorators   []struct {
			Name string   `json:"name"`
			Args []string `json:"args, omitempty"`
		} `json:"decorators, omitempty"`
	} `json:"metadata"`
	Attributes []ResourceAttribute `json:"attributes"`
}

type ResourceAttribute struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	IsRef      bool   `json:"isRef, omitempty"`
	IsArray    bool   `json:"isArray, omitempty"`
	IsOptional bool   `json:"isOptional, omitempty"`
	PartyType  string `json:"partyType, omitempty`
}

func ParseResourceModel(jsonResource string) *ResourceMetadataModel {
	resource := ResourceMetadataModel{}

	json.Unmarshal([]byte(jsonResource), &resource)
	return &resource
}

// ParseApp parses the model file into an app.Config struct
func ParseApp(modelfile string) (*app.Config, error) {

	flowjson, err := ioutil.ReadFile(modelfile)
	if err != nil {
		return nil, err
	}

	return DecodeApp(bytes.NewReader(flowjson))
}

// decodeApp decodes the model file into an app.Config struct
func DecodeApp(r io.Reader) (*app.Config, error) {
	appCfg := &app.Config{}

	jsonParser := json.NewDecoder(r)
	err := jsonParser.Decode(&appCfg)
	if err != nil {
		return nil, err
	}

	return appCfg, nil
}

func GetFlowName(flowuri string) string {
	return flowuri[11:]
}
