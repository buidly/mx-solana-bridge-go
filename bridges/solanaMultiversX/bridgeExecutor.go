package solmultiversx

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-solana-bridge-go/clients"
	"github.com/multiversx/mx-solana-bridge-go/core"
)

// splits - represent the number of times we split the maximum interval
// we wait for the transfer confirmation on Solana
const splits = 3

const minRetries = 1

// ArgsBridgeExecutor is the arguments DTO struct used in both bridges
type ArgsBridgeExecutor struct {
	Log                          logger.Logger
	TopologyProvider             TopologyProvider
	MultiversXClient             MultiversXClient
	SolanaClient                 SolanaClient
	TimeForWaitOnSolana          time.Duration
	StatusHandler                core.StatusHandler
	SignaturesHolder             SignaturesHolder
	BatchValidator               clients.BatchValidator
	MaxQuorumRetriesOnSolana     uint64
	MaxQuorumRetriesOnMultiversX uint64
	MaxRestriesOnWasProposed     uint64
}

type bridgeExecutor struct {
	log                          logger.Logger
	topologyProvider             TopologyProvider
	multiversXClient             MultiversXClient
	solanaClient                 SolanaClient
	timeForWaitOnSolana          time.Duration
	statusHandler                core.StatusHandler
	sigsHolder                   SignaturesHolder
	batchValidator               clients.BatchValidator
	maxQuorumRetriesOnSolana     uint64
	maxQuorumRetriesOnMultiversX uint64
	maxRetriesOnWasProposed      uint64

	batch                     *clients.TransferBatch
	actionID                  uint64
	depositMsgHashMap         map[uint64]common.Hash
	quorumRetriesOnSolana     uint64
	quorumRetriesOnMultiversX uint64
	retriesOnWasProposed      uint64
}

// NewBridgeExecutor creates a bridge executor, which can be used for both half-bridges
func NewBridgeExecutor(args ArgsBridgeExecutor) (*bridgeExecutor, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	executor := createBridgeExecutor(args)
	return executor, nil
}

func checkArgs(args ArgsBridgeExecutor) error {
	if check.IfNil(args.Log) {
		return ErrNilLogger
	}
	if check.IfNil(args.MultiversXClient) {
		return ErrNilMultiversXClient
	}
	if check.IfNil(args.SolanaClient) {
		return ErrNilSolanaClient
	}
	if check.IfNil(args.TopologyProvider) {
		return ErrNilTopologyProvider
	}
	if check.IfNil(args.StatusHandler) {
		return ErrNilStatusHandler
	}
	if args.TimeForWaitOnSolana < durationLimit {
		return ErrInvalidDuration
	}
	if check.IfNil(args.SignaturesHolder) {
		return ErrNilSignaturesHolder
	}
	if check.IfNil(args.BatchValidator) {
		return ErrNilBatchValidator
	}
	if args.MaxQuorumRetriesOnSolana < minRetries {
		return fmt.Errorf("%w for args.MaxQuorumRetriesOnSolana, got: %d, minimum: %d",
			clients.ErrInvalidValue, args.MaxQuorumRetriesOnSolana, minRetries)
	}
	if args.MaxQuorumRetriesOnMultiversX < minRetries {
		return fmt.Errorf("%w for args.MaxQuorumRetriesOnMultiversX, got: %d, minimum: %d",
			clients.ErrInvalidValue, args.MaxQuorumRetriesOnMultiversX, minRetries)
	}
	if args.MaxRestriesOnWasProposed < minRetries {
		return fmt.Errorf("%w for args.MaxRestriesOnWasProposed, got: %d, minimum: %d",
			clients.ErrInvalidValue, args.MaxRestriesOnWasProposed, minRetries)
	}
	return nil
}

func createBridgeExecutor(args ArgsBridgeExecutor) *bridgeExecutor {
	return &bridgeExecutor{
		log:                          args.Log,
		multiversXClient:             args.MultiversXClient,
		solanaClient:                 args.SolanaClient,
		topologyProvider:             args.TopologyProvider,
		statusHandler:                args.StatusHandler,
		timeForWaitOnSolana:          args.TimeForWaitOnSolana,
		sigsHolder:                   args.SignaturesHolder,
		batchValidator:               args.BatchValidator,
		maxQuorumRetriesOnSolana:     args.MaxQuorumRetriesOnSolana,
		maxQuorumRetriesOnMultiversX: args.MaxQuorumRetriesOnMultiversX,
		maxRetriesOnWasProposed:      args.MaxRestriesOnWasProposed,
	}
}

