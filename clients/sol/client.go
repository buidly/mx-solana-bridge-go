package sol

import "C"
import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
	chainCore "github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-solana-bridge-go/bridges/solanaMultiversX"
	"github.com/multiversx/mx-solana-bridge-go/clients"
	"github.com/multiversx/mx-solana-bridge-go/clients/sol/contract/bridge"
	"github.com/multiversx/mx-solana-bridge-go/clients/sol/contract/tokens_safe"
	"github.com/multiversx/mx-solana-bridge-go/core"
	"math"
	"math/big"
	"strconv"
	"sync"
	"time"
)

const (
	messagePrefix   = "\u0019Solana Signed Message:\n32"
	minQuorumValue  = uint64(1)
	minAllowedDelta = 1
)

type argListsBatch struct {
	tokens     []solana.PublicKey
	recipients []solana.PublicKey
	amounts    []*big.Int
	nonces     []*big.Int
}

type ArgsSolanaClient struct {
	SolanaRpcClient         *rpc.Client
	SolanaWssClientAddress  string
	DecimalDiffCalculator   clients.DecimalDiffCalculator
	ProgramAddresses        ProgramAddresses
	SftTokenProgram         SftTokenProgram
	StatusHandler           core.StatusHandler
	RelayerPrivateKey       solana.PrivateKey
	Log                     chainCore.Logger
	AddressConverter        core.AddressConverter
	Broadcaster             Broadcaster
	PrivateKey              *ecdsa.PrivateKey
	TokensMapper            TokensMapper
	SignatureHolder         SignaturesHolder
	GasHandler              GasHandler
	TransferGasLimitBase    uint64
	TransferGasLimitForEach uint64
	AllowDelta              uint64
}

type client struct {
	solanaRpcClient          *rpc.Client
	solanaWssClientAddress   string
	decimalDiffCalculator    clients.DecimalDiffCalculator
	programAddresses         ProgramAddresses
	sftTokenProgram          SftTokenProgram
	statusHandler            core.StatusHandler
	relayerPrivateKey        solana.PrivateKey
	log                      chainCore.Logger
	addressConverter         core.AddressConverter
	broadcaster              Broadcaster
	privateKey               *ecdsa.PrivateKey
	publicKey                *ecdsa.PublicKey
	tokensMapper             TokensMapper
	signatureHolder          SignaturesHolder
	gasHandler               GasHandler
	transferGasLimitBase     uint64
	transferGasLimitForEach  uint64
	allowDelta               uint64
	lastBlockNumber          uint64
	retriesAvailabilityCheck uint64
	mut                      sync.RWMutex
}

func NewSolanaClient(args ArgsSolanaClient) (*client, error) {
	err := checkArgs(args)
	if err != nil {
		return nil, err
	}

	publicKey := args.PrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errPublicKeyCast
	}

	c := &client{
		solanaRpcClient:         args.SolanaRpcClient,
		solanaWssClientAddress:  args.SolanaWssClientAddress,
		decimalDiffCalculator:   args.DecimalDiffCalculator,
		programAddresses:        args.ProgramAddresses,
		statusHandler:           args.StatusHandler,
		sftTokenProgram:         args.SftTokenProgram,
		relayerPrivateKey:       args.RelayerPrivateKey,
		log:                     args.Log,
		addressConverter:        args.AddressConverter,
		broadcaster:             args.Broadcaster,
		privateKey:              args.PrivateKey,
		publicKey:               publicKeyECDSA,
		tokensMapper:            args.TokensMapper,
		signatureHolder:         args.SignatureHolder,
		gasHandler:              args.GasHandler,
		transferGasLimitBase:    args.TransferGasLimitBase,
		transferGasLimitForEach: args.TransferGasLimitForEach,
		allowDelta:              args.AllowDelta,
	}

	c.log.Info("NewSolanaClient",
		"relayer solana address", c.relayerPrivateKey.PublicKey(),
		"relayer ecdsa public key address", crypto.PubkeyToAddress(*publicKeyECDSA),
		"safe contract address", c.programAddresses.TokensSafeProgramAddress,
		"bridge contract address", c.programAddresses.BridgeProgramAddress)

	return c, err
}

