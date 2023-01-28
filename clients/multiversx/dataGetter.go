package multiversx

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"sync"

	"github.com/multiversx/mx-bridge-eth-go/clients"
	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
)

const (
	okCodeAfterExecution                                      = "ok"
	internalError                                             = "internal error"
	getCurrentTxBatchFuncName                                 = "getCurrentTxBatch"
	wasTransferActionProposedFuncName                         = "wasTransferActionProposed"
	wasActionExecutedFuncName                                 = "wasActionExecuted"
	getActionIdForTransferBatchFuncName                       = "getActionIdForTransferBatch"
	wasSetCurrentTransactionBatchStatusActionProposedFuncName = "wasSetCurrentTransactionBatchStatusActionProposed"
	getStatusesAfterExecutionFuncName                         = "getStatusesAfterExecution"
	getActionIdForSetCurrentTransactionBatchStatusFuncName    = "getActionIdForSetCurrentTransactionBatchStatus"
	getTokenIdForErc20AddressFuncName                         = "getTokenIdForErc20Address"
	getErc20AddressForTokenIdFuncName                         = "getErc20AddressForTokenId"
	quorumReachedFuncName                                     = "quorumReached"
	getLastExecutedEthBatchIdFuncName                         = "getLastExecutedEthBatchId"
	getLastExecutedEthTxId                                    = "getLastExecutedEthTxId"
	signedFuncName                                            = "signed"
	getAllStakedRelayersFuncName                              = "getAllStakedRelayers"
	isPausedFuncName                                          = "isPaused"
)

// ArgsDataGetter is the arguments DTO used in the NewDataGetter constructor
type ArgsDataGetter struct {
	MultisigContractAddress core.AddressHandler
	RelayerAddress          core.AddressHandler
	Proxy                   Proxy
	Log                     logger.Logger
}

type multiversXClientDataGetter struct {
	multisigContractAddress core.AddressHandler
	relayerAddress          core.AddressHandler
	proxy                   Proxy
	log                     logger.Logger
	mutNodeStatus           sync.Mutex
	wasShardIDFetched       bool
	shardID                 uint32
}

// NewDataGetter creates a new instance of the dataGetter type
func NewDataGetter(args ArgsDataGetter) (*multiversXClientDataGetter, error) {
	if check.IfNil(args.Log) {
		return nil, errNilLogger
	}
	if check.IfNil(args.Proxy) {
		return nil, errNilProxy
	}
	if check.IfNil(args.RelayerAddress) {
		return nil, fmt.Errorf("%w for the RelayerAddress argument", errNilAddressHandler)
	}
	if check.IfNil(args.MultisigContractAddress) {
		return nil, fmt.Errorf("%w for the MultisigContractAddress argument", errNilAddressHandler)
	}

	return &multiversXClientDataGetter{
		multisigContractAddress: args.MultisigContractAddress,
		relayerAddress:          args.RelayerAddress,
		proxy:                   args.Proxy,
		log:                     args.Log,
	}, nil
}

// ExecuteQueryReturningBytes will try to execute the provided query and return the result as slice of byte slices
func (dg *multiversXClientDataGetter) ExecuteQueryReturningBytes(ctx context.Context, request *data.VmValueRequest) ([][]byte, error) {
	if request == nil {
		return nil, errNilRequest
	}

	response, err := dg.proxy.ExecuteVMQuery(ctx, request)
	if err != nil {
		dg.log.Error("got error on VMQuery", "FuncName", request.FuncName,
			"Args", request.Args, "SC address", request.Address, "Caller", request.CallerAddr, "error", err)
		return nil, err
	}
	dg.log.Debug("executed VMQuery", "FuncName", request.FuncName,
		"Args", request.Args, "SC address", request.Address, "Caller", request.CallerAddr,
		"response.ReturnCode", response.Data.ReturnCode,
		"response.ReturnData", fmt.Sprintf("%+v", response.Data.ReturnData))
	if response.Data.ReturnCode != okCodeAfterExecution {
		return nil, NewQueryResponseError(
			response.Data.ReturnCode,
			response.Data.ReturnMessage,
			request.FuncName,
			request.Address,
			request.Args...,
		)
	}
	return response.Data.ReturnData, nil
}

