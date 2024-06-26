package multiversxtosol

import (
	"context"
	"encoding/json"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-solana-bridge-go/bridges/solanaMultiversX/steps"
	"github.com/multiversx/mx-solana-bridge-go/core"
)

type getPendingStep struct {
	bridge steps.Executor
}

// Execute will execute this step returning the next step to be executed
func (step *getPendingStep) Execute(ctx context.Context) core.StepIdentifier {
	err := step.bridge.CheckMultiversXClientAvailability(ctx)
	if err != nil {
		step.bridge.PrintInfo(logger.LogDebug, "MultiversX client unavailable", "message", err)
	}
	err = step.bridge.CheckSolanaClientAvailability(ctx)
	if err != nil {
		step.bridge.PrintInfo(logger.LogDebug, "Solana client unavailable", "message", err)
	}
	step.bridge.ResetRetriesCountOnSolana()
	step.resetCountersOnMultiversX()

	batch, err := step.bridge.GetBatchFromMultiversX(ctx)
	if err != nil {
		step.bridge.PrintInfo(logger.LogDebug, "cannot fetch MultiversX batch", "message", err)
		return step.Identifier()
	}
	if batch == nil {
		step.bridge.PrintInfo(logger.LogDebug, "no new batch found on MultiversX")
		return step.Identifier()
	}

	err = step.bridge.StoreBatchFromMultiversX(batch)
	if err != nil {
		step.bridge.PrintInfo(logger.LogError, "error storing MultiversX batch", "error", err)
		return step.Identifier()
	}

	isValid, err := step.bridge.ValidateBatch(ctx, batch)
	if err != nil {
		body, _ := json.Marshal(batch)
		step.bridge.PrintInfo(logger.LogError, "error validating MultiversX batch", "error", err, "batch", string(body))
		return step.Identifier()
	}

	if !isValid {
		step.bridge.PrintInfo(logger.LogError, "batch not valid "+batch.String())
		return step.Identifier()
	}

	step.bridge.PrintInfo(logger.LogInfo, "fetched new batch from MultiversX "+batch.String())

	wasPerformed, err := step.bridge.WasTransferPerformedOnSolana(ctx)
	if err != nil {
		step.bridge.PrintInfo(logger.LogError, "error determining if transfer was performed or not", "error", err)
		return step.Identifier()
	}
	if wasPerformed {
		step.bridge.PrintInfo(logger.LogInfo, "transfer performed")
		return WaitingTransferConfirmation
	}

	return SigningProposedTransferOnSolana
}

// Identifier returns the step's identifier
func (step *getPendingStep) Identifier() core.StepIdentifier {
	return GettingPendingBatchFromMultiversX
}

// IsInterfaceNil returns true if there is no value under the interface
func (step *getPendingStep) IsInterfaceNil() bool {
	return step == nil
}

func (step *getPendingStep) resetCountersOnMultiversX() {
	step.bridge.ResetRetriesCountOnMultiversX()
	step.bridge.ResetRetriesOnWasTransferProposedOnMultiversX()
}