func checkArgs(args ArgsSolanaClient) error {
	if check.IfNil(args.StatusHandler) {
		return clients.ErrNilStatusHandler
	}
	if check.IfNil(args.Log) {
		return clients.ErrNilLogger
	}
	if check.IfNil(args.AddressConverter) {
		return clients.ErrNilAddressConverter
	}
	if check.IfNil(args.Broadcaster) {
		return errNilBroadcaster
	}
	if args.PrivateKey == nil {
		return clients.ErrNilPrivateKey
	}
	if check.IfNil(args.TokensMapper) {
		return clients.ErrNilTokensMapper
	}
	if check.IfNil(args.SignatureHolder) {
		return errNilSignaturesHolder
	}
	if check.IfNil(args.GasHandler) {
		return errNilGasHandler
	}
	if args.TransferGasLimitBase == 0 {
		return errInvalidGasLimit
	}
	if args.TransferGasLimitForEach == 0 {
		return errInvalidGasLimit
	}
	if args.AllowDelta < minAllowedDelta {
		return fmt.Errorf("%w for args.AllowedDelta, got: %d, minimum: %d",
			clients.ErrInvalidValue, args.AllowDelta, minAllowedDelta)
	}
	return nil
}

func (c *client) GetBatch(ctx context.Context, nonce uint64) (*clients.TransferBatch, error) {
	c.log.Info("Getting batch", "nonce", nonce)
	batch, err := c.getBatchData(nonce)
	if err != nil {
		return nil, err
	}
	if !batch.IsInitialized {
		return nil, errBatchNotInitialized
	}
	c.statusHandler.AddIntMetric(core.MetricNumSolClientRequests, 1)
	deposits := make([]tokens_safe.Deposit, batch.DepositsCount)
	for i := uint16(0); i < batch.DepositsCount; i++ {
		deposit, err := c.getDepositData(
			batch.FirstDepositNonce.BigInt().Add(batch.FirstDepositNonce.BigInt(), big.NewInt(0).SetUint64(uint64(i))),
		)
		if err != nil {
			return nil, err
		}
		deposits[i] = *deposit
	}
	if err != nil {
		return nil, err
	}
	if int(batch.DepositsCount) != len(deposits) {
		return nil, fmt.Errorf("%w, batch.DepositsCount: %d, fetched deposits len: %d",
			errDepositsAndBatchDepositsCountDiffer, batch.DepositsCount, len(deposits))
	}

	transferBatch := &clients.TransferBatch{
		ID:                    batch.Nonce.BigInt().Uint64(),
		Deposits:              make([]*clients.DepositTransfer, 0, batch.DepositsCount),
		LastUpdatedSlotNumber: batch.LastUpdatedSlotNumber,
		SlotNumber:            batch.SlotNumber,
	}
	cachedTokens := make(map[string][]byte)
	cachedTokensDecimalDifference := make(map[string]int)
	for i := range deposits {
		deposit := deposits[i]
		toBytes := deposit.Recipient[:]
		fromBytes := deposit.Depositor[:]
		tokenBytes := deposit.TokenMint[:]

		depositTransfer := &clients.DepositTransfer{
			Nonce:            deposit.Nonce.BigInt().Uint64(),
			ToBytes:          toBytes,
			DisplayableTo:    c.addressConverter.ToBech32String(toBytes),
			FromBytes:        fromBytes,
			DisplayableFrom:  c.addressConverter.ToHexString(fromBytes),
			TokenBytes:       tokenBytes,
			DisplayableToken: c.addressConverter.ToHexString(tokenBytes),
			Amount:           new(big.Int).SetUint64(deposit.Amount),
		}
		//here
		storedConvertedTokenBytes, exists := cachedTokens[depositTransfer.DisplayableToken]
		if !exists {
			depositTransfer.ConvertedTokenBytes, err = c.tokensMapper.ConvertToken(ctx, depositTransfer.TokenBytes)
			if err != nil {
				return nil, err
			}
			cachedTokens[depositTransfer.DisplayableToken] = depositTransfer.ConvertedTokenBytes
		} else {
			depositTransfer.ConvertedTokenBytes = storedConvertedTokenBytes
		}

		decimalDifference, exists := cachedTokensDecimalDifference[depositTransfer.DisplayableToken]
		if !exists {
			decimalDifference, err = c.decimalDiffCalculator.GetDecimalDifference(ctx, depositTransfer.TokenBytes, depositTransfer.ConvertedTokenBytes)
			if err != nil {
				return nil, err
			}
			cachedTokensDecimalDifference[depositTransfer.DisplayableToken] = decimalDifference
		}
		amountWithDecimals := big.NewFloat(0).Mul(
			big.NewFloat(0).SetInt(depositTransfer.Amount),
			big.NewFloat(0).SetFloat64(math.Pow10(decimalDifference*-1)),
		)
		depositTransfer.AmountAdjustedToDecimals, _ = amountWithDecimals.Int(nil)

		transferBatch.Deposits = append(transferBatch.Deposits, depositTransfer)
	}

	transferBatch.Statuses = make([]byte, len(transferBatch.Deposits))

	return transferBatch, nil
}

