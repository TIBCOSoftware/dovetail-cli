package config

const (
	// DEFAULT_FOLDER_NAME default name for dovetail folder
	DEFAULT_FOLDER_NAME = ".dovetail-cli"

	// Types of supported blockchains

	// HYPERLEDGER_FABRIC blockchain name
	HYPERLEDGER_FABRIC = "fabric"
	// CORDA blockchain name
	CORDA = "corda"
	// SAWTOOTH blockchain name
	SAWTOOTH = "sawtooth"

	// Node vars
	NODE_VERBOSE_KEY    = "verbose"
	NODE_PORT_KEY       = "port"
	NODE_PORT_SHORT_KEY = "p"
	NODE_PORT_DEFAULT   = "8080"
)