// PrintInfo will print the provided data through the inner logger instance
func (executor *bridgeExecutor) PrintInfo(logLevel logger.LogLevel, message string, extras ...interface{}) {
	executor.log.Log(logLevel, message, extras...)

	switch logLevel {
	case logger.LogWarning, logger.LogError:
		executor.setExecutionMessageInStatusHandler(logLevel, message, extras...)
	}
}

func (executor *bridgeExecutor) setExecutionMessageInStatusHandler(level logger.LogLevel, message string, extras ...interface{}) {
	msg := fmt.Sprintf("%s: %s", level, message)
	for i := 0; i < len(extras)-1; i += 2 {
		msg += fmt.Sprintf(" %s = %s", convertObjectToString(extras[i]), convertObjectToString(extras[i+1]))
	}

	executor.statusHandler.SetStringMetric(core.MetricLastError, msg)
}

// MyTurnAsLeader returns true if the current relayer node is the leader
func (executor *bridgeExecutor) MyTurnAsLeader() bool {
	return executor.topologyProvider.MyTurnAsLeader()
}

// GetBatchFromMultiversX fetches the pending batch from MultiversX
func (executor *bridgeExecutor) GetBatchFromMultiversX(ctx context.Context) (*clients.TransferBatch, error) {
	batch, err := executor.multiversXClient.GetPending(ctx)
	if err == nil {
		executor.statusHandler.SetIntMetric(core.MetricNumBatches, int(batch.ID)-1)
	}
	return batch, err
}

// StoreBatchFromMultiversX saves the pending batch from MultiversX
func (executor *bridgeExecutor) StoreBatchFromMultiversX(batch *clients.TransferBatch) error {
	if batch == nil {
		return ErrNilBatch
	}

	executor.batch = batch
	return nil
}

// GetStoredBatch returns the stored batch
func (executor *bridgeExecutor) GetStoredBatch() *clients.TransferBatch {
	return executor.batch
}

// GetLastExecutedSolBatchIDFromMultiversX returns the last executed batch ID that is stored on the MultiversX SC
func (executor *bridgeExecutor) GetLastExecutedSolBatchIDFromMultiversX(ctx context.Context) (uint64, error) {
	batchID, err := executor.multiversXClient.GetLastExecutedSolBatchID(ctx)
	if err == nil {
		executor.statusHandler.SetIntMetric(core.MetricNumBatches, int(batchID))
	}
	return batchID, err
}

// VerifyLastDepositNonceExecutedOnSolanaBatch will check the deposit nonces from the fetched batch from Solana client
func (executor *bridgeExecutor) VerifyLastDepositNonceExecutedOnSolanaBatch(ctx context.Context) error {
	if executor.batch == nil {
		return ErrNilBatch
	}

	lastNonce, err := executor.multiversXClient.GetLastExecutedSolTxID(ctx)
	if err != nil {
		return err
	}

	return executor.verifyDepositNonces(lastNonce)
}

func (executor *bridgeExecutor) verifyDepositNonces(lastNonce uint64) error {
	startNonce := lastNonce + 1
	for _, dt := range executor.batch.Deposits {
		if dt.Nonce != startNonce {
			return fmt.Errorf("%w for deposit %s, expected: %d", ErrInvalidDepositNonce, dt.String(), startNonce)
		}

		startNonce++
	}

	return nil
}

// GetAndStoreActionIDForProposeTransferOnMultiversX fetches the action ID for ProposeTransfer by using the stored batch. Stores the action ID and returns it
func (executor *bridgeExecutor) GetAndStoreActionIDForProposeTransferOnMultiversX(ctx context.Context) (uint64, error) {
	if executor.batch == nil {
		return InvalidActionID, ErrNilBatch
	}

	actionID, err := executor.multiversXClient.GetActionIDForProposeTransfer(ctx, executor.batch)
	if err != nil {
		return InvalidActionID, err
	}

	executor.actionID = actionID

	return actionID, nil
}

// GetAndStoreActionIDForProposeSetStatusFromMultiversX fetches the action ID for SetStatus by using the stored batch. Stores the action ID and returns it
func (executor *bridgeExecutor) GetAndStoreActionIDForProposeSetStatusFromMultiversX(ctx context.Context) (uint64, error) {
	if executor.batch == nil {
		return InvalidActionID, ErrNilBatch
	}

	actionID, err := executor.multiversXClient.GetActionIDForSetStatusOnPendingTransfer(ctx, executor.batch)
	if err != nil {
		return InvalidActionID, err
	}

	executor.actionID = actionID

	return actionID, nil
}

