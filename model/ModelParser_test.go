package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	jsonstring := `{
		"name": "IOU",
		"description": " ",
		"version": "1.0.0",
		"type": "flogo:app",
		"appModel": "1.0.0",
		"resources": [
		 {
		  "id": "flow:flow1",
		  "data": {
		   "name": "flow1",
		   "description": "",
		   "tasks": [],
		   "links": [],
		   "metadata": {
			"input": [],
			"output": []
		   }
		  }
		 }
		],
		"triggers": [
		 {
		  "ref": "chaincode/trigger/chaincode",
		  "name": "ChaincodeTrigger",
		  "description": "",
		  "settings": {
		   "clearAppCache": false,
		   "model": "b949d740-53cf-11e8-86f2-7b5f42c923d6",
		   "createAll": true,
		   "assets": "[\"com.tibco.cp.IOU\",\"com.tibco.dovetail.system.Cash\"]",
		   "transactions": "[\"com.tibco.cp.IssueIOU\",\"com.tibco.cp.TransferIOU\",\"com.tibco.cp.SettleIOU\"]",
		   "schemas": "[[\"com.tibco.cp.IOU\",\"{\\\"$schema\\\":\\\"http://json-schema.org/draft-04/schema#\\\",\\\"title\\\":\\\"IOU\\\",\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"lender\\\":{\\\"type\\\":\\\"string\\\"},\\\"borrower\\\":{\\\"type\\\":\\\"string\\\"},\\\"amt\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]},\\\"paid\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]},\\\"linearId\\\":{\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[\\\"lender\\\",\\\"borrower\\\",\\\"amt\\\",\\\"paid\\\",\\\"linearId\\\"],\\\"description\\\":\\\"{\\\\\\\"metadata\\\\\\\":{\\\\\\\"type\\\\\\\":\\\\\\\"Asset\\\\\\\",\\\\\\\"parent\\\\\\\":\\\\\\\"com.tibco.dovetail.system.LinearState\\\\\\\",\\\\\\\"identifiedBy\\\\\\\":\\\\\\\"linearId\\\\\\\"},\\\\\\\"refAttributes\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"lender\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.dovetail.system.Party\\\\\\\"},{\\\\\\\"name\\\\\\\":\\\\\\\"borrower\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.dovetail.system.Party\\\\\\\"}]}\\\"}\"],[\"com.tibco.cp.IssueIOU\",\"{\\\"$schema\\\":\\\"http://json-schema.org/draft-04/schema#\\\",\\\"title\\\":\\\"IssueIOU\\\",\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"iou\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"lender\\\":{\\\"type\\\":\\\"string\\\"},\\\"borrower\\\":{\\\"type\\\":\\\"string\\\"},\\\"amt\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]},\\\"paid\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]},\\\"linearId\\\":{\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[\\\"lender\\\",\\\"borrower\\\",\\\"amt\\\",\\\"paid\\\",\\\"linearId\\\"]},\\\"transactionId\\\":{\\\"type\\\":\\\"string\\\",\\\"description\\\":\\\"The instance identifier for this type\\\"},\\\"timestamp\\\":{\\\"format\\\":\\\"date-time\\\",\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[\\\"iou\\\",\\\"transactionId\\\",\\\"timestamp\\\"],\\\"description\\\":\\\"{\\\\\\\"metadata\\\\\\\":{\\\\\\\"type\\\\\\\":\\\\\\\"Transaction\\\\\\\",\\\\\\\"parent\\\\\\\":\\\\\\\"org.hyperledger.composer.system.Transaction\\\\\\\",\\\\\\\"identifiedBy\\\\\\\":\\\\\\\"transactionId\\\\\\\"},\\\\\\\"refAttributes\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"lender\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.dovetail.system.Party\\\\\\\"},{\\\\\\\"name\\\\\\\":\\\\\\\"borrower\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.dovetail.system.Party\\\\\\\"}]}\\\"}\"],[\"com.tibco.cp.TransferIOU\",\"{\\\"$schema\\\":\\\"http://json-schema.org/draft-04/schema#\\\",\\\"title\\\":\\\"TransferIOU\\\",\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"iou\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"lender\\\":{\\\"type\\\":\\\"string\\\"},\\\"borrower\\\":{\\\"type\\\":\\\"string\\\"},\\\"amt\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]},\\\"paid\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]},\\\"linearId\\\":{\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[\\\"lender\\\",\\\"borrower\\\",\\\"amt\\\",\\\"paid\\\",\\\"linearId\\\"]},\\\"newLender\\\":{\\\"type\\\":\\\"string\\\"},\\\"transactionId\\\":{\\\"type\\\":\\\"string\\\",\\\"description\\\":\\\"The instance identifier for this type\\\"},\\\"timestamp\\\":{\\\"format\\\":\\\"date-time\\\",\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[\\\"iou\\\",\\\"newLender\\\",\\\"transactionId\\\",\\\"timestamp\\\"],\\\"description\\\":\\\"{\\\\\\\"metadata\\\\\\\":{\\\\\\\"type\\\\\\\":\\\\\\\"Transaction\\\\\\\",\\\\\\\"parent\\\\\\\":\\\\\\\"org.hyperledger.composer.system.Transaction\\\\\\\",\\\\\\\"identifiedBy\\\\\\\":\\\\\\\"transactionId\\\\\\\"},\\\\\\\"refAttributes\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"iou\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.cp.IOU\\\\\\\"},{\\\\\\\"name\\\\\\\":\\\\\\\"newLender\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.dovetail.system.Party\\\\\\\"}]}\\\"}\"],[\"com.tibco.cp.SettleIOU\",\"{\\\"$schema\\\":\\\"http://json-schema.org/draft-04/schema#\\\",\\\"title\\\":\\\"SettleIOU\\\",\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"iou\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"lender\\\":{\\\"type\\\":\\\"string\\\"},\\\"borrower\\\":{\\\"type\\\":\\\"string\\\"},\\\"amt\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]},\\\"paid\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]},\\\"linearId\\\":{\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[\\\"lender\\\",\\\"borrower\\\",\\\"amt\\\",\\\"paid\\\",\\\"linearId\\\"]},\\\"payments\\\":{\\\"type\\\":\\\"array\\\",\\\"items\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"assetId\\\":{\\\"type\\\":\\\"string\\\",\\\"description\\\":\\\"The instance identifier for this type\\\"},\\\"amt\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]},\\\"owner\\\":{\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[\\\"assetId\\\",\\\"amt\\\",\\\"owner\\\"]}},\\\"transactionId\\\":{\\\"type\\\":\\\"string\\\",\\\"description\\\":\\\"The instance identifier for this type\\\"},\\\"timestamp\\\":{\\\"format\\\":\\\"date-time\\\",\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[\\\"iou\\\",\\\"payments\\\",\\\"transactionId\\\",\\\"timestamp\\\"],\\\"description\\\":\\\"{\\\\\\\"metadata\\\\\\\":{\\\\\\\"type\\\\\\\":\\\\\\\"Transaction\\\\\\\",\\\\\\\"parent\\\\\\\":\\\\\\\"org.hyperledger.composer.system.Transaction\\\\\\\",\\\\\\\"identifiedBy\\\\\\\":\\\\\\\"transactionId\\\\\\\"},\\\\\\\"refAttributes\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"iou\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.cp.IOU\\\\\\\"},{\\\\\\\"name\\\\\\\":\\\\\\\"payments\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.dovetail.system.Cash\\\\\\\"}]}\\\"}\"],[\"com.tibco.dovetail.system.Cash\",\"{\\\"$schema\\\":\\\"http://json-schema.org/draft-04/schema#\\\",\\\"title\\\":\\\"Cash\\\",\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"assetId\\\":{\\\"type\\\":\\\"string\\\",\\\"description\\\":\\\"The instance identifier for this type\\\"},\\\"amt\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]},\\\"owner\\\":{\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[\\\"assetId\\\",\\\"amt\\\",\\\"owner\\\"],\\\"description\\\":\\\"{\\\\\\\"metadata\\\\\\\":{\\\\\\\"type\\\\\\\":\\\\\\\"Asset\\\\\\\",\\\\\\\"parent\\\\\\\":\\\\\\\"com.tibco.dovetail.system.FungibleState\\\\\\\",\\\\\\\"identifiedBy\\\\\\\":\\\\\\\"assetId\\\\\\\"},\\\\\\\"refAttributes\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"owner\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.dovetail.system.Party\\\\\\\"}]}\\\"}\"],[\"com.tibco.dovetail.system.Amount\",\"{\\\"$schema\\\":\\\"http://json-schema.org/draft-04/schema#\\\",\\\"title\\\":\\\"Amount\\\",\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"],\\\"description\\\":\\\"{\\\\\\\"metadata\\\\\\\":{\\\\\\\"type\\\\\\\":\\\\\\\"Concept\\\\\\\"},\\\\\\\"refAttributes\\\\\\\":[]}\\\"}\"],[\"com.tibco.dovetail.system.TimeWindow\",\"{\\\"$schema\\\":\\\"http://json-schema.org/draft-04/schema#\\\",\\\"title\\\":\\\"TimeWindow\\\",\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"from\\\":{\\\"format\\\":\\\"date-time\\\",\\\"type\\\":\\\"string\\\"},\\\"until\\\":{\\\"format\\\":\\\"date-time\\\",\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[],\\\"description\\\":\\\"{\\\\\\\"metadata\\\\\\\":{\\\\\\\"type\\\\\\\":\\\\\\\"Concept\\\\\\\"},\\\\\\\"refAttributes\\\\\\\":[]}\\\"}\"]]"
		  },
		  "id": "ChaincodeTrigger",
		  "handlers": [
		   {
			"description": "",
			"settings": {
			 "transaction": "com.tibco.cp.IssueIOU"
			},
			"outputs": {
			 "transactionInput": {
			  "metadata": "{\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"title\":\"IssueIOU\",\"type\":\"object\",\"properties\":{\"iou\":{\"type\":\"object\",\"properties\":{\"lender\":{\"type\":\"string\"},\"borrower\":{\"type\":\"string\"},\"amt\":{\"type\":\"object\",\"properties\":{\"currency\":{\"type\":\"string\"},\"quantity\":{\"type\":\"integer\",\"default\":\"0\"}},\"required\":[\"currency\",\"quantity\"]},\"paid\":{\"type\":\"object\",\"properties\":{\"currency\":{\"type\":\"string\"},\"quantity\":{\"type\":\"integer\",\"default\":\"0\"}},\"required\":[\"currency\",\"quantity\"]},\"linearId\":{\"type\":\"string\"}},\"required\":[\"lender\",\"borrower\",\"amt\",\"paid\",\"linearId\"]},\"transactionId\":{\"type\":\"string\",\"description\":\"The instance identifier for this type\"},\"timestamp\":{\"format\":\"date-time\",\"type\":\"string\"}},\"required\":[\"iou\",\"transactionId\",\"timestamp\"],\"description\":\"{\\\"metadata\\\":{\\\"type\\\":\\\"Transaction\\\",\\\"parent\\\":\\\"org.hyperledger.composer.system.Transaction\\\",\\\"identifiedBy\\\":\\\"transactionId\\\"},\\\"refAttributes\\\":[{\\\"name\\\":\\\"lender\\\",\\\"type\\\":\\\"com.tibco.dovetail.system.Party\\\"},{\\\"name\\\":\\\"borrower\\\",\\\"type\\\":\\\"com.tibco.dovetail.system.Party\\\"}]}\"}",
			  "value": ""
			 }
			},
			"action": {
			 "ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
			 "data": {
			  "flowURI": "res://flow:flow1"
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
		"contrib": "W3sicmVmIjoiY2hhaW5jb2RlIiwiczNsb2NhdGlvbiI6Imh0dHA6Ly9sb2NhbGhvc3Q6ODkwMC93aWNvbnRyaWJ1dGlvbnMvZmxvZ28vY2hhaW5jb2RlIn1d"
	   }`

	model, err := ParseFlowApp(jsonstring)
	assert.Nil(t, err)

	fmt.Printf("%v\n", model.Assets)
	fmt.Printf("%v\n", model.Transactions)
	fmt.Printf("%v\n", model.Schemas)

}
