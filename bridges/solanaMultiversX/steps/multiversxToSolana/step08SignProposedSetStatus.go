package multiversxtosol

import (
	"context"
	solmultiversx "github.com/multiversx/mx-solana-bridge-go/bridges/solanaMultiversX"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-solana-bridge-go/bridges/solanaMultiversX/steps"
	"github.com/multiversx/mx-solana-bridge-go/core"
)

type signProposedSetStatusStep struct {
	bridge steps.Executor
}

// Execute will execute this step returning the next step to be executed
func (step *signProposedSetStatusStep) Execute(ctx context.Context) core.StepIdentifier {
	storedBatch := step.bridge.GetStoredBatch()
	if storedBatch == nil {
		step.bridge.PrintInfo(logger.LogDebug, "nil stored batch")
		return GettingPendingBatchFromMultiversX
	}

	actionID, err := step.bridge.GetAndStoreActionIDForProposeSetStatusFromMultiversX(ctx)
	if err != nil {
		step.bridge.PrintInfo(logger.LogError, "error fetching action ID", "batch ID", storedBatch.ID, "error", err)
		return GettingPendingBatchFromMultiversX
	}
	if actionID == solmultiversx.InvalidActionID {
		step.bridge.PrintInfo(logger.LogError, "contract error, got invalid action ID",
			"batch ID", storedBatch.ID, "error", err, "action ID", actionID)
		return GettingPendingBatchFromMultiversX
	}

	step.bridge.PrintInfo(logger.LogInfo, "fetched action ID", "action ID", actionID, "batch ID", storedBatch.ID)

	wasSigned, err := step.bridge.WasActionSignedOnMultiversX(ctx)
	if err != nil {
		step.bridge.PrintInfo(logger.LogError, "error determining if the proposed transfer was signed or not",
			"batch ID", storedBatch.ID, "error", err)
		return GettingPendingBatchFromMultiversX
	}

	if wasSigned {
		return WaitingForQuorumOnSetStatus
	}

	err = step.bridge.SignActionOnMultiversX(ctx)
	if err != nil {
		step.bridge.PrintInfo(logger.LogError, "error signing the proposed transfer",
			"batch ID", storedBatch.ID, "error", err)
		return GettingPendingBatchFromMultiversX
	}

	return WaitingForQuorumOnSetStatus
}

// Identifier returns the step's identifier
func (step *signProposedSetStatusStep) Identifier() core.StepIdentifier {
	return SigningProposedSetStatusOnMultiversX
}

// IsInterfaceNil returns true if there is no value under the interface
func (step *signProposedSetStatusStep) IsInterfaceNil() bool {
	return step == nil
}
