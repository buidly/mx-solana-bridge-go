package config

import (
	"github.com/multiversx/mx-chain-go/config"
	"github.com/multiversx/mx-solana-bridge-go/clients/chain"
)

// Configs is a holder for the relayer configuration parameters
type Configs struct {
	GeneralConfig   Config
	ApiRoutesConfig ApiRoutesConfig
	FlagsConfig     ContextFlagsConfig
}

// Config general configuration struct
type Config struct {
	Solana            SolanaConfig
	MultiversX        MultiversXConfig
	P2P               ConfigP2P
	StateMachine      map[string]ConfigStateMachine
	Relayer           ConfigRelayer
	Logs              LogsConfig
	Antiflood         AntifloodConfig
	BatchValidator    BatchValidatorConfig
	PeersRatingConfig PeersRatingConfig
}

// SolanaConfig represents the Solana Config parameters
type SolanaConfig struct {
	Chain                              chain.Chain
	NetworkAddress                     string
	WSSNetworkAddress                  string
	BridgeProgramAddress               string
	SafeProgramAddress                 string
	SolanaPrivateKeyFile               string
	PrivateKeyFile                     string
	IntervalToResendTxsInSeconds       uint64
	GasLimitBase                       uint64
	GasLimitForEach                    uint64
	MaxRetriesOnQuorumReached          uint64
	IntervalToWaitForTransferInSeconds uint64
	MaxBlocksDelta                     uint64
}

// ConfigP2P configuration for the P2P communication
type ConfigP2P struct {
	Port            string
	InitialPeerList []string
	ProtocolID      string
	AntifloodConfig config.AntifloodConfig
}

// ConfigRelayer configuration for general relayer configuration
type ConfigRelayer struct {
	Marshalizer          config.MarshalizerConfig
	RoleProvider         RoleProviderConfig
	StatusMetricsStorage config.StorageConfig
}

// ConfigStateMachine the configuration for the state machine
type ConfigStateMachine struct {
	StepDurationInMillis       uint64
	IntervalForLeaderInSeconds uint64
}

// ContextFlagsConfig the configuration for flags
type ContextFlagsConfig struct {
	WorkingDir           string
	LogLevel             string
	DisableAnsiColor     bool
	ConfigurationFile    string
	ConfigurationApiFile string
	SaveLogFile          bool
	EnableLogName        bool
	RestApiInterface     string
	EnablePprof          bool
}

// WebServerAntifloodConfig will hold the anti-flooding parameters for the web server
type WebServerAntifloodConfig struct {
	SimultaneousRequests         uint32
	SameSourceRequests           uint32
	SameSourceResetIntervalInSec uint32
}

// AntifloodConfig will hold all p2p antiflood parameters
type AntifloodConfig struct {
	Enabled   bool
	WebServer WebServerAntifloodConfig
}

// BatchValidatorConfig represents the configuration for the batch validator
type BatchValidatorConfig struct {
	Enabled              bool
	URL                  string
	RequestTimeInSeconds int
}

// ApiRoutesConfig holds the configuration related to Rest API routes
type ApiRoutesConfig struct {
	Logging     ApiLoggingConfig
	APIPackages map[string]APIPackageConfig
}

// ApiLoggingConfig holds the configuration related to API requests logging
type ApiLoggingConfig struct {
	LoggingEnabled          bool
	ThresholdInMicroSeconds int
}

// APIPackageConfig holds the configuration for the routes of each package
type APIPackageConfig struct {
	Routes []RouteConfig
}

// RouteConfig holds the configuration for a single route
type RouteConfig struct {
	Name string
	Open bool
}

// LogsConfig will hold settings related to the logging sub-system
type LogsConfig struct {
	LogFileLifeSpanInSec int
	LogFileLifeSpanInMB  int
}

// RoleProviderConfig is the configuration for the role provider component
type RoleProviderConfig struct {
	PollingIntervalInMillis uint64
}

// MultiversXConfig represents the MultiversX Config parameters
type MultiversXConfig struct {
	NetworkAddress                  string
	NetworkApiAddress               string
	MultisigContractAddress         string
	PrivateKeyFile                  string
	IntervalToResendTxsInSeconds    uint64
	GasMap                          MultiversXGasMapConfig
	MaxRetriesOnQuorumReached       uint64
	MaxRetriesOnWasTransferProposed uint64
	ProxyCacherExpirationSeconds    uint64
	ProxyRestAPIEntityType          string
	ProxyMaxNoncesDelta             int
	ProxyFinalityCheck              bool
}

// MultiversXGasMapConfig represents the gas limits for MultiversX operations
type MultiversXGasMapConfig struct {
	Sign                   uint64
	ProposeTransferBase    uint64
	ProposeTransferForEach uint64
	ProposeStatusBase      uint64
	ProposeStatusForEach   uint64
	PerformActionBase      uint64
	PerformActionForEach   uint64
}

// PeersRatingConfig will hold settings related to peers rating
type PeersRatingConfig struct {
	TopRatedCacheCapacity int
	BadRatedCacheCapacity int
}