// GetCurrentNonce will get from the shard containing the multisig contract the latest block's nonce
func (dg *multiversXClientDataGetter) GetCurrentNonce(ctx context.Context) (uint64, error) {
	shardID, err := dg.getShardID(ctx)
	if err != nil {
		return 0, err
	}

	nodeStatus, err := dg.proxy.GetNetworkStatus(ctx, shardID)
	if err != nil {
		return 0, err
	}
	if nodeStatus == nil {
		return 0, errNilNodeStatusResponse
	}

	return nodeStatus.Nonce, nil
}

func (dg *multiversXClientDataGetter) getShardID(ctx context.Context) (uint32, error) {
	dg.mutNodeStatus.Lock()
	defer dg.mutNodeStatus.Unlock()

	if dg.wasShardIDFetched {
		return dg.shardID, nil
	}

	var err error
	dg.shardID, err = dg.proxy.GetShardOfAddress(ctx, dg.multisigContractAddress.AddressAsBech32String())
	if err == nil {
		dg.wasShardIDFetched = true
	}

	return dg.shardID, err
}

// ExecuteQueryReturningBool will try to execute the provided query and return the result as bool
func (dg *multiversXClientDataGetter) ExecuteQueryReturningBool(ctx context.Context, request *data.VmValueRequest) (bool, error) {
	response, err := dg.ExecuteQueryReturningBytes(ctx, request)
	if err != nil {
		return false, err
	}

	if len(response) == 0 {
		return false, nil
	}

	return dg.parseBool(response[0], request.FuncName, request.Address, request.Args...)
}

func (dg *multiversXClientDataGetter) parseBool(buff []byte, funcName string, address string, args ...string) (bool, error) {
	if len(buff) == 0 {
		return false, nil
	}

	result, err := strconv.ParseBool(fmt.Sprintf("%d", buff[0]))
	if err != nil {
		return false, NewQueryResponseError(
			internalError,
			fmt.Sprintf("error converting the received bytes to bool, %s", err.Error()),
			funcName,
			address,
			args...,
		)
	}

	return result, nil
}

// ExecuteQueryReturningUint64 will try to execute the provided query and return the result as uint64
func (dg *multiversXClientDataGetter) ExecuteQueryReturningUint64(ctx context.Context, request *data.VmValueRequest) (uint64, error) {
	response, err := dg.ExecuteQueryReturningBytes(ctx, request)
	if err != nil {
		return 0, err
	}

	if len(response) == 0 {
		return 0, nil
	}
	if len(response[0]) == 0 {
		return 0, nil
	}

	num, err := parseUInt64FromByteSlice(response[0])
	if err != nil {
		return 0, NewQueryResponseError(
			internalError,
			err.Error(),
			request.FuncName,
			request.Address,
			request.Args...,
		)
	}

	return num, nil
}

func parseUInt64FromByteSlice(bytes []byte) (uint64, error) {
	num := big.NewInt(0).SetBytes(bytes)
	if !num.IsUint64() {
		return 0, errNotUint64Bytes
	}

	return num.Uint64(), nil
}

func (dg *multiversXClientDataGetter) executeQueryFromBuilder(ctx context.Context, builder builders.VMQueryBuilder) ([][]byte, error) {
	vmValuesRequest, err := builder.ToVmValueRequest()
	if err != nil {
		return nil, err
	}

	return dg.ExecuteQueryReturningBytes(ctx, vmValuesRequest)
}

func (dg *multiversXClientDataGetter) executeQueryUint64FromBuilder(ctx context.Context, builder builders.VMQueryBuilder) (uint64, error) {
	vmValuesRequest, err := builder.ToVmValueRequest()
	if err != nil {
		return 0, err
	}

	return dg.ExecuteQueryReturningUint64(ctx, vmValuesRequest)
}

func (dg *multiversXClientDataGetter) executeQueryBoolFromBuilder(ctx context.Context, builder builders.VMQueryBuilder) (bool, error) {
	vmValuesRequest, err := builder.ToVmValueRequest()
	if err != nil {
		return false, err
	}

	return dg.ExecuteQueryReturningBool(ctx, vmValuesRequest)
}

func (dg *multiversXClientDataGetter) createDefaultVmQueryBuilder() builders.VMQueryBuilder {
	return builders.NewVMQueryBuilder().Address(dg.multisigContractAddress).CallerAddress(dg.relayerAddress)
}

