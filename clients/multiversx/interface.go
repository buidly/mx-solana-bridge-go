package multiversx

import (
	"context"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/core"
)

// NonceTransactionsHandler represents the interface able to handle the current nonce and the transactions resend mechanism
type NonceTransactionsHandler interface {
	GetNonce(ctx context.Context, address core.AddressHandler) (uint64, error)
	SendTransaction(ctx context.Context, tx *transaction.FrontendTransaction) (string, error)
	Close() error
}

// TokensMapper can convert a token bytes from one chain to another
type TokensMapper interface {
	ConvertToken(ctx context.Context, sourceBytes []byte) ([]byte, error)
	IsInterfaceNil() bool
}

type txHandler interface {
	SendTransactionReturnHash(ctx context.Context, builder builders.TxDataBuilder, gasLimit uint64) (string, error)
	Close() error
}

type roleProvider interface {
	IsWhitelisted(address core.AddressHandler) bool
	IsInterfaceNil() bool
}
