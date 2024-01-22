package sol

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go"
	solanaRPC "github.com/gagliardetto/solana-go/rpc"
	"github.com/multiversx/mx-solana-bridge-go/clients/sol/contract/tokens_safe"
	"github.com/multiversx/mx-solana-bridge-go/core"
)

// ClientWrapper represents the Solana client wrapper that the sol client can rely on
type ClientWrapper interface {
	core.StatusHandler
	GetBatch(ctx context.Context, batchNonce *big.Int) (tokens_safe.Batch, error)
	GetBatchDeposits(ctx context.Context, batchNonce *big.Int) ([]tokens_safe.Deposit, error)
	GetRelayers(ctx context.Context) ([]common.Address, error)
	WasBatchExecuted(ctx context.Context, batchNonce *big.Int) (bool, error)
	ChainID(ctx context.Context) (*big.Int, error)
	BlockNumber(ctx context.Context) (uint64, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	ExecuteTransfer(opts *bind.TransactOpts, tokens []common.Address,
		recipients []common.Address, amounts []*big.Int, nonces []*big.Int, batchNonce *big.Int,
		signatures [][]byte) (*types.Transaction, error)
	Quorum(ctx context.Context) (*big.Int, error)
	GetStatusesAfterExecution(ctx context.Context, batchID *big.Int) ([]byte, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	IsPaused(ctx context.Context) (bool, error)
}

// SftTokenProgram defines the Solana ERC20 contract operations
type SftTokenProgram interface {
	BalanceOf(ctx context.Context, mintTokenAddress solana.PublicKey, associatedTokenAccount solana.PublicKey) (*solanaRPC.UiTokenAmount, error)
	IsInterfaceNil() bool
}

// Broadcaster defines the operations for a component used for communication with other peers
type Broadcaster interface {
	BroadcastSignature([]*core.EthereumSignature)
	IsInterfaceNil() bool
}

// TokensMapper can convert a token bytes from one chain to another
type TokensMapper interface {
	ConvertToken(ctx context.Context, sourceBytes []byte) ([]byte, error)
	IsInterfaceNil() bool
}

// GasHandler defines the component able to fetch the current gas price
type GasHandler interface {
	GetCurrentGasPrice() (*big.Int, error)
	IsInterfaceNil() bool
}

// SignaturesHolder defines the operations for a component that can hold and manage signatures
type SignaturesHolder interface {
	Signatures(messageHash []byte) [][]byte
	ClearStoredSignatures()
	IsInterfaceNil() bool
}

type erc20ContractWrapper interface {
	BalanceOf(ctx context.Context, account common.Address) (*big.Int, error)
	IsInterfaceNil() bool
}
