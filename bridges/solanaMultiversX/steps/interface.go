package steps

import (
	"context"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-solana-bridge-go/clients"
)

// Executor defines a generic bridge interface able to handle both halves of the bridge
type Executor interface {
	PrintInfo(logLevel logger.LogLevel, message string, extras ...interface{})
	MyTurnAsLeader() bool

	GetBatchFromMultiversX(ctx context.Context) (*clients.TransferBatch, error)
	StoreBatchFromMultiversX(batch *clients.TransferBatch) error
	GetStoredBatch() *clients.TransferBatch
	GetLastExecutedSolBatchIDFromMultiversX(ctx context.Context) (uint64, error)
	VerifyLastDepositNonceExecutedOnSolanaBatch(ctx context.Context) error

	GetAndStoreActionIDForProposeTransferOnMultiversX(ctx context.Context) (uint64, error)
	GetAndStoreActionIDForProposeSetStatusFromMultiversX(ctx context.Context) (uint64, error)
	GetStoredActionID() uint64

	WasTransferProposedOnMultiversX(ctx context.Context) (bool, error)
	ProposeTransferOnMultiversX(ctx context.Context) error
	ProcessMaxRetriesOnWasTransferProposedOnMultiversX() bool
	ResetRetriesOnWasTransferProposedOnMultiversX()

	WasSetStatusProposedOnMultiversX(ctx context.Context) (bool, error)
	ProposeSetStatusOnMultiversX(ctx context.Context) error

	WasActionSignedOnMultiversX(ctx context.Context) (bool, error)
	SignActionOnMultiversX(ctx context.Context) error

	ProcessQuorumReachedOnMultiversX(ctx context.Context) (bool, error)
	WasActionPerformedOnMultiversX(ctx context.Context) (bool, error)
	PerformActionOnMultiversX(ctx context.Context) error
	ResolveNewDepositsStatuses(numDeposits uint64)

	ProcessMaxQuorumRetriesOnMultiversX() bool
	ResetRetriesCountOnMultiversX()

	GetAndStoreBatchFromSolana(ctx context.Context, nonce uint64) error
	WasTransferPerformedOnSolana(ctx context.Context) (bool, error)
	SignTransferOnSolana() error
	PerformTransferOnSolana(ctx context.Context) error
	ProcessQuorumReachedOnSolana(ctx context.Context) (bool, error)
	WaitForTransferConfirmation(ctx context.Context) bool
	WaitAndReturnFinalBatchStatuses(ctx context.Context) []byte
	GetBatchStatusesFromSolana(ctx context.Context) ([]byte, error)

	ProcessMaxQuorumRetriesOnSolana() bool
	ResetRetriesCountOnSolana()
	ClearStoredP2PSignaturesForSolana()

	ValidateBatch(ctx context.Context, batch *clients.TransferBatch) (bool, error)
	CheckMultiversXClientAvailability(ctx context.Context) error
	CheckSolanaClientAvailability(ctx context.Context) error

	IsInterfaceNil() bool
}