// WasExecuted returns true if the batch ID was executed
func (c *client) WasExecuted(ctx context.Context, batchID uint64) (bool, uint64, error) {
	c.statusHandler.AddIntMetric(core.MetricNumSolClientRequests, 1)
	batchExecutedData, err := c.getExecutedBatchData(batchID)
	if err != nil && err.Error() == "not found" {
		return false, 0, nil
	} else if err != nil {
		return false, 0, err
	}
	return batchExecutedData.IsExecuted, batchExecutedData.LastUpdatedSlotNumber, nil
}

// BroadcastSignatureForMessageHashes will send the signature for the provided message hash
func (c *client) BroadcastSignatureForMessageHashes(hashes map[uint64]common.Hash) {
	ethSigs := make([]*core.EthereumSignature, 0)
	for _, hash := range hashes {
		signature, err := crypto.Sign(hash.Bytes(), c.privateKey)
		if err != nil {
			c.log.Error("error generating signature", "msh hash", hash, "error", err)
			return
		}
		ethSigs = append(ethSigs, &core.EthereumSignature{
			MessageHash: hash.Bytes(),
			Signature:   signature,
		})
	}

	c.broadcaster.BroadcastSignature(ethSigs)
}

// GenerateMessageHash will generate the message hash based on the provided batch
func (c *client) GenerateMessageHash(
	recipient solana.PublicKey,
	mintToken solana.PublicKey,
	amount *big.Int,
	depositNonce *big.Int,
	batchNonce uint64,
) (common.Hash, error) {

	bytesArray := append(recipient.Bytes(), mintToken.Bytes()...)
	bytesArray = append(bytesArray, bigIntToLittleEndianBytes(amount, 8)...)
	bytesArray = append(bytesArray, bigIntToLittleEndianBytes(depositNonce, 16)...)
	bytesArray = append(bytesArray, convertUint64ToBuffer(batchNonce)...)
	bytesArray = append(bytesArray, []byte("executeTransferAction")...)

	hash := crypto.Keccak256Hash(bytesArray)
	return crypto.Keccak256Hash(append([]byte(messagePrefix), hash.Bytes()...)), nil
}

func (c *client) extractList(batch *clients.TransferBatch) (argListsBatch, error) {
	arg := argListsBatch{}

	for _, dt := range batch.Deposits {
		recipient := solana.PublicKeyFromBytes(dt.ToBytes)
		arg.recipients = append(arg.recipients, recipient)

		token := solana.PublicKeyFromBytes(dt.ConvertedTokenBytes)
		arg.tokens = append(arg.tokens, token)

		amount := big.NewInt(0).Set(dt.AmountAdjustedToDecimals)
		arg.amounts = append(arg.amounts, amount)

		nonce := dt.Nonce
		arg.nonces = append(arg.nonces, big.NewInt(0).SetUint64(nonce))
	}

	return arg, nil
}

