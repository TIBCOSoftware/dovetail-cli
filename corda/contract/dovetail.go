package contract

import (
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/dovetail-cli/model"
	"github.com/project-flogo/core/trigger"
)

func processDovetailTrigger(txnTrigger *trigger.Config) ([]ModelResources, error) {
	txnsbyasset := make(map[string]ModelResources)
	for _, h := range txnTrigger.Handlers {
		asset := h.Settings["assetname"].(string)
		ns, name := splitNamespace(asset, "")
		txn := model.GetFlowName(h.Actions[0].Settings["flowURI"].(string))

		if r, ok := txnsbyasset[asset]; !ok {
			r = ModelResources{Assets: make([]string, 0), Transactions: make([]string, 0), Schemas: make(map[string]string)}
			r.AppName = name
			r.Assets = append(r.Assets, asset)
			r.Schemas[asset] = h.Settings["assetschema"].(string)
			txnsbyasset[asset] = r
		}
		r := txnsbyasset[asset]
		r.Transactions = append(r.Transactions, ns+"."+txn)
		r.Schemas[ns+"."+txn] = h.Schemas.Output["transactionInput"].(map[string]interface{})["value"].(string)
		txnsbyasset[asset] = r
	}

	resources := make([]ModelResources, 0)
	for _, res := range txnsbyasset {
		resources = append(resources, res)
	}

	return resources, nil
}

func createDovetailContractFiles(kotlindir string, opts *Options, resources []ModelResources, schedulables map[string]string, appclass string) error {
	for _, r := range resources {
		models := parseAllResources(r.Schemas)
		data, err := prepareContractData(opts, r, models, schedulables)
		if err != nil {
			return fmt.Errorf("prepareContractStateData err %v", err)
		}
		data.AppClass = appclass

		//create ContractState
		for _, s := range data.States {
			err = createKotlinFile(kotlindir, s.NS, s, "kotlin.state.template", fmt.Sprintf("%s%s", s.Class, ".kt"))
			if err != nil {
				return fmt.Errorf("createContractStateFile kotlin.state.template err %v", err)
			}
		}

		//create Contract
		err = createKotlinFile(kotlindir, data.NS, data, "kotlin.contract.template", fmt.Sprintf("%s%s", data.ContractClass, ".kt"))
		if err != nil {
			return fmt.Errorf("createContractFile kotlin.contract.template err %v", err)
		}
	}
	return nil
}

func prepareContractData(opts *Options, r ModelResources, models map[string]*model.ResourceMetadataModel, schedulables map[string]string) (data ContractData, err error) {
	logger.Println("Prepare contract state data ....")
	asset := r.Assets[0]
	commands := r.Transactions
	ns, name := splitNamespace(asset, "")
	data = ContractData{NS: ns, Flow: opts.ModelFile}

	logger.Printf("***** Contract name = %s.%s%s\n", ns, name, "Contract")
	data.ContractClass = fmt.Sprintf("%s%s", name, "Contract")
	states := make([]DataState, 0)

	metadata, ok := models[asset]
	if ok {
		state, err := prepareStateData(asset, metadata)
		if err != nil {
			return data, err
		}

		if schedulable, ok := schedulables[asset]; ok {
			state.IsSchedulable = true
			state.ScheduledActivity = schedulable
		}
		states = append(states, state)
	}

	data.States = states

	//transactions
	commandData, err := prepareCmdData(commands, models)
	if err != nil {
		return data, err
	}
	data.Commands = commandData

	return data, nil
}
func prepareStateData(asset string, r *model.ResourceMetadataModel) (state DataState, err error) {
	state = DataState{}
	ns, class := splitNamespace(asset, "")

	state.NS = ns
	state.Class = class
	state.ContractClass = state.Class + "Contract"
	state.CordaClass = strings.TrimSpace(getCordaClass(asset, r))

	state.Attributes = make([]model.ResourceAttribute, len(r.Attributes))
	copy(state.Attributes, r.Attributes)

	//get participants
	//TODO: exit keys for fungible asset
	participants, exits := prepareParticipants(asset, r)

	if len(participants) == 0 {
		return state, fmt.Errorf("There must be at least one Party assigned as move signers")
	}

	state.Participants = participants
	state.ExitKeys = exits

	return state, nil
}
func prepareCmdData(commands []string, models map[string]*model.ResourceMetadataModel) ([]Command, error) {
	data := make([]Command, 0)
	for _, cmd := range commands {
		logger.Printf("Procesing command %s\n", cmd)
		command := Command{}
		command.Attributes = make([]model.ResourceAttribute, len(models[cmd].Attributes))
		ns, nm := splitNamespace(cmd, "")

		command.Name = nm
		command.NS = ns

		m := models[cmd]
		if m == nil {
			return nil, fmt.Errorf("%s is not found in contract", cmd)
		}

		copy(command.Attributes, models[cmd].Attributes)
		data = append(data, command)
	}
	return data, nil
}

func getCordaClass(asset string, model *model.ResourceMetadataModel) string {

	if model.Metadata.Parent == "com.tibco.dovetail.system.LinearState" {
		return "net.corda.core.contracts.LinearState"
	} else if model.Metadata.Parent == "com.tibco.dovetail.system.FungibleAsset" {
		return "net.corda.core.contracts.FungibleAsset"
	}
	return ""
}

func prepareParticipants(asset string, model *model.ResourceMetadataModel) (participants []string, exits []string) {
	return model.Metadata.Participants, model.Metadata.ExitSigners
}
