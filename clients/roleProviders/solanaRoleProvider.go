package roleproviders

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/multiversx/mx-solana-bridge-go/clients/sol"
	"github.com/multiversx/mx-solana-bridge-go/clients/sol/contract/bridge"
	"golang.org/x/crypto/sha3"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-solana-bridge-go/clients"
)

const ethSignatureSize = 64

// ArgsSolanaRoleProvider is the argument for the sol role provider constructor
type ArgsSolanaRoleProvider struct {
	SolanaRpcClient  *rpc.Client
	ProgramAddresses sol.ProgramAddresses
	Log              logger.Logger
}

type solanaRoleProvider struct {
	solanaRpcClient            *rpc.Client
	programAddresses           sol.ProgramAddresses
	log                        logger.Logger
	whitelistedSolanaAddresses map[solana.PublicKey]struct{}
	whitelistedEthAddresses    map[common.Address]struct{}
	mut                        sync.RWMutex
}

// NewSolanaRoleProvider creates a new sol role provider instance able to fetch the
// whitelisted addresses and able to check sol signatures
func NewSolanaRoleProvider(args ArgsSolanaRoleProvider) (*solanaRoleProvider, error) {
	err := checkSolanaRoleProviderSpecificArgs(args)
	if err != nil {
		return nil, err
	}

	erp := &solanaRoleProvider{
		whitelistedSolanaAddresses: make(map[solana.PublicKey]struct{}),
		whitelistedEthAddresses:    make(map[common.Address]struct{}),
		solanaRpcClient:            args.SolanaRpcClient,
		programAddresses:           args.ProgramAddresses,
		log:                        args.Log,
	}

	return erp, nil
}

func checkSolanaRoleProviderSpecificArgs(args ArgsSolanaRoleProvider) error {
	if check.IfNil(args.Log) {
		return clients.ErrNilLogger
	}

	return nil
}

// Execute will fetch the available whitelisted addresses and store them in the inner map
func (erp *solanaRoleProvider) Execute(ctx context.Context) error {
	relayersPDA, _, err := solana.FindProgramAddress([][]byte{[]byte(RELAYER_SETTINGS_SEED_PREFIX)}, erp.programAddresses.BridgeProgramAddress)
	if err != nil {
		return err
	}
	resp, err := erp.solanaRpcClient.GetAccountInfoWithOpts(
		context.TODO(),
		relayersPDA,
		&rpc.GetAccountInfoOpts{
			Encoding:   solana.EncodingBase64Zstd,
			Commitment: rpc.CommitmentFinalized,
		},
	)
	if err != nil {
		return err
	}
	var relayersSettings bridge.RelayerSettings
	//TODO check if can be fixed differently
	moveBytesToEnd(resp.Value.Data.GetBinary(), 10, 12)
	err = bin.NewBinDecoder(resp.Value.Data.GetBinary()).Decode(&relayersSettings)
	if err != nil {
		panic(err)
	}

	erp.processResults(relayersSettings.Relayers)

	return nil
}

func (erp *solanaRoleProvider) processResults(results []bridge.RelayerDetails) {
	currentList := make([]string, 0, len(results))
	ethCurrentList := make([]string, 0, len(results))

	erp.mut.Lock()
	erp.whitelistedSolanaAddresses = make(map[solana.PublicKey]struct{})
	erp.whitelistedEthAddresses = make(map[common.Address]struct{})

	for _, addr := range results {
		erp.whitelistedSolanaAddresses[addr.SolanaAddress] = struct{}{}
		erp.whitelistedEthAddresses[publicKeyToAddress(addr.SecpAddress.Bytes)] = struct{}{}
		currentList = append(currentList, addr.SolanaAddress.String())
		ethCurrentList = append(ethCurrentList, publicKeyToAddress(addr.SecpAddress.Bytes).String())
	}
	erp.mut.Unlock()

	erp.log.Debug("fetched Solana whitelisted addresses:\n" + strings.Join(currentList, "\n"))
	erp.log.Debug("fetched Solana Eth whitelisted addresses:\n" + strings.Join(ethCurrentList, "\n"))
}

// VerifyEthSignature will verify the provided signature against the message hash. It will also checks if the
// resulting public key is whitelisted or not
func (erp *solanaRoleProvider) VerifyEthSignature(signature []byte, messageHash []byte) error {
	pkBytes, err := crypto.Ecrecover(messageHash, signature)
	if err != nil {
		return err
	}

	pk, err := crypto.UnmarshalPubkey(pkBytes)
	if err != nil {
		return err
	}

	address := crypto.PubkeyToAddress(*pk)
	if !erp.isWhitelisted(address) {
		return ErrAddressIsNotWhitelisted
	}

	if len(signature) > ethSignatureSize {
		// signatures might contain the recovery byte
		signature = signature[:ethSignatureSize]
	}

	sigOk := crypto.VerifySignature(pkBytes, messageHash, signature)
	if !sigOk {
		return ErrInvalidSignature
	}

	return nil
}

func publicKeyToAddress(publicKeyBytes [64]uint8) common.Address {
	keccakHash := sha3.NewLegacyKeccak256()
	keccakHash.Write(publicKeyBytes[:])
	hashed := keccakHash.Sum(nil)

	var ethereumAddress common.Address
	copy(ethereumAddress[:], hashed[12:])

	return ethereumAddress
}

func (erp *solanaRoleProvider) isWhitelisted(address common.Address) bool {
	erp.mut.RLock()
	defer erp.mut.RUnlock()

	_, exists := erp.whitelistedEthAddresses[address]

	return exists
}

// IsInterfaceNil returns true if there is no value under the interface
func (erp *solanaRoleProvider) IsInterfaceNil() bool {
	return erp == nil
}

func moveBytesToEnd(byteArray []byte, start, end int) {
	if start < 0 || start >= len(byteArray) || end < 0 || end >= len(byteArray) || start > end {
		return
	}

	tempBuffer := make([]byte, end-start+1)
	copy(tempBuffer, byteArray[start:end+1])

	copy(byteArray[start:], byteArray[end+1:])

	copy(byteArray[len(byteArray)-len(tempBuffer):], tempBuffer)
}
