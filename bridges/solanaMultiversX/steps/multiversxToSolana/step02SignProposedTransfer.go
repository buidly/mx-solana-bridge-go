package multiversxtosol

import (
	"context"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-solana-bridge-go/bridges/solanaMultiversX/steps"
	"github.com/multiversx/mx-solana-bridge-go/core"
)

type signProposedTransferStep struct {
	bridge steps.Executor
}

// Execute will execute this step returning the next step to be executed
func (step *signProposedTransferStep) Execute(_ context.Context) core.StepIdentifier {
	storedBatch := step.bridge.GetStoredBatch()
	if storedBatch == nil {
		step.bridge.PrintInfo(logger.LogDebug, "nil batch stored")
		return GettingPendingBatchFromMultiversX
	}

	err := step.bridge.SignTransferOnSolana()
	if err != nil {
		step.bridge.PrintInfo(logger.LogError, "error signing", "batch ID", storedBatch.ID, "error", err)
		return GettingPendingBatchFromMultiversX
	}

	return WaitingForQuorumOnTransfer
}

// Identifier returns the step's identifier
func (step *signProposedTransferStep) Identifier() core.StepIdentifier {
	return SigningProposedTransferOnSolana
}

// IsInterfaceNil returns true if there is no value under the interface
func (step *signProposedTransferStep) IsInterfaceNil() bool {
	return step == nil
}
