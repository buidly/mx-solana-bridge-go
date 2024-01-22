package soltomultiversx

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
	step.bridge.ResetRetriesCountOnMultiversX()
	lastSolBatchExecuted, err := step.bridge.GetLastExecutedSolBatchIDFromMultiversX(ctx)
	if err != nil {
		step.bridge.PrintInfo(logger.LogError, "error fetching last executed sol batch ID", "error", err)
		return step.Identifier()
	}

	err = step.bridge.GetAndStoreBatchFromSolana(ctx, lastSolBatchExecuted+1)
	if err != nil {
		step.bridge.PrintInfo(logger.LogDebug, "cannot fetch sol batch", "batch ID", lastSolBatchExecuted+1, "message", err)
		return step.Identifier()
	}

	batch := step.bridge.GetStoredBatch()
	if batch == nil {
		step.bridge.PrintInfo(logger.LogDebug, "no new batch found on sol", "last executed on MultiversX", lastSolBatchExecuted)
		return step.Identifier()
	}

	isValid, err := step.bridge.ValidateBatch(ctx, batch)
	if err != nil {
		body, _ := json.Marshal(batch)
		step.bridge.PrintInfo(logger.LogError, "error validating Solana batch", "error", err, "batch", string(body))
		return step.Identifier()
	}

	if !isValid {
		step.bridge.PrintInfo(logger.LogError, "batch not valid "+batch.String())
		return step.Identifier()
	}

	step.bridge.PrintInfo(logger.LogInfo, "fetched new batch from Solana "+batch.String())

	err = step.bridge.VerifyLastDepositNonceExecutedOnSolanaBatch(ctx)
	if err != nil {
		step.bridge.PrintInfo(logger.LogError, "verification failed on the new batch from Solana", "batch ID", lastSolBatchExecuted+1, "error", err)
		return step.Identifier()
	}

	return ProposingTransferOnMultiversX
}

// Identifier returns the step's identifier
func (step *getPendingStep) Identifier() core.StepIdentifier {
	return GettingPendingBatchFromSolana
}

// IsInterfaceNil returns true if there is no value under the interface
func (step *getPendingStep) IsInterfaceNil() bool {
	return step == nil
}