// ExecuteTransfer will initiate and send the transaction from the transfer batch struct
func (c *client) ExecuteTransfer(
	ctx context.Context,
	depositMsgHashMap map[uint64]common.Hash,
	batch *clients.TransferBatch,
	quorum int,
) ([]string, error) {
	if batch == nil {
		return nil, clients.ErrNilBatch
	}

	c.statusHandler.AddIntMetric(core.MetricNumSolClientRequests, 1)

	bridgeSettings, err := c.getBridgeSettingsData()
	if err != nil {
		return nil, err
	}
	isPaused := bridgeSettings.PausableSettings.IsPaused
	if err != nil {
		return nil, fmt.Errorf("%w in client.ExecuteTransfer", err)
	}
	if isPaused {
		return nil, fmt.Errorf("%w in client.ExecuteTransfer", clients.ErrMultisigContractPaused)
	}

	depositNonceStart, err := c.getFirstDepositNonce(batch)
	if err != nil {
		return nil, err
	}

	c.log.Info("executing transfer " + batch.String())

	if err != nil {
		return nil, err
	}

	c.statusHandler.AddIntMetric(core.MetricNumSolClientRequests, 1)

	if err != nil {
		return nil, err
	}

	depositSignatureHashMap := make(map[uint64][][]byte)
	for depositNonce, hash := range depositMsgHashMap {

		signatures := c.signatureHolder.Signatures(hash.Bytes())
		if len(signatures) < quorum {
			return nil, fmt.Errorf("%w num signatures: %d, quorum: %d, deposit: %d", errQuorumNotReached, len(signatures), quorum, depositNonce)
		}
		if len(signatures) > quorum {
			c.log.Debug("reducing the size of the signatures set",
				"quorum", quorum, "total signatures", len(signatures), "deposit", depositNonce)
			signatures = signatures[:quorum]
		}
		depositSignatureHashMap[depositNonce] = signatures
	}

	relayerSettingsPDA, err := c.getRelayerSettingsPDA()
	if err != nil {
		return nil, err
	}
	bridgeSettingsPDA, err := c.getBridgeSettingsPDA()
	if err != nil {
		return nil, err
	}
	safeSettingsPDA, err := c.getSafeSettingsPDA()
	if err != nil {
		return nil, err
	}
	executedBatchPDA, err := c.getExecutedBatchPDA(batch.ID)
	if err != nil {
		return nil, err
	}

	var transferHashes []string

	c.log.Info("Connecting to WS", "batchID", batch.ID)
	wsClient, err := ws.Connect(ctx, c.solanaWssClientAddress)
	if err != nil {
		c.log.Info("WS failed to connect", "batchID", batch.ID)
		return nil, err
	}
	c.log.Info("WS connected successfully", "batchID", batch.ID)
	defer wsClient.Close()

	for _, dt := range batch.Deposits {

		if dt.Nonce < depositNonceStart {
			continue
		}

		c.log.Info("Executing transfer for deposit", "batchID", batch.ID, "depositID", dt.Nonce)

		recipient := solana.PublicKeyFromBytes(dt.ToBytes)
		mintTokenAddress := solana.PublicKeyFromBytes(dt.ConvertedTokenBytes)
		amount := dt.AmountAdjustedToDecimals
		err = c.checkCumulatedTransfers(ctx, mintTokenAddress, amount)
		if err != nil {
			return nil, err
		}

		batchNonceUint128, err := ag_binary.NewBinDecoder(convertUint64ToBuffer(batch.ID)).ReadUint128(ag_binary.LE)
		if err != nil {
			return nil, err
		}
		depositNonceUint128, err := ag_binary.NewBinDecoder(convertUint64ToBuffer(dt.Nonce)).ReadUint128(ag_binary.LE)
		if err != nil {
			return nil, err
		}
		transferSignatures := make([]bridge.TransferSignature, len(depositSignatureHashMap[dt.Nonce]))
		for i, item := range depositSignatureHashMap[dt.Nonce] {
			recoveryId := item[64]
			var bytes [64]uint8
			copy(bytes[:], item[:64])
			transferSignatures[i] = bridge.TransferSignature{
				RecoveryId: recoveryId,
				Signature: bridge.Secp256k1Signature{
					Bytes: bytes,
				},
			}
		}
		transferRequest := bridge.TransferRequest{
			Amount:       amount.Uint64(),
			DepositNonce: depositNonceUint128,
			Signatures:   transferSignatures,
		}

		recipientAta, err := c.getAssociatedTokenAccountForMintToken(recipient, mintTokenAddress)
		if err != nil {
			return nil, err
		}
		safeSettingsAta, err := c.getAssociatedTokenAccountForMintToken(safeSettingsPDA, mintTokenAddress)
		if err != nil {
			return nil, err
		}
		whitelistedTokenPDA, err := c.getWhitelistedTokenPDA(mintTokenAddress)
		if err != nil {
			return nil, err
		}
		executedTransferPDA, err := c.getExecutedTransferPDA(dt.Nonce)
		if err != nil {
			return nil, err
		}

		c.statusHandler.AddIntMetric(core.MetricNumSolClientRequests, 1)
		executeTransferTransaction := bridge.NewExecuteTransferInstructionBuilder().
			SetBatchNonce(batchNonceUint128).
			SetTransferRequest(transferRequest).
			SetIsLastTransfer(dt.Nonce == batch.Deposits[len(batch.Deposits)-1].Nonce).
			SetRelayerAccount(c.relayerPrivateKey.PublicKey()).
			SetBridgeSettingsAccount(bridgeSettingsPDA).
			SetRelayerSettingsAccount(relayerSettingsPDA).
			SetSafeSettingsAccount(safeSettingsPDA).
			SetExecutedTransferAccount(executedTransferPDA).
			SetWhitelistedTokenAccount(whitelistedTokenPDA).
			SetExecutedBatchAccount(executedBatchPDA).
			SetTokenMintAccount(mintTokenAddress).
			SetSafeSettingsAtaAccount(*safeSettingsAta).
			SetSendToAccount(*recipientAta).
			SetRecipientAccount(recipient).
			SetTokensSafeProgramAccount(c.programAddresses.TokensSafeProgramAddress).
			SetSystemProgramAccount(solana.SystemProgramID).
			SetTokenProgramAccount(solana.TokenProgramID).
			SetAssociatedTokenProgramAccount(solana.SPLAssociatedTokenAccountProgramID).
			Build()

		recent, err := c.solanaRpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
		if err != nil {
			return nil, err
		}

		additionalComputeUnitsInstructionData := append(
			bigIntToLittleEndianBytes(big.NewInt(0).SetUint64(2), 1),
			bigIntToLittleEndianBytes(big.NewInt(0).SetUint64(300000), 4)...,
		)
		additionalComputeUnitsInstruction := solana.NewInstruction(
			solana.MustPublicKeyFromBase58("ComputeBudget111111111111111111111111111111"),
			solana.AccountMetaSlice{},
			additionalComputeUnitsInstructionData,
		)

		tx, err := solana.NewTransaction(
			[]solana.Instruction{executeTransferTransaction, additionalComputeUnitsInstruction},
			recent.Value.Blockhash,
			solana.TransactionPayer(c.relayerPrivateKey.PublicKey()),
		)
		if err != nil {
			return nil, err
		}

		fee, err := c.solanaRpcClient.GetFeeForMessage(context.TODO(), tx.Message.ToBase64(), rpc.CommitmentFinalized)
		if err != nil {
			return nil, err
		}
		if fee == nil || fee.Value == nil {
			return nil, errNilFee
		}
		err = c.checkRelayerFundsForFee(ctx, *fee.Value)
		if err != nil {
			return nil, err
		}

		_, err = tx.Sign(
			func(key solana.PublicKey) *solana.PrivateKey {
				if c.relayerPrivateKey.PublicKey().Equals(key) {
					return &c.relayerPrivateKey
				}
				return nil
			},
		)
		if err != nil {
			return nil, err
		}
		c.log.Info("Sending transaction", "batchID", batch.ID, "depositID", dt.Nonce)
		sig, err := confirm.SendAndConfirmTransactionWithTimeout(
			ctx,
			c.solanaRpcClient,
			wsClient,
			tx,
			2*time.Minute,
		)

		if err != nil {
			c.log.Error("Transaction failed", "batchID", batch.ID, "depositID", dt.Nonce, "hash", sig.String(), err)
			return nil, err
		}

		c.log.Info("Executed transfer transaction", "batchID", batch.ID, "depositID", dt.Nonce, "hash", sig.String())
		transferHashes = append(transferHashes, sig.String())
	}

	return transferHashes, err
}