// GetStoredActionID returns the stored action ID
func (executor *bridgeExecutor) GetStoredActionID() uint64 {
	return executor.actionID
}

// WasTransferProposedOnMultiversX checks if the transfer was proposed on MultiversX
func (executor *bridgeExecutor) WasTransferProposedOnMultiversX(ctx context.Context) (bool, error) {
	if executor.batch == nil {
		return false, ErrNilBatch
	}

	return executor.multiversXClient.WasProposedTransfer(ctx, executor.batch)
}

// ProposeTransferOnMultiversX propose the transfer on MultiversX
func (executor *bridgeExecutor) ProposeTransferOnMultiversX(ctx context.Context) error {
	if executor.batch == nil {
		return ErrNilBatch
	}

	hash, err := executor.multiversXClient.ProposeTransfer(ctx, executor.batch)
	if err != nil {
		return err
	}

	executor.log.Info("proposed transfer", "hash", hash,
		"batch ID", executor.batch.ID, "action ID", executor.actionID)

	return nil
}

// ProcessMaxRetriesOnWasTransferProposedOnMultiversX checks if the retries on MultiversX were reached and increments the counter
func (executor *bridgeExecutor) ProcessMaxRetriesOnWasTransferProposedOnMultiversX() bool {
	if executor.retriesOnWasProposed < executor.maxRetriesOnWasProposed {
		executor.retriesOnWasProposed++
		return false
	}

	return true
}

// ResetRetriesOnWasTransferProposedOnMultiversX resets the number of retries on was transfer proposed
func (executor *bridgeExecutor) ResetRetriesOnWasTransferProposedOnMultiversX() {
	executor.retriesOnWasProposed = 0
}

// WasSetStatusProposedOnMultiversX checks if set status was proposed on MultiversX
func (executor *bridgeExecutor) WasSetStatusProposedOnMultiversX(ctx context.Context) (bool, error) {
	if executor.batch == nil {
		return false, ErrNilBatch
	}

	return executor.multiversXClient.WasProposedSetStatus(ctx, executor.batch)
}

// ProposeSetStatusOnMultiversX propose set status on MultiversX
func (executor *bridgeExecutor) ProposeSetStatusOnMultiversX(ctx context.Context) error {
	if executor.batch == nil {
		return ErrNilBatch
	}

	hash, err := executor.multiversXClient.ProposeSetStatus(ctx, executor.batch)
	if err != nil {
		return err
	}

	executor.log.Info("proposed set status", "hash", hash,
		"batch ID", executor.batch.ID)

	return nil
}

// WasActionSignedOnMultiversX returns true if the current relayer already signed the action
func (executor *bridgeExecutor) WasActionSignedOnMultiversX(ctx context.Context) (bool, error) {
	return executor.multiversXClient.WasSigned(ctx, executor.actionID)
}

// SignActionOnMultiversX calls the MultiversX client to generate and send the signature
func (executor *bridgeExecutor) SignActionOnMultiversX(ctx context.Context) error {
	hash, err := executor.multiversXClient.Sign(ctx, executor.actionID)
	if err != nil {
		return err
	}

	executor.log.Info("signed proposed transfer", "hash", hash, "action ID", executor.actionID)

	return nil
}

// ProcessQuorumReachedOnMultiversX returns true if the proposed transfer reached the set quorum
func (executor *bridgeExecutor) ProcessQuorumReachedOnMultiversX(ctx context.Context) (bool, error) {
	return executor.multiversXClient.QuorumReached(ctx, executor.actionID)
}

// WaitForTransferConfirmation waits for the confirmation of a transfer
func (executor *bridgeExecutor) WaitForTransferConfirmation(ctx context.Context) bool {
	wasPerformed := false
	executor.log.Info("Waiting for batch execution to be final", "batchID", executor.batch.ID)
	for i := 0; i < splits && !wasPerformed; i++ {
		if executor.waitWithContextSucceeded(ctx) {
			wasExecuted, lastUpdatedSlotNumber, err := executor.solanaClient.WasExecuted(ctx, executor.batch.ID)
			if err != nil {
				executor.log.Error("Could not check if batch was executed", "batchID", executor.batch.ID, "error", err)
				continue
			}
			lastFinalizedBlock, err := executor.solanaClient.GetLastFinalizedSlotNumber(ctx)
			if err != nil {
				executor.log.Error("Could not fetch last finalized block from Solana", "batchID", executor.batch.ID, "error", err)
				continue
			}
			isFinalized := lastFinalizedBlock > lastUpdatedSlotNumber
			wasPerformed = wasExecuted && isFinalized
			if !isFinalized {
				executor.log.Info("Batch execution is not final", "lastFinalizedBlock", lastFinalizedBlock, "lastUpdatedSlotNumber", lastUpdatedSlotNumber)
			}
		}
	}
	return wasPerformed
}