// GetCurrentBatchAsDataBytes will assemble a builder and query the proxy for the current pending batch
func (dg *multiversXClientDataGetter) GetCurrentBatchAsDataBytes(ctx context.Context) ([][]byte, error) {
	builder := dg.createDefaultVmQueryBuilder()
	builder.Function(getCurrentTxBatchFuncName)

	return dg.executeQueryFromBuilder(ctx, builder)
}

// GetTokenIdForErc20Address will assemble a builder and query the proxy for a token id given a specific erc20 address
func (dg *multiversXClientDataGetter) GetTokenIdForErc20Address(ctx context.Context, erc20Address []byte) ([][]byte, error) {
	builder := dg.createDefaultVmQueryBuilder()
	builder.Function(getTokenIdForErc20AddressFuncName)
	builder.ArgBytes(erc20Address)

	return dg.executeQueryFromBuilder(ctx, builder)
}

// GetERC20AddressForTokenId will assemble a builder and query the proxy for an erc20 address given a specific token id
func (dg *multiversXClientDataGetter) GetERC20AddressForTokenId(ctx context.Context, tokenId []byte) ([][]byte, error) {
	builder := dg.createDefaultVmQueryBuilder()
	builder.Function(getErc20AddressForTokenIdFuncName)
	builder.ArgBytes(tokenId)
	return dg.executeQueryFromBuilder(ctx, builder)
}

// WasProposedTransfer returns true if the transfer action proposed was triggered
func (dg *multiversXClientDataGetter) WasProposedTransfer(ctx context.Context, batch *clients.TransferBatch) (bool, error) {
	if batch == nil {
		return false, clients.ErrNilBatch
	}

	builder := dg.createDefaultVmQueryBuilder()
	builder.Function(wasTransferActionProposedFuncName).ArgInt64(int64(batch.ID))
	addBatchInfo(builder, batch)

	return dg.executeQueryBoolFromBuilder(ctx, builder)
}

// WasExecuted returns true if the provided actionID was executed or not
func (dg *multiversXClientDataGetter) WasExecuted(ctx context.Context, actionID uint64) (bool, error) {
	builder := dg.createDefaultVmQueryBuilder()
	builder.Function(wasActionExecutedFuncName).ArgInt64(int64(actionID))

	return dg.executeQueryBoolFromBuilder(ctx, builder)
}

// GetActionIDForProposeTransfer returns the action ID for the proposed transfer operation
func (dg *multiversXClientDataGetter) GetActionIDForProposeTransfer(ctx context.Context, batch *clients.TransferBatch) (uint64, error) {
	if batch == nil {
		return 0, clients.ErrNilBatch
	}

	builder := dg.createDefaultVmQueryBuilder()
	builder.Function(getActionIdForTransferBatchFuncName).ArgInt64(int64(batch.ID))
	addBatchInfo(builder, batch)

	return dg.executeQueryUint64FromBuilder(ctx, builder)
}

// WasProposedSetStatus returns true if the proposed set status was triggered
func (dg *multiversXClientDataGetter) WasProposedSetStatus(ctx context.Context, batch *clients.TransferBatch) (bool, error) {
	if batch == nil {
		return false, clients.ErrNilBatch
	}

	builder := dg.createDefaultVmQueryBuilder()
	builder.Function(wasSetCurrentTransactionBatchStatusActionProposedFuncName).ArgInt64(int64(batch.ID))
	for _, stat := range batch.Statuses {
		builder.ArgBytes([]byte{stat})
	}

	return dg.executeQueryBoolFromBuilder(ctx, builder)
}

// GetTransactionsStatuses will return the transactions statuses from the batch ID
func (dg *multiversXClientDataGetter) GetTransactionsStatuses(ctx context.Context, batchID uint64) ([]byte, error) {
	builder := dg.createDefaultVmQueryBuilder()
	builder.Function(getStatusesAfterExecutionFuncName).ArgInt64(int64(batchID))

	values, err := dg.executeQueryFromBuilder(ctx, builder)
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("%w for batch ID %v", errNoStatusForBatchID, batchID)
	}

	isFinished, err := dg.parseBool(values[0], getStatusesAfterExecutionFuncName, dg.multisigContractAddress.AddressAsBech32String())
	if err != nil {
		return nil, err
	}
	if !isFinished {
		return nil, fmt.Errorf("%w for batch ID %v", errBatchNotFinished, batchID)
	}

	results := make([]byte, len(values)-1)
	for i := 1; i < len(values); i++ {
		results[i-1], err = getStatusFromBuff(values[i])
		if err != nil {
			return nil, fmt.Errorf("%w for result index %d", err, i-1)
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("%w status is finished, no results are given", errMalformedBatchResponse)
	}

	return results, nil
}