func (c *client) getFirstDepositNonce(batch *clients.TransferBatch) (uint64, error) {
	executedBatchData, err := c.getExecutedBatchData(batch.ID)
	depositNonceStart := batch.Deposits[0].Nonce
	if err != nil && err.Error() != "not found" {
		return 0, fmt.Errorf("%w in client.ExecuteTransfer", err)
	} else if err == nil {
		depositNonceStart = executedBatchData.LastDepositNonce.BigInt().Uint64() + 1
	}
	return depositNonceStart, nil
}

func (c *client) GetLastFinalizedSlotNumber(ctx context.Context) (uint64, error) {
	currentBlock, err := c.solanaRpcClient.GetSlot(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return 0, err
	}
	return currentBlock, nil
}

// CheckClientAvailability will check the client availability and set the metric accordingly
func (c *client) CheckClientAvailability(ctx context.Context) error {
	c.mut.Lock()
	defer c.mut.Unlock()

	c.statusHandler.AddIntMetric(core.MetricNumSolClientRequests, 1)
	currentBlock, err := c.solanaRpcClient.GetBlockHeight(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		c.setStatusForAvailabilityCheck(solmultiversx.Unavailable, err.Error(), currentBlock)

		return err
	}

	if currentBlock != c.lastBlockNumber {
		c.retriesAvailabilityCheck = 0
		c.lastBlockNumber = currentBlock
	}

	// if we reached this point we will need to increment the retries counter
	defer c.incrementRetriesAvailabilityCheck()

	if c.retriesAvailabilityCheck > c.allowDelta {
		message := fmt.Sprintf("block %d fetched for %d times in a row", currentBlock, c.retriesAvailabilityCheck)
		c.setStatusForAvailabilityCheck(solmultiversx.Unavailable, message, currentBlock)

		return nil
	}

	c.setStatusForAvailabilityCheck(solmultiversx.Available, "", currentBlock)

	return nil
}

