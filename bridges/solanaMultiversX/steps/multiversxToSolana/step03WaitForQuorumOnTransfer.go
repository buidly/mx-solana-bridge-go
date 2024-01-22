package multiversxtosol

import (
	"context"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-solana-bridge-go/bridges/solanaMultiversX/steps"
	"github.com/multiversx/mx-solana-bridge-go/core"
)

type waitForQuorumOnTransferStep struct {
	bridge steps.Executor
}

// Execute will execute this step returning the next step to be executed
func (step *waitForQuorumOnTransferStep) Execute(ctx context.Context) core.StepIdentifier {
	if step.bridge.ProcessMaxQuorumRetriesOnSolana() {
		step.bridge.PrintInfo(logger.LogDebug, "max number of retries reached, resetting counter")
		return GettingPendingBatchFromMultiversX
	}

	isQuorumReached, err := step.bridge.ProcessQuorumReachedOnSolana(ctx)
	if err != nil {
		step.bridge.PrintInfo(logger.LogError, "error while checking the quorum on Solana", "error", err)
		return GettingPendingBatchFromMultiversX
	}

	step.bridge.PrintInfo(logger.LogDebug, "quorum reached check", "is reached", isQuorumReached)

	if !isQuorumReached {
		return step.Identifier()
	}

	return PerformingTransfer
}

// Identifier returns the step's identifier
func (step *waitForQuorumOnTransferStep) Identifier() core.StepIdentifier {
	return WaitingForQuorumOnTransfer
}

// IsInterfaceNil returns true if there is no value under the interface
func (step *waitForQuorumOnTransferStep) IsInterfaceNil() bool {
	return step == nil
}
