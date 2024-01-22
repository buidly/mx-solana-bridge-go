package core

const (
	// WebServerOffString represents the constant used to switch off the web server
	WebServerOffString = "off"
)

const (
	// MetricNumBatches represents the metric used for counting the number of executed batches
	MetricNumBatches = "num batches"

	// MetricLastError represents the metric used to store the last encountered error
	MetricLastError = "last encountered error"

	// MetricCurrentStateMachineStep represents the metric used to store the current running machine step
	MetricCurrentStateMachineStep = "current state machine step"

	// MetricNumSolClientRequests represents the metric used to count the number of sol client requests
	MetricNumSolClientRequests = "num sol client requests"

	// MetricNumSolClientTransactions represents the metric used to count the number of sol sent transactions
	MetricNumSolClientTransactions = "num sol client transactions"

	// MetricLastQueriedSolanaBlockNumber represents the metric used to store the last sol block number that was
	// fetched from the sol client
	MetricLastQueriedSolanaBlockNumber = "sol last queried block number"

	// MetricSolanaClientStatus represents the metric used to store the status of the sol client
	MetricSolanaClientStatus = "sol client status"

	// MetricLastSolanaClientError represents the metric used to store the last encountered error from the sol client
	MetricLastSolanaClientError = "sol client last encountered error"

	// MetricLastQueriedMultiversXBlockNumber represents the metric used to store the last MultiversX block number that was
	// fetched from the MultiversX client
	MetricLastQueriedMultiversXBlockNumber = "multiversx last queried block number"

	// MetricMultiversXClientStatus represents the metric used to store the status of the MultiversX client
	MetricMultiversXClientStatus = "multiversx client status"

	// MetricLastMultiversXClientError represents the metric used to store the last encountered error from the MultiversX client
	MetricLastMultiversXClientError = "multiversx client last encountered error"

	// MetricRelayerP2PAddresses represents the metric used to store all the P2P addresses the messenger has bound to
	MetricRelayerP2PAddresses = "relayer P2P addresses"

	// MetricConnectedP2PAddresses represents the metric used to store all the P2P addresses the messenger has connected to
	MetricConnectedP2PAddresses = "connected P2P addresses"

	// MetricLastBlockNonce represents the last block nonce queried
	MetricLastBlockNonce = "last block nonce"
)

// PersistedMetrics represents the array of metrics that should be persisted
var PersistedMetrics = []string{MetricNumBatches, MetricNumSolClientRequests, MetricNumSolClientTransactions,
	MetricLastQueriedSolanaBlockNumber, MetricLastQueriedMultiversXBlockNumber, MetricSolanaClientStatus,
	MetricMultiversXClientStatus, MetricLastSolanaClientError, MetricLastMultiversXClientError, MetricLastBlockNonce}

const (
	// SolClientStatusHandlerName is the Solana client status handler name
	SolClientStatusHandlerName = "sol-client"

	// MultiversXClientStatusHandlerName is the MultiversX client status handler name
	MultiversXClientStatusHandlerName = "multiversx-client"
)
