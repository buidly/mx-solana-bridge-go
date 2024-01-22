package soltomultiversx

import (
	"context"
	solmultiversx "github.com/multiversx/mx-solana-bridge-go/bridges/solanaMultiversX"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-solana-bridge-go/bridges/solanaMultiversX/steps"
	"github.com/multiversx/mx-solana-bridge-go/core"
)

type signProposedTransferStep struct {
	bridge steps.Executor
}

// Execute will execute this step returning the next step to be executed
func (step *signProposedTransferStep) Execute(ctx context.Context) core.StepIdentifier {
	batch := step.bridge.GetStoredBatch()
	if batch == nil {
		step.bridge.PrintInfo(logger.LogDebug, "no batch found")
		return GettingPendingBatchFromSolana
	}

	actionID, err := step.bridge.GetAndStoreActionIDForProposeTransferOnMultiversX(ctx)
	if err != nil {
		step.bridge.PrintInfo(logger.LogError, "error fetching action ID", "batch ID", batch.ID, "error", err)
		return GettingPendingBatchFromSolana
	}
	if actionID == solmultiversx.InvalidActionID {
		step.bridge.PrintInfo(logger.LogError, "contract error, got invalid action ID",
			"batch ID", batch.ID, "error", err, "action ID", actionID)
		return GettingPendingBatchFromSolana
	}

	step.bridge.PrintInfo(logger.LogInfo, "fetched action ID", "action ID", actionID, "batch ID", batch.ID)

	wasSigned, err := step.bridge.WasActionSignedOnMultiversX(ctx)
	if err != nil {
		step.bridge.PrintInfo(logger.LogError, "error determining if the proposed transfer was signed or not",
			"batch ID", batch.ID, "error", err)
		return GettingPendingBatchFromSolana
	}

	if wasSigned {
		return WaitingForQuorum
	}

	err = step.bridge.SignActionOnMultiversX(ctx)
	if err != nil {
		step.bridge.PrintInfo(logger.LogError, "error signing the proposed transfer",
			"batch ID", batch.ID, "error", err)
		return GettingPendingBatchFromSolana
	}

	return WaitingForQuorum
}

// Identifier returns the step's identifier
func (step *signProposedTransferStep) Identifier() core.StepIdentifier {
	return SigningProposedTransferOnMultiversX
}

// IsInterfaceNil returns true if there is no value under the interface
func (step *signProposedTransferStep) IsInterfaceNil() bool {
	return step == nil
}
