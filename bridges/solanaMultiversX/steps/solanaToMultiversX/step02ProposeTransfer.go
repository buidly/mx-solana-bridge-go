package soltomultiversx

import (
	"context"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-solana-bridge-go/bridges/solanaMultiversX/steps"
	"github.com/multiversx/mx-solana-bridge-go/core"
)

type proposeTransferStep struct {
	bridge steps.Executor
}

// Execute will execute this step returning the next step to be executed
func (step *proposeTransferStep) Execute(ctx context.Context) core.StepIdentifier {
	batch := step.bridge.GetStoredBatch()
	if batch == nil {
		step.bridge.PrintInfo(logger.LogDebug, "no batch found")
		return GettingPendingBatchFromSolana
	}

	wasTransferProposed, err := step.bridge.WasTransferProposedOnMultiversX(ctx)
	if err != nil {
		step.bridge.PrintInfo(logger.LogError, "error determining if the batch was proposed or not on MultiversX",
			"batch ID", batch.ID, "error", err)
		return GettingPendingBatchFromSolana
	}

	if wasTransferProposed {
		return SigningProposedTransferOnMultiversX
	}

	if !step.bridge.MyTurnAsLeader() {
		step.bridge.PrintInfo(logger.LogDebug, "not my turn as leader in this round")
		return step.Identifier()
	}

	err = step.bridge.ProposeTransferOnMultiversX(ctx)
	if err != nil {
		step.bridge.PrintInfo(logger.LogError, "error proposing transfer on MultiversX",
			"batch ID", batch.ID, "error", err)
		return GettingPendingBatchFromSolana
	}

	return SigningProposedTransferOnMultiversX
}

// Identifier returns the step's identifier
func (step *proposeTransferStep) Identifier() core.StepIdentifier {
	return ProposingTransferOnMultiversX
}

// IsInterfaceNil returns true if there is no value under the interface
func (step *proposeTransferStep) IsInterfaceNil() bool {
	return step == nil
}