func (c *client) incrementRetriesAvailabilityCheck() {
	c.retriesAvailabilityCheck++
}

func (c *client) setStatusForAvailabilityCheck(status solmultiversx.ClientStatus, message string, nonce uint64) {
	c.statusHandler.SetStringMetric(core.MetricMultiversXClientStatus, status.String())
	c.statusHandler.SetStringMetric(core.MetricLastMultiversXClientError, message)
	c.statusHandler.SetIntMetric(core.MetricLastBlockNonce, int(nonce))
}

func (c *client) checkCumulatedTransfers(ctx context.Context, sftTokenProgramAddress solana.PublicKey, amount *big.Int) error {
	safeSettingsPDA, err := c.getSafeSettingsPDA()
	if err != nil {
		return err
	}
	tokensSafeAta, err := c.getAssociatedTokenAccountForMintToken(safeSettingsPDA, sftTokenProgramAddress)
	spew.Dump(c.programAddresses.TokensSafeProgramAddress, sftTokenProgramAddress, tokensSafeAta)
	if err != nil {
		return err
	}
	existingBalance, err := c.sftTokenProgram.BalanceOf(ctx, sftTokenProgramAddress, *tokensSafeAta)
	if err != nil {
		return fmt.Errorf("%w for address %s for ERC20 token %s", err, c.programAddresses.TokensSafeProgramAddress, sftTokenProgramAddress.String())
	}

	existingBalanceAmount, err := strconv.ParseUint(existingBalance.Amount, 10, 64)
	valueToBeTransferred := amount.Uint64()
	if valueToBeTransferred > existingBalanceAmount {
		return fmt.Errorf("%w, existing: %s, required: %s for ERC20 token %s and address %s",
			errInsufficientSftBalance, existingBalance.UiAmountString, amount.String(), sftTokenProgramAddress.String(), c.programAddresses.TokensSafeProgramAddress)
	}

	c.log.Debug("checked ERC20 balance",
		"ERC20 token", sftTokenProgramAddress.String(),
		"address", c.programAddresses.TokensSafeProgramAddress,
		"existing balance", existingBalanceAmount,
		"needed", valueToBeTransferred)

	return nil
}

