package soltomultiversx

import (
	"context"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-solana-bridge-go/bridges/solanaMultiversX/steps"
	"github.com/multiversx/mx-solana-bridge-go/core"
)

type performActionIDStep struct {
	bridge steps.Executor
}

// Execute will execute this step returning the next step to be executed
func (step *performActionIDStep) Execute(ctx context.Context) core.StepIdentifier {
	wasPerformed, err := step.bridge.WasActionPerformedOnMultiversX(ctx)
	if err != nil {
		step.bridge.PrintInfo(logger.LogError, "error determining if the action ID was proposed or not",
			"action ID", step.bridge.GetStoredActionID(), "error", err)
		return GettingPendingBatchFromSolana
	}

	if wasPerformed {
		step.bridge.PrintInfo(logger.LogInfo, "action ID performed",
			"action ID", step.bridge.GetStoredActionID())
		return GettingPendingBatchFromSolana
	}

	if !step.bridge.MyTurnAsLeader() {
		step.bridge.PrintInfo(logger.LogDebug, "not my turn as leader in this round")
		return step.Identifier()
	}

	err = step.bridge.PerformActionOnMultiversX(ctx)
	if err != nil {
		step.bridge.PrintInfo(logger.LogError, "error performing action ID",
			"action ID", step.bridge.GetStoredActionID(), "error", err)
		return GettingPendingBatchFromSolana
	}

	return step.Identifier()
}

// Identifier returns the step's identifier
func (step *performActionIDStep) Identifier() core.StepIdentifier {
	return PerformingActionID
}

// IsInterfaceNil returns true if there is no value under the interface
func (step *performActionIDStep) IsInterfaceNil() bool {
	return step == nil
}
