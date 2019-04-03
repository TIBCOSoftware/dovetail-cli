package cordapp

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestParser(t *testing.T) {
	/*	jsonstring := `{
			"name": "IOUApp",
			"description": " ",
			"version": "1.0.0",
			"type": "flogo:app",
			"appModel": "1.0.0",
			"resources": [
			  {
				"id": "flow:IssueIOUInitiator",
				"data": {
				  "name": "IssueIOUInitiator",
				  "description": "",
				  "tasks": [
					{
					  "id": "BuildTransactionProposal",
					  "name": "BuildTransactionProposal",
					  "activity": {
						"ref": "github.com/TIBCOSoftware/dovetail-contrib/CorDApp/activity/txnbuilder",
						"input": {
						  "notary": "C=FR,L=Paris,O=Notary",
						  "contract": "daaee9a0-ec38-11e8-8227-afbc072d2792",
						  "contractClass": "com.example.iou.IOUContract",
						  "command": "com.example.iou.IssueIOU",
						  "input": {
							"metadata": "{\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"title\":\"IssueIOU\",\"type\":\"object\",\"properties\":{\"iou\":{\"type\":\"object\",\"properties\":{\"issuer\":{\"type\":\"string\"},\"owner\":{\"type\":\"string\"},\"amt\":{\"type\":\"object\",\"properties\":{\"currency\":{\"type\":\"string\"},\"quantity\":{\"type\":\"integer\",\"default\":\"0\"}},\"required\":[\"currency\",\"quantity\"]},\"linearId\":{\"type\":\"string\"}},\"required\":[\"issuer\",\"owner\",\"amt\",\"linearId\"]},\"transactionId\":{\"type\":\"string\"},\"timestamp\":{\"format\":\"date-time\",\"type\":\"string\"}},\"required\":[\"iou\",\"transactionId\",\"timestamp\"],\"description\":\"{\\\"metadata\\\":{\\\"type\\\":\\\"Transaction\\\",\\\"parent\\\":\\\"org.hyperledger.composer.system.Transaction\\\",\\\"isAbstract\\\":false,\\\"identifiedBy\\\":\\\"transactionId\\\",\\\"decorators\\\":[{\\\"name\\\":\\\"InitiatedBy\\\",\\\"args\\\":[\\\"$tx.iou.issuer\\\"]}]},\\\"attributes\\\":[{\\\"name\\\":\\\"iou\\\",\\\"isOptional\\\":false,\\\"type\\\":\\\"com.example.iou.IOU\\\"},{\\\"name\\\":\\\"transactionId\\\",\\\"isOptional\\\":false,\\\"type\\\":\\\"String\\\"},{\\\"name\\\":\\\"timestamp\\\",\\\"isOptional\\\":false,\\\"type\\\":\\\"DateTime\\\"}]}\"}",
							"value": ""
						  }
						},
						"output": {},
						"mappings": {
						  "input": [
							{
							  "mapTo": "$INPUT['input']['iou']['issuer']",
							  "type": "expression",
							  "value": "$flow.ourIdentity"
							},
							{
							  "mapTo": "$INPUT['input']['iou']['owner']",
							  "type": "expression",
							  "value": "$flow.transactionInput.owner"
							},
							{
							  "mapTo": "$INPUT['input']['iou']['amt']",
							  "type": "expression",
							  "value": "$flow.transactionInput.amt"
							},
							{
							  "mapTo": "$INPUT['input']['iou']['linearId']",
							  "type": "expression",
							  "value": "$flow.transactionInput.id"
							},
							{
							  "mapTo": "$INPUT['input']['transactionId']",
							  "type": "expression",
							  "value": "txnId"
							},
							{
							  "mapTo": "$INPUT['input']['timestamp']",
							  "type": "expression",
							  "value": "2019-3-25T13:00:00"
							}
						  ]
						}
					  }
					},
					{
					  "id": "InitiatorSignandFinalizeaTransaction",
					  "name": "InitiatorSignandFinalizeaTransaction",
					  "activity": {
						"ref": "github.com/TIBCOSoftware/dovetail-contrib/CorDApp/activity/finalize",
						"input": {},
						"output": {}
					  }
					}
				  ],
				  "links": [
					{
					  "id": 1,
					  "from": "BuildTransactionProposal",
					  "to": "InitiatorSignandFinalizeaTransaction",
					  "type": "default"
					}
				  ],
				  "metadata": {
					"input": [],
					"output": []
				  }
				}
			  },
			  {
				"id": "flow:IssueIOUResponder",
				"data": {
				  "name": "IssueIOUResponder",
				  "description": "",
				  "tasks": [
					{
					  "id": "ReceiverVerifyandSignTransaction",
					  "name": "ReceiverVerifyandSignTransaction",
					  "activity": {
						"ref": "github.com/TIBCOSoftware/dovetail-contrib/CorDApp/activity/receiversign",
						"input": {},
						"output": {}
					  }
					},
					{
					  "id": "CommitSignedandNotorizedTxn",
					  "name": "CommitSignedandNotorizedTxn",
					  "activity": {
						"ref": "github.com/TIBCOSoftware/dovetail-contrib/CorDApp/activity/finalize",
						"input": {},
						"output": {}
					  }
					}
				  ],
				  "links": [
					{
					  "id": 1,
					  "from": "ReceiverVerifyandSignTransaction",
					  "to": "CommitSignedandNotorizedTxn",
					  "type": "default"
					}
				  ],
				  "metadata": {
					"input": [],
					"output": []
				  }
				}
			  }
			],
			"triggers": [
			  {
				"ref": "github.com/TIBCOSoftware/dovetail-contrib/CorDApp/trigger/flowinitiator",
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
					  "allowRPCClient": true,
					  "useAnonymousIdentity": false,
					  "inputParams": {
						"metadata": "{\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"type\":\"array\",\"items\":{\"type\":\"object\",\"properties\":{\"parameterName\":{\"type\":\"string\"},\"type\":{\"type\":\"string\"},\"repeating\":{\"type\":\"string\"},\"required\":{\"type\":\"string\"},\"isEditable\":{\"type\":\"boolean\"}},\"required\":[\"parameterName\",\"type\",\"repeating\",\"required\",\"isEditable\"]}}",
						"value": ""
					  }
					},
					"outputs": {
					  "transactionInput": {
						"metadata": "{\"schema\":\"http://json-schema.org/draft-04/schema#\",\"type\":\"object\",\"properties\":{\"owner\":{\"type\":\"string\"},\"amt\":{\"type\":\"object\",\"properties\":{\"currency\":{\"type\":\"string\"},\"quantity\":{\"type\":\"number\"}}},\"id\":{\"type\":\"string\"}},\"description\":\"{\\\"metadata\\\":{\\\"type\\\":\\\"Transaction\\\"},\\\"attributes\\\":[{\\\"name\\\":\\\"owner\\\",\\\"type\\\":\\\"net.corda.core.identity.Party\\\",\\\"isRef\\\":true,\\\"isArray\\\":false},{\\\"name\\\":\\\"amt\\\",\\\"type\\\":\\\"net.corda.core.contracts.Amount<Currency>\\\",\\\"isRef\\\":false,\\\"isArray\\\":false},{\\\"name\\\":\\\"id\\\",\\\"type\\\":\\\"string\\\",\\\"isRef\\\":false,\\\"isArray\\\":false}]}\"}",
						"value": ""
					  },
					  "ourIdentity": ""
					},
					"action": {
					  "ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
					  "data": {
						"flowURI": "res://flow:IssueIOUInitiator"
					  },
					  "mappings": {
						"input": [],
						"output": []
					  }
					}
				  }
				]
			  },
			  {
				"ref": "github.com/TIBCOSoftware/dovetail-contrib/CorDApp/trigger/flowreceiver",
				"name": "CorDAppFlowReceiver",
				"description": "",
				"settings": {
				  "flowType": "receiver"
				},
				"id": "CorDAppFlowReceiver",
				"handlers": [
				  {
					"description": "",
					"settings": {
					  "initiatorFlow": "IssueIOUInitiator"
					},
					"outputs": {},
					"action": {
					  "ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
					  "data": {
						"flowURI": "res://flow:IssueIOUResponder"
					  },
					  "mappings": {
						"input": [],
						"output": []
					  }
					}
				  }
				]
			  }
			],
			"ui": "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",
			"contrib": "W3sicmVmIjoiZ2l0aHViLmNvbS9USUJDT1NvZnR3YXJlL2RvdmV0YWlsLWNvbnRyaWIvQ29yREFwcCIsInMzbG9jYXRpb24iOiJ7VVNFUklEfS9Db3JEQXBwIn1d"
		  }`
	*/
	wdir, _ := filepath.Abs("/Users/mwenyan@tibco.com/Downloads/tmp")

	opts := NewOptions("/Users/mwenyan@tibco.com/dovetail/src/github.com/TIBCOSoftware/dovetail-cli/corda/cordapp/iou.json", "1.0.0", wdir, "com.example.flows.iou")
	generator := NewGenerator(opts)

	//config, err := model.DecodeApp(bytes.NewReader([]byte(jsonstring)))
	err := generator.Generate()
	if err != nil {
		fmt.Errorf("generate err %v", err)
	}
}
