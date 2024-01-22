package soltomultiversx

import (
	"context"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-solana-bridge-go/bridges/solanaMultiversX/steps"
	"github.com/multiversx/mx-solana-bridge-go/core"
)

type waitForQuorumStep struct {
	bridge steps.Executor
}

// Execute will execute this step returning the next step to be executed
func (step *waitForQuorumStep) Execute(ctx context.Context) core.StepIdentifier {
	if step.bridge.ProcessMaxQuorumRetriesOnMultiversX() {
		step.bridge.PrintInfo(logger.LogDebug, "max number of retries reached, resetting counter")
		return GettingPendingBatchFromSolana
	}

	isQuorumReached, err := step.bridge.ProcessQuorumReachedOnMultiversX(ctx)
	if err != nil {
		step.bridge.PrintInfo(logger.LogError, "error while checking the quorum", "error", err)
		return GettingPendingBatchFromSolana
	}

	step.bridge.PrintInfo(logger.LogDebug, "quorum reached check", "is reached", isQuorumReached)

	if !isQuorumReached {
		return step.Identifier()
	}

	return PerformingActionID
}

// Identifier returns the step's identifier
func (step *waitForQuorumStep) Identifier() core.StepIdentifier {
	return WaitingForQuorum
}

// IsInterfaceNil returns true if there is no value under the interface
func (step *waitForQuorumStep) IsInterfaceNil() bool {
	return step == nil
}