func (c *client) getBatchData(nonce uint64) (*tokens_safe.Batch, error) {
	batchPDA, _, err := solana.FindProgramAddress(
		[][]byte{[]byte(BATCH_SEED_PREFIX), convertUint64ToBuffer(nonce)},
		c.programAddresses.TokensSafeProgramAddress,
	)
	if err != nil {
		return nil, err
	}
	var batch tokens_safe.Batch
	err = c.solanaRpcClient.GetAccountDataInto(
		context.TODO(),
		batchPDA,
		&batch,
	)
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

func (c *client) GetSafeSettings() (*tokens_safe.SafeSettings, error) {
	safeSettingsPDA, err := c.getSafeSettingsPDA()
	if err != nil {
		return nil, err
	}
	var safeSettings tokens_safe.SafeSettings
	err = c.solanaRpcClient.GetAccountDataInto(
		context.TODO(),
		safeSettingsPDA,
		&safeSettings,
	)
	if err != nil {
		return nil, err
	}
	return &safeSettings, nil
}

func (c *client) getDepositData(depositNonce *big.Int) (*tokens_safe.Deposit, error) {
	depositPDA, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(DEPOSIT_SEED_PREFIX),
			bigIntToLittleEndianBytes(depositNonce, 16),
		},
		c.programAddresses.TokensSafeProgramAddress,
	)
	if err != nil {
		return nil, err
	}
	var deposit tokens_safe.Deposit
	err = c.solanaRpcClient.GetAccountDataInto(
		context.TODO(),
		depositPDA,
		&deposit,
	)
	if err != nil {
		return nil, err
	}
	return &deposit, nil
}

func (c *client) getBridgeSettingsData() (*bridge.BridgeSettings, error) {
	bridgeSettingsPDA, err := c.getBridgeSettingsPDA()
	if err != nil {
		return nil, err
	}
	var bridgeSettings bridge.BridgeSettings
	err = c.solanaRpcClient.GetAccountDataInto(
		context.TODO(),
		bridgeSettingsPDA,
		&bridgeSettings,
	)
	return &bridgeSettings, nil
}

func (c *client) getExecutedTransferData(depositNonce uint64) (*bridge.ExecutedTransfer, error) {
	executedTransferPDA, err := c.getExecutedTransferPDA(depositNonce)
	if err != nil {
		return nil, err
	}
	var executedTransfer bridge.ExecutedTransfer
	err = c.solanaRpcClient.GetAccountDataInto(
		context.TODO(),
		executedTransferPDA,
		&executedTransfer,
	)
	if err != nil {
		return nil, err
	}
	return &executedTransfer, nil
}

func (c *client) getExecutedTransferPDA(depositNonce uint64) (solana.PublicKey, error) {
	executedTransferPDA, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(EXECUTED_TRANFER_PREFIX),
			convertUint64ToBuffer(depositNonce),
		},
		c.programAddresses.BridgeProgramAddress,
	)
	return executedTransferPDA, err
}

func (c *client) getBridgeSettingsPDA() (solana.PublicKey, error) {
	bridgeSettingsPDA, _, err := solana.FindProgramAddress(
		[][]byte{},
		c.programAddresses.BridgeProgramAddress,
	)
	return bridgeSettingsPDA, err
}

func (c *client) getRelayerSettingsPDA() (solana.PublicKey, error) {
	relayerSettingsPDA, _, err := solana.FindProgramAddress(
		[][]byte{[]byte(RELAYER_SETTINGS_SEED_PREFIX)},
		c.programAddresses.BridgeProgramAddress,
	)
	return relayerSettingsPDA, err
}

func (c *client) getExecutedBatchData(batchNonce uint64) (*bridge.ExecutedBatch, error) {
	executedBatchPDA, err := c.getExecutedBatchPDA(batchNonce)
	if err != nil {
		return nil, err
	}
	var executedBatch bridge.ExecutedBatch
	err = c.solanaRpcClient.GetAccountDataInto(
		context.TODO(),
		executedBatchPDA,
		&executedBatch,
	)
	if err != nil {
		return nil, err
	}
	return &executedBatch, nil
}

func (c *client) getExecutedBatchPDA(batchNonce uint64) (solana.PublicKey, error) {
	executedBatchPDA, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(EXECUTED_BATCH_SEED_PREFIX),
			convertUint64ToBuffer(batchNonce),
		},
		c.programAddresses.BridgeProgramAddress,
	)
	return executedBatchPDA, err
}

