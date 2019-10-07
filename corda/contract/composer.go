package contract

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/dovetail-cli/model"
	"github.com/project-flogo/core/trigger"
)

func processComposerTrigger(appName string, txnTrigger *trigger.Config) (ModelResources, error) {
	model := ModelResources{}
	model.AppName = appName
	model.Schemas = make(map[string]string)

	json.Unmarshal([]byte(txnTrigger.Settings["assets"].(string)), &model.Assets)
	json.Unmarshal([]byte(txnTrigger.Settings["transactions"].(string)), &model.Transactions)

	schemas := struct{ schemas [][]string }{}
	json.Unmarshal([]byte(txnTrigger.Settings["schemas"].(string)), &schemas.schemas)

	for _, value := range schemas.schemas {
		model.Schemas[value[0]] = value[1]
	}

	return model, nil
}

func createComposerContractFiles(kotlindir string, opts *Options, r ModelResources, schedulables map[string]string, appclass string) error {
	models := parseAllResources(r.Schemas)

	data, concepts, err := prepareContractStateData(opts, r, models, schedulables)
	if err != nil {
		return fmt.Errorf("prepareContractStateData err %v", err)
	}

	//create ContractState
	for _, s := range data.States {
		err = createKotlinFile(kotlindir, opts.Namespace, s, "kotlin.state.template", fmt.Sprintf("%s%s", s.Class, ".kt"))
		if err != nil {
			return fmt.Errorf("createContractStateFile kotlin.state.template err %v", err)
		}
	}

	//create Contract
	data.AppClass = appclass
	err = createKotlinFile(kotlindir, opts.Namespace, data, "kotlin.contract.template", fmt.Sprintf("%s%s", data.ContractClass, ".kt"))
	if err != nil {
		return fmt.Errorf("createContractFile kotlin.contract.template err %v", err)
	}

	//create common Concept files
	err = createConceptKotlinFiles(kotlindir, "kotlin.concept.template", concepts)
	if err != nil {
		return fmt.Errorf("createConceptJavaFiles kotlin.concept.template err %v", err)
	}
	return nil
}
func prepareContractStateData(opts *Options, flow ModelResources, models map[string]*model.ResourceMetadataModel, schedulables map[string]string) (data ContractData, conceptdata []DataState, err error) {
	logger.Println("Prepare contract state data ....")
	assets := flow.Assets
	transactions := flow.Transactions
	data = ContractData{NS: opts.Namespace, Flow: opts.ModelFile}

	fmt.Printf("***** Contract name = %s.%s%s\n", opts.Namespace, flow.AppName, "Contract")
	data.ContractClass = fmt.Sprintf("%s%s", flow.AppName, "Contract")
	conceptdata = make([]DataState, 0)
	states := make([]DataState, 0)
	concepts := make(map[string]string)

	//assets
	if opts.State != "" {
		state, process, err := prepareState(opts.State, opts.Namespace, models, concepts)
		if err != nil {
			return data, conceptdata, err
		}
		state.ContractClass = fmt.Sprintf("%s.%s%s", opts.Namespace, flow.AppName, "Contract")
		if process {
			if schedulable, ok := schedulables[opts.State]; ok {
				state.IsSchedulable = true
				state.ScheduledActivity = schedulable
			}
			states = append(states, state)
		}
	} else {
		for _, asset := range assets {
			metadata, ok := models[asset]
			if ok && !metadata.Metadata.IsAbstract {
				state, process, err := prepareState(asset, opts.Namespace, models, concepts)
				if err != nil {
					return data, conceptdata, err
				}
				state.ContractClass = fmt.Sprintf("%s.%s%s", opts.Namespace, flow.AppName, "Contract")
				if process {
					if schedulable, ok := schedulables[asset]; ok {
						state.IsSchedulable = true
						state.ScheduledActivity = schedulable
					}
					states = append(states, state)
				}
			}
		}
	}
	data.States = states

	//transactions
	var commands []string
	if len(opts.Commands) > 0 {
		commands = opts.Commands
	} else {
		commands = transactions
	}
	commandData, err := prepareCommands(opts.Namespace, commands, models, concepts)
	if err != nil {
		return data, conceptdata, err
	}
	data.Commands = commandData

	//common concepts
	for _, nm := range concepts {
		concept := DataState{}
		ns, clazz := splitNamespace(nm, opts.Namespace)
		concept.NS = ns
		concept.Class = clazz
		concept.Parent = models[nm].Metadata.Parent
		concept.Attributes = make([]model.ResourceAttribute, len(models[nm].Attributes))
		copy(concept.Attributes, models[nm].Attributes)
		conceptdata = append(conceptdata, concept)
	}
	return data, conceptdata, nil
}
func prepareState(asset, targetns string, models map[string]*model.ResourceMetadataModel, concepts map[string]string) (state DataState, generate bool, err error) {
	state = DataState{}
	ns, class := splitNamespace(asset, targetns)
	if ns != targetns {
		return DataState{}, false, nil
	}

	state.NS = ns
	state.Class = class
	state.CordaClass = strings.TrimSpace(getParentCordaClass(asset, models))

	state.Attributes = make([]model.ResourceAttribute, len(models[asset].Attributes))
	metadata := models[asset]
	copy(state.Attributes, metadata.Attributes)

	//check CordaParticipants decorator
	participants := make([]string, 0)
	for _, decorator := range metadata.Metadata.Decorators {
		if decorator.Name == "CordaParticipants" {
			for _, arg := range decorator.Args {
				participants = append(participants, arg)
			}
			break
		}
	}

	for idx, p := range participants {
		if !strings.HasPrefix(p, "$tx.") {
			return state, false, fmt.Errorf("CordaParticipants %s should be in the format of $tx.path.to.party", p)
		}
		participants[idx] = strings.TrimPrefix(p, "$tx.")
	}

	if len(participants) == 0 {
		//default to all top level party objects
		for _, attr := range state.Attributes {
			if attr.Type == "com.tibco.dovetail.system.Party" && attr.IsRef {
				participants = append(participants, toParticipant(attr.Name, attr))
			}
		}
	}

	if len(participants) == 0 {
		return state, false, fmt.Errorf("There must be at least one Party objects defined either as top Asset attribute or through CordaParticipants decorator")
	}

	state.Participants = participants

	for idx, attr := range state.Attributes {
		m := models[attr.Type]
		if m != nil {
			if m.Metadata.CordaClass != "" {
				state.Attributes[idx].Type = strings.TrimSpace(m.Metadata.CordaClass)
			}
		}
	}

	addCommonConcept(asset, models, concepts)
	return state, true, nil
}
func prepareCommands(targetns string, commands []string, models map[string]*model.ResourceMetadataModel, concepts map[string]string) ([]Command, error) {
	data := make([]Command, 0)
	for _, cmd := range commands {
		logger.Printf("Procesing command %s\n", cmd)
		command := Command{}
		command.Attributes = make([]model.ResourceAttribute, len(models[cmd].Attributes))
		ns, nm := splitNamespace(cmd, targetns)
		if ns != targetns {
			continue
		}
		command.Name = nm
		command.NS = ns

		m := models[cmd]
		if m == nil {
			return nil, fmt.Errorf("%s is not found in model", cmd)
		}

		isQuery := false
		for _, decorator := range m.Metadata.Decorators {
			if decorator.Name == "Query" {
				isQuery = true
				break
			}
		}

		if isQuery {
			continue
		}
		copy(command.Attributes, models[cmd].Attributes)
		data = append(data, command)

		addCommonConcept(cmd, models, concepts)
	}
	return data, nil
}
func addCommonConcept(resource string, models map[string]*model.ResourceMetadataModel, concepts map[string]string) map[string]string {
	am, ok := models[resource]
	if ok == true {
		if am.Metadata.Parent != "" {
			concepts = addCommonConcept(am.Metadata.Parent, models, concepts)
		}

		for _, attr := range am.Attributes {
			m := models[attr.Type]
			if m != nil && strings.Compare(strings.ToUpper(m.Metadata.Type), "CONCEPT") == 0 && m.Metadata.CordaClass == "" {
				concepts[attr.Type] = attr.Type
				concepts = addCommonConcept(attr.Type, models, concepts)
			}
		}
	}
	return concepts
}

func getParentCordaClass(asset string, models map[string]*model.ResourceMetadataModel) string {
	model := models[asset]
	cordaClass := ""
	if model.Metadata.Parent != "" {
		parent := models[model.Metadata.Parent]
		if parent != nil {
			if parent.Metadata.CordaClass != "" {
				cordaClass = parent.Metadata.CordaClass
			} else {
				if parent.Metadata.Parent != "" {
					return getParentCordaClass(model.Metadata.Parent, models)
				}
			}
		}
	}
	return cordaClass
}