// WaitAndReturnFinalBatchStatuses waits for the statuses to be final
func (executor *bridgeExecutor) WaitAndReturnFinalBatchStatuses(ctx context.Context) []byte {
	for i := 0; i < splits; i++ {
		if !executor.waitWithContextSucceeded(ctx) {
			return nil
		}

		statuses, err := executor.GetBatchStatusesFromSolana(ctx)
		if err != nil {
			executor.log.Debug("got message while fetching batch statuses", "message", err)
			continue
		}
		if len(statuses) == 0 {
			executor.log.Debug("no status available")
			continue
		}

		executor.log.Debug("bridgeExecutor.WaitAndReturnFinalBatchStatuses", "statuses", statuses)
		return statuses
	}

	return nil
}

func (executor *bridgeExecutor) waitWithContextSucceeded(ctx context.Context) bool {
	timer := time.NewTimer(executor.timeForWaitOnSolana / splits)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		executor.log.Debug("closing due to context expiration")
		return false
	case <-timer.C:
		return true
	}
}

// GetBatchStatusesFromSolana gets statuses for the batch
func (executor *bridgeExecutor) GetBatchStatusesFromSolana(ctx context.Context) ([]byte, error) {
	if executor.batch == nil {
		return nil, ErrNilBatch
	}

	statuses, err := executor.solanaClient.GetTransactionsStatuses(ctx, executor.batch)
	if err != nil {
		return nil, err
	}

	return statuses, nil
}

// WasActionPerformedOnMultiversX returns true if the action was already performed
func (executor *bridgeExecutor) WasActionPerformedOnMultiversX(ctx context.Context) (bool, error) {
	return executor.multiversXClient.WasExecuted(ctx, executor.actionID)
}

// PerformActionOnMultiversX sends the perform-action transaction on the MultiversX chain
func (executor *bridgeExecutor) PerformActionOnMultiversX(ctx context.Context) error {
	if executor.batch == nil {
		return ErrNilBatch
	}

	hash, err := executor.multiversXClient.PerformAction(ctx, executor.actionID, executor.batch)
	if err != nil {
		return err
	}

	executor.log.Info("sent perform action transaction", "hash", hash,
		"batch ID", executor.batch.ID, "action ID", executor.actionID)

	return nil
}

// ResolveNewDepositsStatuses resolves the new deposits statuses for batch
func (executor *bridgeExecutor) ResolveNewDepositsStatuses(numDeposits uint64) {
	executor.batch.ResolveNewDeposits(int(numDeposits))
}

// ProcessMaxQuorumRetriesOnMultiversX checks if the retries on MultiversX were reached and increments the counter
func (executor *bridgeExecutor) ProcessMaxQuorumRetriesOnMultiversX() bool {
	if executor.quorumRetriesOnMultiversX < executor.maxQuorumRetriesOnMultiversX {
		executor.quorumRetriesOnMultiversX++
		return false
	}

	return true
}

// ResetRetriesCountOnMultiversX resets the number of retries on MultiversX
func (executor *bridgeExecutor) ResetRetriesCountOnMultiversX() {
	executor.quorumRetriesOnMultiversX = 0
}

// GetAndStoreBatchFromSolana fetches and stores the batch from the sol client
func (executor *bridgeExecutor) GetAndStoreBatchFromSolana(ctx context.Context, nonce uint64) error {
	batch, err := executor.solanaClient.GetBatch(ctx, nonce)
	if err != nil {
		return err
	}
	lastFinalizedBlock, err := executor.solanaClient.GetLastFinalizedSlotNumber(ctx)
	if err != nil {
		return err
	}

	if batch.LastUpdatedSlotNumber > lastFinalizedBlock {
		return ErrBatchNotFinalized
	}

	safeSettings, err := executor.solanaClient.GetSafeSettings()
	if err != nil {
		return err
	}

	isBatchCompleted := len(batch.Deposits) > 0 &&
		(batch.SlotNumber+uint64(safeSettings.BatchSettings.BatchSlotLimit) < lastFinalizedBlock ||
			uint16(len(batch.Deposits)) >= safeSettings.BatchSettings.BatchSize)
	isBatchInvalid := batch.ID != nonce || !isBatchCompleted
	if isBatchInvalid {
		return fmt.Errorf("%w, requested nonce: %d, fetched nonce: %d, num deposits: %d",
			ErrBatchNotFound, nonce, batch.ID, len(batch.Deposits))
	}

	executor.batch = batch

	return nil
}