// GetActionIDForSetStatusOnPendingTransfer returns the action ID for setting the status on the pending transfer batch
func (dg *multiversXClientDataGetter) GetActionIDForSetStatusOnPendingTransfer(ctx context.Context, batch *clients.TransferBatch) (uint64, error) {
	if batch == nil {
		return 0, clients.ErrNilBatch
	}

	builder := dg.createDefaultVmQueryBuilder()
	builder.Function(getActionIdForSetCurrentTransactionBatchStatusFuncName).ArgInt64(int64(batch.ID))
	for _, stat := range batch.Statuses {
		builder.ArgBytes([]byte{stat})
	}

	return dg.executeQueryUint64FromBuilder(ctx, builder)
}

// QuorumReached returns true if the provided action ID reached the set quorum
func (dg *multiversXClientDataGetter) QuorumReached(ctx context.Context, actionID uint64) (bool, error) {
	builder := dg.createDefaultVmQueryBuilder()
	builder.Function(quorumReachedFuncName).ArgInt64(int64(actionID))

	return dg.executeQueryBoolFromBuilder(ctx, builder)
}

// GetLastExecutedEthBatchID returns the last executed Ethereum batch ID
func (dg *multiversXClientDataGetter) GetLastExecutedEthBatchID(ctx context.Context) (uint64, error) {
	builder := dg.createDefaultVmQueryBuilder().Function(getLastExecutedEthBatchIdFuncName)

	return dg.executeQueryUint64FromBuilder(ctx, builder)
}

// GetLastExecutedEthTxID returns the last executed Ethereum deposit ID
func (dg *multiversXClientDataGetter) GetLastExecutedEthTxID(ctx context.Context) (uint64, error) {
	builder := dg.createDefaultVmQueryBuilder().Function(getLastExecutedEthTxId)

	return dg.executeQueryUint64FromBuilder(ctx, builder)
}

// WasSigned returns true if the action was already signed by the current relayer
func (dg *multiversXClientDataGetter) WasSigned(ctx context.Context, actionID uint64) (bool, error) {
	builder := dg.createDefaultVmQueryBuilder()
	builder.Function(signedFuncName).ArgAddress(dg.relayerAddress).ArgInt64(int64(actionID))

	return dg.executeQueryBoolFromBuilder(ctx, builder)
}

// GetAllStakedRelayers returns all staked relayers defined in MultiversX SC
func (dg *multiversXClientDataGetter) GetAllStakedRelayers(ctx context.Context) ([][]byte, error) {
	builder := dg.createDefaultVmQueryBuilder()
	builder.Function(getAllStakedRelayersFuncName)

	return dg.executeQueryFromBuilder(ctx, builder)
}

// IsPaused returns true if the multisig contract is paused
func (dg *multiversXClientDataGetter) IsPaused(ctx context.Context) (bool, error) {
	builder := dg.createDefaultVmQueryBuilder()
	builder.Function(isPausedFuncName)

	return dg.executeQueryBoolFromBuilder(ctx, builder)
}

func getStatusFromBuff(buff []byte) (byte, error) {
	if len(buff) == 0 {
		return 0, errMalformedBatchResponse
	}

	return buff[len(buff)-1], nil
}

func addBatchInfo(builder builders.VMQueryBuilder, batch *clients.TransferBatch) {
	for _, dt := range batch.Deposits {
		builder.ArgBytes(dt.FromBytes).
			ArgBytes(dt.ToBytes).
			ArgBytes(dt.ConvertedTokenBytes).
			ArgBigInt(dt.Amount).
			ArgInt64(int64(dt.Nonce))
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (dg *multiversXClientDataGetter) IsInterfaceNil() bool {
	return dg == nil
}
