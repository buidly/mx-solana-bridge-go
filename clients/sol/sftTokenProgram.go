package sol

import (
	"context"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-solana-bridge-go/clients"
	"github.com/multiversx/mx-solana-bridge-go/core"
)

type ArgsSftTokenProgram struct {
	SolanaRPC                 *rpc.Client
	SolanaClientStatusHandler core.StatusHandler
}

// sftTokenProgram represents the SftTokenProgram implementation
type sftTokenProgram struct {
	mut                       sync.RWMutex
	solanaRPC                 *rpc.Client
	solanaClientStatusHandler core.StatusHandler
}

// NewSftTokenProgram returns a new sftTokenProgram instance
func NewSftTokenProgram(args ArgsSftTokenProgram) (*sftTokenProgram, error) {
	if check.IfNilReflect(args.SolanaRPC) {
		return nil, errNilSolClient
	}
	if check.IfNil(args.SolanaClientStatusHandler) {
		return nil, clients.ErrNilStatusHandler
	}
	return &sftTokenProgram{
		solanaRPC:                 args.SolanaRPC,
		solanaClientStatusHandler: args.SolanaClientStatusHandler,
	}, nil
}

func (h *sftTokenProgram) BalanceOf(ctx context.Context, mintTokenAddress solana.PublicKey, associatedTokenAccount solana.PublicKey) (*rpc.UiTokenAmount, error) {
	h.mut.Lock()
	defer h.mut.Unlock()

	var accountInfo token.Account
	err := h.solanaRPC.GetAccountDataInto(
		context.TODO(),
		associatedTokenAccount,
		&accountInfo,
	)
	if err != nil {
		return nil, err
	}

	if !(accountInfo.Mint.Equals(mintTokenAddress)) {
		return nil, errMintAccountDidNotMatch
	}

	accountBalance, err := h.solanaRPC.GetTokenAccountBalance(
		context.TODO(),
		associatedTokenAccount,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return nil, err
	}

	h.solanaClientStatusHandler.AddIntMetric(core.MetricNumSolClientRequests, 1)

	return accountBalance.Value, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (h *sftTokenProgram) IsInterfaceNil() bool {
	return h == nil
}