// WasTransferPerformedOnSolana returns true if the batch was performed on Solana
func (executor *bridgeExecutor) WasTransferPerformedOnSolana(ctx context.Context) (bool, error) {
	if executor.batch == nil {
		return false, ErrNilBatch
	}
	wasExecuted, _, err := executor.solanaClient.WasExecuted(ctx, executor.batch.ID)
	return wasExecuted, err
}

// SignTransferOnSolana generates the message hash for batch and broadcast the signature
func (executor *bridgeExecutor) SignTransferOnSolana() error {
	if executor.batch == nil {
		return ErrNilBatch
	}

	depositHashesMap := make(map[uint64]common.Hash)
	for _, dt := range executor.batch.Deposits {

		hash, err := executor.solanaClient.GenerateMessageHash(
			solana.PublicKeyFromBytes(dt.ToBytes),
			solana.PublicKeyFromBytes(dt.ConvertedTokenBytes),
			big.NewInt(0).Set(dt.AmountAdjustedToDecimals),
			big.NewInt(0).SetUint64(dt.Nonce),
			executor.batch.ID,
		)
		if err != nil {
			return err
		}
		depositHashesMap[dt.Nonce] = hash
	}

	executor.log.Info("generated message hash on Solana for each deposit ",
		"batch ID", executor.batch.ID, "deposits", len(executor.batch.Deposits))

	executor.depositMsgHashMap = depositHashesMap

	executor.solanaClient.BroadcastSignatureForMessageHashes(depositHashesMap)

	return nil
}

// PerformTransferOnSolana transfers a batch to Solana
func (executor *bridgeExecutor) PerformTransferOnSolana(ctx context.Context) error {
	if executor.batch == nil {
		return ErrNilBatch
	}

	quorumSize, err := executor.solanaClient.GetQuorumSize(ctx)
	if err != nil {
		return err
	}

	executor.log.Debug("fetched quorum size", "quorum", quorumSize.Int64())

	hash, err := executor.solanaClient.ExecuteTransfer(ctx, executor.depositMsgHashMap, executor.batch, int(quorumSize.Int64()))
	if err != nil {
		return err
	}

	executor.log.Info("sent execute transfer for each remaining deposit", "hashes", hash,
		"batch ID", executor.batch.ID)

	return nil
}

// ProcessQuorumReachedOnSolana returns true if the proposed transfer reached the set quorum
func (executor *bridgeExecutor) ProcessQuorumReachedOnSolana(ctx context.Context) (bool, error) {
	return executor.solanaClient.IsQuorumReached(ctx, executor.batch, executor.depositMsgHashMap)
}

// ProcessMaxQuorumRetriesOnSolana checks if the retries on Solana were reached and increments the counter
func (executor *bridgeExecutor) ProcessMaxQuorumRetriesOnSolana() bool {
	if executor.quorumRetriesOnSolana < executor.maxQuorumRetriesOnSolana {
		executor.quorumRetriesOnSolana++
		return false
	}

	return true
}

// ResetRetriesCountOnSolana resets the number of retries on Solana
func (executor *bridgeExecutor) ResetRetriesCountOnSolana() {
	executor.quorumRetriesOnSolana = 0
}

// ClearStoredP2PSignaturesForSolana deletes all stored P2P signatures used for Solana client
func (executor *bridgeExecutor) ClearStoredP2PSignaturesForSolana() {
	executor.sigsHolder.ClearStoredSignatures()
	executor.log.Info("cleared stored P2P signatures")
}

// ValidateBatch returns true if the given batch is validated on microservice side
func (executor *bridgeExecutor) ValidateBatch(ctx context.Context, batch *clients.TransferBatch) (bool, error) {
	return executor.batchValidator.ValidateBatch(ctx, batch)
}

// CheckMultiversXClientAvailability trigger a self availability check for the MultiversX client
func (executor *bridgeExecutor) CheckMultiversXClientAvailability(ctx context.Context) error {
	return executor.multiversXClient.CheckClientAvailability(ctx)
}

// CheckSolanaClientAvailability trigger a self availability check for the Solana client
func (executor *bridgeExecutor) CheckSolanaClientAvailability(ctx context.Context) error {
	return executor.solanaClient.CheckClientAvailability(ctx)
}

// IsInterfaceNil returns true if there is no value under the interface
func (executor *bridgeExecutor) IsInterfaceNil() bool {
	return executor == nil
}
