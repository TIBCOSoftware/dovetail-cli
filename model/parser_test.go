package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	jsonstring := `{
		"imports": [
			"github.com/project-flogo/flow",
			"github.com/TIBCOSoftware/dovetail-contrib/General/activity/mapper",
			"github.com/TIBCOSoftware/dovetail-contrib/CorDApp/trigger/flowinitiator"
		],
		"name": "testing",
		"description": " ",
		"version": "1.1.0",
		"type": "flogo:app",
		"appModel": "1.1.0",
		"feVersion": "2.7.0",
		"triggers": [
			{
				"ref": "#flowinitiator",
				"name": "CorDAppFlowInitiator",
				"description": "",
				"settings": {
					"flowType": "initiator"
				},
				"id": "CorDAppFlowInitiator",
				"handlers": [
					{
						"description": "",
						"settings": {
							"useAnonymousIdentity": true,
							"hasObservers": true,
							"observerManual": true,
							"schemaSelection": "user",
							"inputParams": {
								"metadata": "{\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"type\":\"array\",\"items\":{\"type\":\"object\",\"properties\":{\"parameterName\":{\"type\":\"string\"},\"type\":{\"type\":\"string\"},\"repeating\":{\"type\":\"string\"},\"required\":{\"type\":\"string\"},\"isEditable\":{\"type\":\"boolean\"},\"partyType\":{\"type\":\"string\"}}}}",
								"value": "",
								"fe_metadata": "[{\"parameterName\":\"linearId\",\"type\":\"String\",\"repeating\":\"false\",\"required\":\"false\",\"isEditable\":true,\"partyType\":\"\"}]"
							}
						},
						"action": {
							"ref": "github.com/project-flogo/flow",
							"settings": {
								"flowURI": "res://flow:flow1"
							},
							"input": {
								"transactionInput": {
									"mapping": {
										"linearId": "=$.transactionInput.linearId"
									}
								},
								"ourIdentity": "=$.ourIdentity"
							}
						},
						"schemas": {
							"output": {
								"transactionInput": {
									"type": "json",
									"value": "{\"schema\":\"http://json-schema.org/draft-04/schema#\",\"type\":\"object\",\"properties\":{\"linearId\":{\"type\":\"string\"}},\"description\":\"{\\\"metadata\\\":{\\\"type\\\":\\\"Transaction\\\"},\\\"attributes\\\":[{\\\"name\\\":\\\"linearId\\\",\\\"type\\\":\\\"String\\\",\\\"isRef\\\":false,\\\"isArray\\\":false,\\\"partyType\\\":\\\"\\\"}]}\"}",
									"fe_metadata": "{\"schema\":\"http://json-schema.org/draft-04/schema#\",\"type\":\"object\",\"properties\":{\"linearId\":{\"type\":\"string\"}},\"description\":\"{\\\"metadata\\\":{\\\"type\\\":\\\"Transaction\\\"},\\\"attributes\\\":[{\\\"name\\\":\\\"linearId\\\",\\\"type\\\":\\\"String\\\",\\\"isRef\\\":false,\\\"isArray\\\":false,\\\"partyType\\\":\\\"\\\"}]}\"}"
								}
							}
						}
					}
				]
			}
		],
		"resources": [
			{
				"id": "flow:flow1",
				"data": {
					"name": "flow1",
					"description": "",
					"links": [],
					"tasks": [
						{
							"id": "Mapper",
							"name": "Mapper",
							"activity": {
								"ref": "#mapper",
								"input": {
									"model": "",
									"dataType": "String",
									"isArray": false,
									"inputArrayType": "Object Array",
									"outputArrayType": "Object Array",
									"rounding": "HALF_EVEN",
									"format": "yyyy-MM-dd HH:mm:ss ZZZ",
									"input": {
										"mapping": {
											"field": "=$flow.transactionInput.linearId"
										}
									}
								},
								"schemas": {
									"input": {
										"input": {
											"type": "json",
											"value": "{\"type\":\"object\",\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"required\":[\"field\"],\"properties\":{\"field\":{\"type\":\"string\"}}}",
											"fe_metadata": "{\"type\":\"object\",\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"required\":[\"field\"],\"properties\":{\"field\":{\"type\":\"string\"}}}"
										}
									},
									"output": {
										"output": {
											"type": "json",
											"value": "{\"type\":\"object\",\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"required\":[\"field\"],\"properties\":{\"field\":{\"type\":\"string\"}}}",
											"fe_metadata": "{\"type\":\"object\",\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"required\":[\"field\"],\"properties\":{\"field\":{\"type\":\"string\"}}}"
										}
									}
								}
							}
						}
					],
					"metadata": {
						"input": [
							{
								"name": "transactionInput",
								"type": "object",
								"schema": {
									"type": "json",
									"value": "{\"linearId\":{\"type\":\"string\"}}"
								}
							},
							{
								"name": "ourIdentity",
								"type": "string"
							}
						],
						"output": [],
						"fe_metadata": {
							"input": "{\"type\":\"object\",\"title\":\"CorDAppFlowInitiator\",\"properties\":{\"transactionInput\":{\"schema\":\"http://json-schema.org/draft-04/schema#\",\"type\":\"object\",\"properties\":{\"linearId\":{\"type\":\"string\"}},\"description\":\"{\\\"metadata\\\":{\\\"type\\\":\\\"Transaction\\\"},\\\"attributes\\\":[{\\\"name\\\":\\\"linearId\\\",\\\"type\\\":\\\"String\\\",\\\"isRef\\\":false,\\\"isArray\\\":false,\\\"partyType\\\":\\\"\\\"}]}\"},\"ourIdentity\":{\"type\":\"string\",\"required\":false}}}"
						}
					}
				}
			}
		],
		"properties": [],
		"contrib": "W3sicmVmIjoiZ2l0aHViLmNvbS9USUJDT1NvZnR3YXJlL2RvdmV0YWlsLWNvbnRyaWIvQ29yREFwcCIsInMzbG9jYXRpb24iOiJ7VVNFUklEfS9Db3JEQXBwIn0seyJyZWYiOiJnaXRodWIuY29tL1RJQkNPU29mdHdhcmUvZG92ZXRhaWwtY29udHJpYi9HZW5lcmFsIiwiczNsb2NhdGlvbiI6IntVU0VSSUR9L0dlbmVyYWwifV0=",
		"fe_metadata": "UEsDBAoAAAAIAMwE+E6rqSRGFQAAABMAAAAIAAAAYXBwLmpzb26rViopykxPTy1yy8kvL1ayio6tBQBQSwECFAAKAAAACADMBPhOq6kkRhUAAAATAAAACAAAAAAAAAAAAAAAAAAAAAAAYXBwLmpzb25QSwUGAAAAAAEAAQA2AAAAOwAAAAAA"
	}`

	model, err := DecodeApp(strings.NewReader(jsonstring))
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, model)

	assert.True(t, len(model.Triggers) > 0, "Triggers should not be 0")
	assert.True(t, len(model.Resources) > 0, "Resources should not be 0")

}
