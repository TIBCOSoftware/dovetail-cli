package config

const (
	// DefaultFolderName default name for dovetail folder
	DefaultFolderName = ".dovetail-cli"

	// Types of supported blockchains

	// HyperledgerFabric blockchain name
	HyperledgerFabric = "fabric"
	// Corda blockchain name
	Corda = "corda"
	// Sawtooth blockchain name
	Sawtooth = "sawtooth"

	// NodeVerboseKey is the viper key for node verbose
	NodeVerboseKey = "verbose"
	// NodePortKey is the viper key for node port
	NodePortKey = "port"
	// NodePortShortKey is the viper key for node port short
	NodePortShortKey = "p"
	// NodePortDefault is the default value for the node port
	NodePortDefault = "8080"
)