func (c *client) getSafeSettingsPDA() (solana.PublicKey, error) {
	safeSettingsPDA, _, err := solana.FindProgramAddress(
		[][]byte{},
		c.programAddresses.TokensSafeProgramAddress,
	)
	return safeSettingsPDA, err
}

func (c *client) getWhitelistedTokenPDA(tokenMint solana.PublicKey) (solana.PublicKey, error) {
	whitelistedTokenPDA, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(WHITELISTED_TOKEN_SEED_PREFIX),
			tokenMint.Bytes(),
		},
		c.programAddresses.TokensSafeProgramAddress,
	)
	return whitelistedTokenPDA, err
}

func (c *client) getAssociatedTokenAccountForMintToken(owner solana.PublicKey, sftTokenProgramAddress solana.PublicKey) (*solana.PublicKey, error) {
	ata, _, err := solana.FindProgramAddress(
		[][]byte{
			owner.Bytes(),
			solana.TokenProgramID.Bytes(),
			sftTokenProgramAddress.Bytes(),
		},
		solana.SPLAssociatedTokenAccountProgramID,
	)
	if err != nil {
		return nil, err
	}

	return &ata, nil
}

func (c *client) checkRelayerFundsForFee(ctx context.Context, transferFee uint64) error {

	c.statusHandler.AddIntMetric(core.MetricNumSolClientRequests, 1)
	existingBalance, err := c.solanaRpcClient.GetBalance(context.TODO(), c.relayerPrivateKey.PublicKey(), rpc.CommitmentFinalized)
	if err != nil {
		return err
	}

	if transferFee > existingBalance.Value {
		return fmt.Errorf("%w, existing: %s, required: %s",
			errInsufficientBalance, existingBalance.Value, transferFee)
	}

	c.log.Debug("checked balance",
		"existing balance", existingBalance.Value,
		"needed", transferFee)

	return nil
}

func (c *client) GetTransactionsStatuses(ctx context.Context, batch *clients.TransferBatch) ([]byte, error) {
	var depositStatuses []byte
	for _, dt := range batch.Deposits {
		executedTransfer, err := c.getExecutedTransferData(dt.Nonce)
		//TODO discuss if it stays like this
		if err != nil && err.Error() == "not found" {
			depositStatuses = append(depositStatuses, bridge.TransferStatusREJECTED.Byte())
			continue
		}
		if err != nil {
			return nil, err
		}
		depositStatuses = append(depositStatuses, executedTransfer.Status.Byte())
	}
	return depositStatuses, nil
}

func (c *client) GetQuorumSize(ctx context.Context) (*big.Int, error) {
	c.statusHandler.AddIntMetric(core.MetricNumSolClientRequests, 1)
	bridgeSettings, err := c.getBridgeSettingsData()
	if err != nil {
		return nil, err
	}
	return big.NewInt(0).SetUint64(uint64(bridgeSettings.QuorumSettings.Quorum)), nil
}

// IsQuorumReached returns true if the number of signatures is at least the size of quorum
func (c *client) IsQuorumReached(ctx context.Context, batch *clients.TransferBatch, depositMsgHashMap map[uint64]common.Hash) (bool, error) {
	depositNonceStart, err := c.getFirstDepositNonce(batch)
	quorum, err := c.GetQuorumSize(ctx)
	if err != nil {
		return false, fmt.Errorf("%w in IsQuorumReached, Quorum call", err)
	}
	isQuorumReachedForAllDeposits := true
	for _, dt := range batch.Deposits {
		if dt.Nonce < depositNonceStart {
			continue
		}
		signatures := c.signatureHolder.Signatures(depositMsgHashMap[dt.Nonce].Bytes())
		if quorum.Uint64() < minQuorumValue {
			return false, fmt.Errorf("%w in IsQuorumReached, minQuorum %d, got: %s", clients.ErrInvalidValue, minQuorumValue, quorum.String())
		}
		c.log.Info("number of signatures",
			"deposit nonce", dt.Nonce,
			"signatures number", len(signatures),
		)
		isQuorumReachedForAllDeposits = isQuorumReachedForAllDeposits && len(signatures) >= int(quorum.Int64())
	}

	return isQuorumReachedForAllDeposits, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (c *client) IsInterfaceNil() bool {
	return c == nil
}
