package clients

import (
	"context"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"math/big"
)

// GasHandler defines the component able to fetch the current gas price
type GasHandler interface {
	GetCurrentGasPrice() (*big.Int, error)
	IsInterfaceNil() bool
}

// BatchValidator defines the operations for a component that can verify a batch
type BatchValidator interface {
	ValidateBatch(ctx context.Context, batch *TransferBatch) (bool, error)
	IsInterfaceNil() bool
}

// Proxy defines the behavior of a proxy able to serve MultiversX blockchain requests
type Proxy interface {
	GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error)
	SendTransaction(ctx context.Context, tx *transaction.FrontendTransaction) (string, error)
	SendTransactions(ctx context.Context, txs []*transaction.FrontendTransaction) ([]string, error)
	ExecuteVMQuery(ctx context.Context, vmRequest *data.VmValueRequest) (*data.VmValuesResponseData, error)
	GetAccount(ctx context.Context, address core.AddressHandler) (*data.Account, error)
	GetNetworkStatus(ctx context.Context, shardID uint32) (*data.NetworkStatus, error)
	GetShardOfAddress(ctx context.Context, bech32Address string) (uint32, error)
	GetHTTP(ctx context.Context, endpoint string) ([]byte, int, error)
	IsInterfaceNil() bool
}

type DecimalDiffCalculator interface {
	GetDecimalDifference(ctx context.Context, solanaTokenAddress []byte, mvxTokenAddress []byte) (int, error)
}
