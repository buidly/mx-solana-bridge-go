package factory

import (
	"context"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	solmultiversx "github.com/multiversx/mx-solana-bridge-go/bridges/solanaMultiversX"
	multiversxtosol "github.com/multiversx/mx-solana-bridge-go/bridges/solanaMultiversX/steps/multiversxToSolana"
	soltomultiversx "github.com/multiversx/mx-solana-bridge-go/bridges/solanaMultiversX/steps/solanaToMultiversX"
	roleproviders "github.com/multiversx/mx-solana-bridge-go/clients/roleProviders"
	"io"
	"io/ioutil"
	"sync"
	"time"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	chainCore "github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	crypto "github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519/singlesig"
	chainConfig "github.com/multiversx/mx-chain-go/config"
	antifloodFactory "github.com/multiversx/mx-chain-go/process/throttle/antiflood/factory"
	logger "github.com/multiversx/mx-chain-logger-go"
	sdkCore "github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/core/polling"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/interactors"
	"github.com/multiversx/mx-solana-bridge-go/bridges/solanaMultiversX/disabled"
	"github.com/multiversx/mx-solana-bridge-go/bridges/solanaMultiversX/topology"
	"github.com/multiversx/mx-solana-bridge-go/clients"
	batchValidatorManagement "github.com/multiversx/mx-solana-bridge-go/clients/batchValidator"
	batchManagementFactory "github.com/multiversx/mx-solana-bridge-go/clients/batchValidator/factory"
	"github.com/multiversx/mx-solana-bridge-go/clients/chain"
	"github.com/multiversx/mx-solana-bridge-go/clients/gasManagement/factory"
	"github.com/multiversx/mx-solana-bridge-go/clients/multiversx"
	"github.com/multiversx/mx-solana-bridge-go/clients/multiversx/mappers"
	"github.com/multiversx/mx-solana-bridge-go/clients/sol"
	"github.com/multiversx/mx-solana-bridge-go/config"
	"github.com/multiversx/mx-solana-bridge-go/core"
	"github.com/multiversx/mx-solana-bridge-go/core/converters"
	"github.com/multiversx/mx-solana-bridge-go/core/timer"
	"github.com/multiversx/mx-solana-bridge-go/p2p"
	"github.com/multiversx/mx-solana-bridge-go/stateMachine"
	"github.com/multiversx/mx-solana-bridge-go/status"
)

const (
	minTimeForBootstrap     = time.Millisecond * 100
	minTimeBeforeRepeatJoin = time.Second * 30
	pollingDurationOnError  = time.Second * 5
)

var suite = ed25519.NewEd25519()
var keyGen = signing.NewKeyGenerator(suite)
var singleSigner = &singlesig.Ed25519Signer{}

type ArgsSolanaToMultiversXBridge struct {
	ProgramAddresses              sol.ProgramAddresses
	SolStatusHandler              core.StatusHandler
	RelayerPrivateKey             solana.PrivateKey
	Configs                       config.Configs
	Messenger                     p2p.NetMessenger
	StatusStorer                  core.Storer
	Proxy                         clients.Proxy
	MultiversXClientStatusHandler core.StatusHandler
	SftTokenProgram               sol.SftTokenProgram
	SolanaRpcClient               *rpc.Client
	SolanaWssClientAddress        string
	DecimalDiffCalculator         clients.DecimalDiffCalculator
	TimeForBootstrap              time.Duration
	TimeBeforeRepeatJoin          time.Duration
	MetricsHolder                 core.MetricsHolder
	AppStatusHandler              chainCore.AppStatusHandler
}

type solanaMultiversXBridgeComponents struct {
	programAddresses                  sol.ProgramAddresses
	solStatusHandler                  core.StatusHandler
	relayerPrivateKey                 solana.PrivateKey
	solanaRpcClient                   *rpc.Client
	solanaWssClientAddress            string
	decimalDiffCalculator             clients.DecimalDiffCalculator
	baseLogger                        logger.Logger
	messenger                         p2p.NetMessenger
	statusStorer                      core.Storer
	multiversXClient                  solmultiversx.MultiversXClient
	solanaClient                      solmultiversx.SolanaClient
	compatibleChain                   chain.Chain
	multiversXMultisigContractAddress sdkCore.AddressHandler
	multiversXRelayerPrivateKey       crypto.PrivateKey
	multiversXRelayerAddress          sdkCore.AddressHandler
	mxDataGetter                      dataGetter
	proxy                             clients.Proxy
	multiversXRoleProvider            MultiversXRoleProvider
	solanaRoleProvider                SolanaRoleProvider
	broadcaster                       Broadcaster
	timer                             core.Timer
	timeForBootstrap                  time.Duration
	metricsHolder                     core.MetricsHolder
	addressConverter                  core.AddressConverter
	SftTokenProgram                   sol.SftTokenProgram
	solToMultiversXMachineStates      core.MachineStates
	solToMultiversXStepDuration       time.Duration
	solToMultiversXStatusHandler      core.StatusHandler
	solToMultiversXStateMachine       StateMachine
	solToMultiversXSignaturesHolder   solmultiversx.SignaturesHolder

	multiversXToSolMachineStates core.MachineStates
	multiversXToSolStepDuration  time.Duration
	multiversXToSolStatusHandler core.StatusHandler
	multiversXToSolStateMachine  StateMachine

	mutClosableHandlers sync.RWMutex
	closableHandlers    []io.Closer

	pollingHandlers []PollingHandler

	timeBeforeRepeatJoin time.Duration
	cancelFunc           func()
	appStatusHandler     chainCore.AppStatusHandler
}

func NewSolanaMultiversXBridgeComponents(args ArgsSolanaToMultiversXBridge) (*solanaMultiversXBridgeComponents, error) {
	err := checkArgsSolanaToMultiversXBridge(args)
	if err != nil {
		return nil, err
	}
	compatibleChain := args.Configs.GeneralConfig.Solana.Chain
	solToMultiversXName := compatibleChain.SolanaToMultiversXName()
	baseLogId := compatibleChain.BaseLogId()
	components := &solanaMultiversXBridgeComponents{
		solStatusHandler:       args.SolStatusHandler,
		solanaRpcClient:        args.SolanaRpcClient,
		solanaWssClientAddress: args.SolanaWssClientAddress,
		decimalDiffCalculator:  args.DecimalDiffCalculator,
		relayerPrivateKey:      args.RelayerPrivateKey,
		baseLogger:             core.NewLoggerWithIdentifier(logger.GetOrCreate(solToMultiversXName), baseLogId),
		compatibleChain:        compatibleChain,
		messenger:              args.Messenger,
		statusStorer:           args.StatusStorer,
		closableHandlers:       make([]io.Closer, 0),
		proxy:                  args.Proxy,
		timer:                  timer.NewNTPTimer(),
		timeForBootstrap:       args.TimeForBootstrap,
		timeBeforeRepeatJoin:   args.TimeBeforeRepeatJoin,
		metricsHolder:          args.MetricsHolder,
		appStatusHandler:       args.AppStatusHandler,
		SftTokenProgram:        args.SftTokenProgram,
		programAddresses:       args.ProgramAddresses,
	}

	addressConverter, err := converters.NewAddressConverter()
	if err != nil {
		return nil, clients.ErrNilAddressConverter
	}
	components.addressConverter = addressConverter

	components.addClosableComponent(components.timer)

	err = components.createMultiversXKeysAndAddresses(args.Configs.GeneralConfig.MultiversX)
	if err != nil {
		return nil, err
	}

	err = components.createDataGetter()
	if err != nil {
		return nil, err
	}

	err = components.createMultiversXRoleProvider(args)
	if err != nil {
		return nil, err
	}

	err = components.createMultiversXClient(args)
	if err != nil {
		return nil, err
	}

	err = components.createSolanaRoleProvider(args)
	if err != nil {
		return nil, err
	}

	err = components.createSolanaClient(args)
	if err != nil {
		return nil, err
	}

	err = components.createSolanaToMultiversXBridge(args)
	if err != nil {
		return nil, err
	}

	err = components.createSolanaToMultiversXStateMachine()
	if err != nil {
		return nil, err
	}

	err = components.createMultiversXToSolanaBridge(args)
	if err != nil {
		return nil, err
	}

	err = components.createMultiversXToSolanaStateMachine()
	if err != nil {
		return nil, err
	}

	return components, nil
}

func (components *solanaMultiversXBridgeComponents) addClosableComponent(closable io.Closer) {
	components.mutClosableHandlers.Lock()
	components.closableHandlers = append(components.closableHandlers, closable)
	components.mutClosableHandlers.Unlock()
}

func checkArgsSolanaToMultiversXBridge(args ArgsSolanaToMultiversXBridge) error {
	if check.IfNil(args.Proxy) {
		return errNilProxy
	}
	if check.IfNil(args.Messenger) {
		return errNilMessenger
	}
	if check.IfNil(args.StatusStorer) {
		return errNilStatusStorer
	}
	if check.IfNil(args.SftTokenProgram) {
		return errNilErc20ContractsHolder
	}
	if args.TimeForBootstrap < minTimeForBootstrap {
		return fmt.Errorf("%w for TimeForBootstrap, received: %v, minimum: %v", errInvalidValue, args.TimeForBootstrap, minTimeForBootstrap)
	}
	if args.TimeBeforeRepeatJoin < minTimeBeforeRepeatJoin {
		return fmt.Errorf("%w for TimeBeforeRepeatJoin, received: %v, minimum: %v", errInvalidValue, args.TimeBeforeRepeatJoin, minTimeBeforeRepeatJoin)
	}
	if check.IfNil(args.MetricsHolder) {
		return errNilMetricsHolder
	}
	if check.IfNil(args.AppStatusHandler) {
		return errNilStatusHandler
	}

	return nil
}

func (components *solanaMultiversXBridgeComponents) createMultiversXKeysAndAddresses(chainConfigs config.MultiversXConfig) error {
	wallet := interactors.NewWallet()
	multiversXPrivateKeyBytes, err := wallet.LoadPrivateKeyFromPemFile(chainConfigs.PrivateKeyFile)
	if err != nil {
		return err
	}

	components.multiversXRelayerPrivateKey, err = keyGen.PrivateKeyFromByteArray(multiversXPrivateKeyBytes)
	if err != nil {
		return err
	}

	components.multiversXRelayerAddress, err = wallet.GetAddressFromPrivateKey(multiversXPrivateKeyBytes)
	if err != nil {
		return err
	}

	components.multiversXMultisigContractAddress, err = data.NewAddressFromBech32String(chainConfigs.MultisigContractAddress)
	if err != nil {
		return fmt.Errorf("%w for chainConfigs.BridgeProgramAddress", err)
	}

	return nil
}

func (components *solanaMultiversXBridgeComponents) createDataGetter() error {
	multiversXDataGetterLogId := components.compatibleChain.MultiversXDataGetterLogId()
	argsMXClientDataGetter := multiversx.ArgsMXClientDataGetter{
		MultisigContractAddress: components.multiversXMultisigContractAddress,
		RelayerAddress:          components.multiversXRelayerAddress,
		Proxy:                   components.proxy,
		Log:                     core.NewLoggerWithIdentifier(logger.GetOrCreate(multiversXDataGetterLogId), multiversXDataGetterLogId),
	}

	var err error
	components.mxDataGetter, err = multiversx.NewMXClientDataGetter(argsMXClientDataGetter)

	return err
}

func (components *solanaMultiversXBridgeComponents) createMultiversXClient(args ArgsSolanaToMultiversXBridge) error {
	chainConfigs := args.Configs.GeneralConfig.MultiversX
	tokensMapper, err := mappers.NewMultiversXToErc20Mapper(components.mxDataGetter)
	if err != nil {
		return err
	}
	multiversXClientLogId := components.compatibleChain.MultiversXClientLogId()

	clientArgs := multiversx.ClientArgs{
		DecimalDiffCalculator:        args.DecimalDiffCalculator,
		GasMapConfig:                 chainConfigs.GasMap,
		Proxy:                        args.Proxy,
		Log:                          core.NewLoggerWithIdentifier(logger.GetOrCreate(multiversXClientLogId), multiversXClientLogId),
		RelayerPrivateKey:            components.multiversXRelayerPrivateKey,
		MultisigContractAddress:      components.multiversXMultisigContractAddress,
		IntervalToResendTxsInSeconds: chainConfigs.IntervalToResendTxsInSeconds,
		TokensMapper:                 tokensMapper,
		RoleProvider:                 components.multiversXRoleProvider,
		StatusHandler:                args.MultiversXClientStatusHandler,
		AllowDelta:                   uint64(chainConfigs.ProxyMaxNoncesDelta),
	}

	components.multiversXClient, err = multiversx.NewClient(clientArgs)
	components.addClosableComponent(components.multiversXClient)

	return err
}

func (components *solanaMultiversXBridgeComponents) createSolanaClient(args ArgsSolanaToMultiversXBridge) error {
	solanaConfigs := args.Configs.GeneralConfig.Solana

	//GasStation is disabled by default
	gs, err := factory.CreateGasStation()
	if err != nil {
		return err
	}

	antifloodComponents, err := components.createAntifloodComponents(args.Configs.GeneralConfig.P2P.AntifloodConfig)
	if err != nil {
		return err
	}

	peerDenialEvaluator, err := p2p.NewPeerDenialEvaluator(antifloodComponents.BlacklistHandler, antifloodComponents.PubKeysCacher)
	if err != nil {
		return err
	}
	err = args.Messenger.SetPeerDenialEvaluator(peerDenialEvaluator)
	if err != nil {
		return err
	}

	broadcasterLogId := components.compatibleChain.BroadcasterLogId()
	solToMultiversXName := components.compatibleChain.SolanaToMultiversXName()
	argsBroadcaster := p2p.ArgsBroadcaster{
		Messenger:              args.Messenger,
		Log:                    core.NewLoggerWithIdentifier(logger.GetOrCreate(broadcasterLogId), broadcasterLogId),
		MultiversXRoleProvider: components.multiversXRoleProvider,
		SignatureProcessor:     components.solanaRoleProvider,
		KeyGen:                 keyGen,
		SingleSigner:           singleSigner,
		PrivateKey:             components.multiversXRelayerPrivateKey,
		Name:                   solToMultiversXName,
		AntifloodComponents:    antifloodComponents,
	}

	components.broadcaster, err = p2p.NewBroadcaster(argsBroadcaster)
	if err != nil {
		return err
	}

	privateKeyBytes, err := ioutil.ReadFile(solanaConfigs.PrivateKeyFile)
	if err != nil {
		return err
	}
	privateKeyString := converters.TrimWhiteSpaceCharacters(string(privateKeyBytes))
	privateKey, err := ethCrypto.HexToECDSA(privateKeyString)
	if err != nil {
		return err
	}

	tokensMapper, err := mappers.NewErc20ToMultiversXMapper(components.mxDataGetter)
	if err != nil {
		return err
	}

	signaturesHolder := solmultiversx.NewSignatureHolder()
	components.solToMultiversXSignaturesHolder = signaturesHolder
	err = components.broadcaster.AddBroadcastClient(signaturesHolder)
	if err != nil {
		return err
	}

	solClientLogId := components.compatibleChain.SolanaClientLogId()
	argsSolClient := sol.ArgsSolanaClient{
		SolanaRpcClient:         args.SolanaRpcClient,
		SolanaWssClientAddress:  args.SolanaWssClientAddress,
		DecimalDiffCalculator:   args.DecimalDiffCalculator,
		ProgramAddresses:        args.ProgramAddresses,
		StatusHandler:           components.solStatusHandler,
		RelayerPrivateKey:       args.RelayerPrivateKey,
		Log:                     core.NewLoggerWithIdentifier(logger.GetOrCreate(solClientLogId), solClientLogId),
		AddressConverter:        components.addressConverter,
		Broadcaster:             components.broadcaster,
		PrivateKey:              privateKey,
		TokensMapper:            tokensMapper,
		SignatureHolder:         signaturesHolder,
		GasHandler:              gs,
		TransferGasLimitBase:    solanaConfigs.GasLimitBase,
		TransferGasLimitForEach: solanaConfigs.GasLimitForEach,
		AllowDelta:              solanaConfigs.MaxBlocksDelta,
		SftTokenProgram:         components.SftTokenProgram,
	}

	components.solanaClient, err = sol.NewSolanaClient(argsSolClient)

	return err
}

func (components *solanaMultiversXBridgeComponents) createMultiversXRoleProvider(args ArgsSolanaToMultiversXBridge) error {
	configs := args.Configs.GeneralConfig
	multiversXRoleProviderLogId := components.compatibleChain.MultiversXRoleProviderLogId()
	log := core.NewLoggerWithIdentifier(logger.GetOrCreate(multiversXRoleProviderLogId), multiversXRoleProviderLogId)

	argsRoleProvider := roleproviders.ArgsMultiversXRoleProvider{
		DataGetter: components.mxDataGetter,
		Log:        log,
	}

	var err error
	components.multiversXRoleProvider, err = roleproviders.NewMultiversXRoleProvider(argsRoleProvider)
	if err != nil {
		return err
	}

	argsPollingHandler := polling.ArgsPollingHandler{
		Log:              log,
		Name:             "MultiversX role provider",
		PollingInterval:  time.Duration(configs.Relayer.RoleProvider.PollingIntervalInMillis) * time.Millisecond,
		PollingWhenError: pollingDurationOnError,
		Executor:         components.multiversXRoleProvider,
	}

	pollingHandler, err := polling.NewPollingHandler(argsPollingHandler)
	if err != nil {
		return err
	}

	components.addClosableComponent(pollingHandler)
	components.pollingHandlers = append(components.pollingHandlers, pollingHandler)

	return nil
}

func (components *solanaMultiversXBridgeComponents) createSolanaRoleProvider(args ArgsSolanaToMultiversXBridge) error {
	configs := args.Configs.GeneralConfig
	solRoleProviderLogId := components.compatibleChain.SolanaRoleProviderLogId()
	log := core.NewLoggerWithIdentifier(logger.GetOrCreate(solRoleProviderLogId), solRoleProviderLogId)
	argsRoleProvider := roleproviders.ArgsSolanaRoleProvider{
		SolanaRpcClient:  args.SolanaRpcClient,
		ProgramAddresses: args.ProgramAddresses,
		Log:              log,
	}

	var err error
	components.solanaRoleProvider, err = roleproviders.NewSolanaRoleProvider(argsRoleProvider)
	if err != nil {
		return err
	}

	argsPollingHandler := polling.ArgsPollingHandler{
		Log:              log,
		Name:             string(components.compatibleChain) + " role provider",
		PollingInterval:  time.Duration(configs.Relayer.RoleProvider.PollingIntervalInMillis) * time.Millisecond,
		PollingWhenError: pollingDurationOnError,
		Executor:         components.solanaRoleProvider,
	}

	pollingHandler, err := polling.NewPollingHandler(argsPollingHandler)
	if err != nil {
		return err
	}

	components.addClosableComponent(pollingHandler)
	components.pollingHandlers = append(components.pollingHandlers, pollingHandler)

	return nil
}

func (components *solanaMultiversXBridgeComponents) createSolanaToMultiversXBridge(args ArgsSolanaToMultiversXBridge) error {
	solToMultiversXName := components.compatibleChain.SolanaToMultiversXName()
	log := core.NewLoggerWithIdentifier(logger.GetOrCreate(solToMultiversXName), solToMultiversXName)

	configs, found := args.Configs.GeneralConfig.StateMachine[solToMultiversXName]
	if !found {
		return fmt.Errorf("%w for %q", errMissingConfig, solToMultiversXName)
	}

	components.solToMultiversXStepDuration = time.Duration(configs.StepDurationInMillis) * time.Millisecond

	argsTopologyHandler := topology.ArgsTopologyHandler{
		PublicKeysProvider: components.multiversXRoleProvider,
		Timer:              components.timer,
		IntervalForLeader:  time.Second * time.Duration(configs.IntervalForLeaderInSeconds),
		AddressBytes:       components.multiversXRelayerAddress.AddressBytes(),
		Log:                log,
		AddressConverter:   components.addressConverter,
	}

	topologyHandler, err := topology.NewTopologyHandler(argsTopologyHandler)
	if err != nil {
		return err
	}

	components.solToMultiversXStatusHandler, err = status.NewStatusHandler(solToMultiversXName, components.statusStorer)
	if err != nil {
		return err
	}

	err = components.metricsHolder.AddStatusHandler(components.solToMultiversXStatusHandler)
	if err != nil {
		return err
	}

	timeForTransferExecution := time.Second * time.Duration(args.Configs.GeneralConfig.Solana.IntervalToWaitForTransferInSeconds)

	batchValidator, err := components.createBatchValidator(components.compatibleChain, chain.MultiversX, args.Configs.GeneralConfig.BatchValidator)
	if err != nil {
		return err
	}

	argsBridgeExecutor := solmultiversx.ArgsBridgeExecutor{
		Log:                          log,
		TopologyProvider:             topologyHandler,
		MultiversXClient:             components.multiversXClient,
		SolanaClient:                 components.solanaClient,
		StatusHandler:                components.solToMultiversXStatusHandler,
		TimeForWaitOnSolana:          timeForTransferExecution,
		SignaturesHolder:             disabled.NewDisabledSignaturesHolder(),
		BatchValidator:               batchValidator,
		MaxQuorumRetriesOnSolana:     args.Configs.GeneralConfig.Solana.MaxRetriesOnQuorumReached,
		MaxQuorumRetriesOnMultiversX: args.Configs.GeneralConfig.MultiversX.MaxRetriesOnQuorumReached,
		MaxRestriesOnWasProposed:     args.Configs.GeneralConfig.MultiversX.MaxRetriesOnWasTransferProposed,
	}

	bridge, err := solmultiversx.NewBridgeExecutor(argsBridgeExecutor)
	if err != nil {
		return err
	}

	components.solToMultiversXMachineStates, err = soltomultiversx.CreateSteps(bridge)
	if err != nil {
		return err
	}

	return nil
}

func (components *solanaMultiversXBridgeComponents) createMultiversXToSolanaBridge(args ArgsSolanaToMultiversXBridge) error {
	multiversXToSolName := components.compatibleChain.MultiversXToSolanaName()
	log := core.NewLoggerWithIdentifier(logger.GetOrCreate(multiversXToSolName), multiversXToSolName)

	configs, found := args.Configs.GeneralConfig.StateMachine[multiversXToSolName]
	if !found {
		return fmt.Errorf("%w for %q", errMissingConfig, multiversXToSolName)
	}

	components.multiversXToSolStepDuration = time.Duration(configs.StepDurationInMillis) * time.Millisecond
	argsTopologyHandler := topology.ArgsTopologyHandler{
		PublicKeysProvider: components.multiversXRoleProvider,
		Timer:              components.timer,
		IntervalForLeader:  time.Second * time.Duration(configs.IntervalForLeaderInSeconds),
		AddressBytes:       components.multiversXRelayerAddress.AddressBytes(),
		Log:                log,
		AddressConverter:   components.addressConverter,
	}

	topologyHandler, err := topology.NewTopologyHandler(argsTopologyHandler)
	if err != nil {
		return err
	}

	components.multiversXToSolStatusHandler, err = status.NewStatusHandler(multiversXToSolName, components.statusStorer)
	if err != nil {
		return err
	}

	err = components.metricsHolder.AddStatusHandler(components.multiversXToSolStatusHandler)
	if err != nil {
		return err
	}

	timeForWaitOnSolana := time.Second * time.Duration(args.Configs.GeneralConfig.Solana.IntervalToWaitForTransferInSeconds)

	batchValidator, err := components.createBatchValidator(chain.MultiversX, components.compatibleChain, args.Configs.GeneralConfig.BatchValidator)
	if err != nil {
		return err
	}

	argsBridgeExecutor := solmultiversx.ArgsBridgeExecutor{
		Log:                          log,
		TopologyProvider:             topologyHandler,
		MultiversXClient:             components.multiversXClient,
		SolanaClient:                 components.solanaClient,
		StatusHandler:                components.multiversXToSolStatusHandler,
		TimeForWaitOnSolana:          timeForWaitOnSolana,
		SignaturesHolder:             components.solToMultiversXSignaturesHolder,
		BatchValidator:               batchValidator,
		MaxQuorumRetriesOnSolana:     args.Configs.GeneralConfig.Solana.MaxRetriesOnQuorumReached,
		MaxQuorumRetriesOnMultiversX: args.Configs.GeneralConfig.MultiversX.MaxRetriesOnQuorumReached,
		MaxRestriesOnWasProposed:     args.Configs.GeneralConfig.MultiversX.MaxRetriesOnWasTransferProposed,
	}

	bridge, err := solmultiversx.NewBridgeExecutor(argsBridgeExecutor)
	if err != nil {
		return err
	}

	components.multiversXToSolMachineStates, err = multiversxtosol.CreateSteps(bridge)
	if err != nil {
		return err
	}

	return nil
}

func (components *solanaMultiversXBridgeComponents) startPollingHandlers() error {
	for _, pollingHandler := range components.pollingHandlers {
		err := pollingHandler.StartProcessingLoop()
		if err != nil {
			return err
		}
	}

	return nil
}

// Start will start the bridge
func (components *solanaMultiversXBridgeComponents) Start() error {
	err := components.messenger.Bootstrap()
	if err != nil {
		return err
	}

	components.baseLogger.Info("waiting for p2p bootstrap", "time", components.timeForBootstrap)
	time.Sleep(components.timeForBootstrap)

	err = components.broadcaster.RegisterOnTopics()
	if err != nil {
		return err
	}

	components.broadcaster.BroadcastJoinTopic()

	err = components.startPollingHandlers()
	if err != nil {
		return err
	}

	go components.startBroadcastJoinRetriesLoop()

	return nil
}

func (components *solanaMultiversXBridgeComponents) createBatchValidator(sourceChain chain.Chain, destinationChain chain.Chain, args config.BatchValidatorConfig) (clients.BatchValidator, error) {
	argsBatchValidator := batchValidatorManagement.ArgsBatchValidator{
		SourceChain:      sourceChain,
		DestinationChain: destinationChain,
		RequestURL:       args.URL,
		RequestTime:      time.Second * time.Duration(args.RequestTimeInSeconds),
	}

	batchValidator, err := batchManagementFactory.CreateBatchValidator(argsBatchValidator, args.Enabled)
	if err != nil {
		return nil, err
	}
	return batchValidator, err
}

func (components *solanaMultiversXBridgeComponents) createSolanaToMultiversXStateMachine() error {
	solToMultiversXName := components.compatibleChain.SolanaToMultiversXName()
	log := core.NewLoggerWithIdentifier(logger.GetOrCreate(solToMultiversXName), solToMultiversXName)

	argsStateMachine := stateMachine.ArgsStateMachine{
		StateMachineName:     solToMultiversXName,
		Steps:                components.solToMultiversXMachineStates,
		StartStateIdentifier: soltomultiversx.GettingPendingBatchFromSolana,
		Log:                  log,
		StatusHandler:        components.solToMultiversXStatusHandler,
	}

	var err error
	components.solToMultiversXStateMachine, err = stateMachine.NewStateMachine(argsStateMachine)
	if err != nil {
		return err
	}

	argsPollingHandler := polling.ArgsPollingHandler{
		Log:              log,
		Name:             solToMultiversXName + " State machine",
		PollingInterval:  components.solToMultiversXStepDuration,
		PollingWhenError: pollingDurationOnError,
		Executor:         components.solToMultiversXStateMachine,
	}

	pollingHandler, err := polling.NewPollingHandler(argsPollingHandler)
	if err != nil {
		return err
	}

	components.addClosableComponent(pollingHandler)
	components.pollingHandlers = append(components.pollingHandlers, pollingHandler)

	return nil
}

func (components *solanaMultiversXBridgeComponents) createMultiversXToSolanaStateMachine() error {
	multiversXToSolName := components.compatibleChain.MultiversXToSolanaName()
	log := core.NewLoggerWithIdentifier(logger.GetOrCreate(multiversXToSolName), multiversXToSolName)

	argsStateMachine := stateMachine.ArgsStateMachine{
		StateMachineName:     multiversXToSolName,
		Steps:                components.multiversXToSolMachineStates,
		StartStateIdentifier: multiversxtosol.GettingPendingBatchFromMultiversX,
		Log:                  log,
		StatusHandler:        components.multiversXToSolStatusHandler,
	}

	var err error
	components.multiversXToSolStateMachine, err = stateMachine.NewStateMachine(argsStateMachine)
	if err != nil {
		return err
	}

	argsPollingHandler := polling.ArgsPollingHandler{
		Log:              log,
		Name:             multiversXToSolName + " State machine",
		PollingInterval:  components.multiversXToSolStepDuration,
		PollingWhenError: pollingDurationOnError,
		Executor:         components.multiversXToSolStateMachine,
	}

	pollingHandler, err := polling.NewPollingHandler(argsPollingHandler)
	if err != nil {
		return err
	}

	components.addClosableComponent(pollingHandler)
	components.pollingHandlers = append(components.pollingHandlers, pollingHandler)

	return nil
}

func (components *solanaMultiversXBridgeComponents) createAntifloodComponents(antifloodConfig chainConfig.AntifloodConfig) (*antifloodFactory.AntiFloodComponents, error) {
	var err error
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer func() {
		if err != nil {
			cancelFunc()
		}
	}()

	cfg := chainConfig.Config{
		Antiflood: antifloodConfig,
	}
	antiFloodComponents, err := antifloodFactory.NewP2PAntiFloodComponents(ctx, cfg, components.appStatusHandler, components.messenger.ID())
	if err != nil {
		return nil, err
	}
	return antiFloodComponents, nil
}

func (components *solanaMultiversXBridgeComponents) startBroadcastJoinRetriesLoop() {
	broadcastTimer := time.NewTimer(components.timeBeforeRepeatJoin)
	defer broadcastTimer.Stop()

	var ctx context.Context
	ctx, components.cancelFunc = context.WithCancel(context.Background())
	for {
		broadcastTimer.Reset(components.timeBeforeRepeatJoin)

		select {
		case <-broadcastTimer.C:
			components.baseLogger.Info("broadcast again join topic")
			components.broadcaster.BroadcastJoinTopic()
		case <-ctx.Done():
			components.baseLogger.Info("closing broadcast join topic loop")
			return

		}
	}
}

// Close will close any sub-components started
func (components *solanaMultiversXBridgeComponents) Close() error {
	components.mutClosableHandlers.RLock()
	defer components.mutClosableHandlers.RUnlock()

	if components.cancelFunc != nil {
		components.cancelFunc()
	}

	var lastError error
	for _, closable := range components.closableHandlers {
		if closable == nil {
			components.baseLogger.Warn("programming error, nil closable component")
			continue
		}

		err := closable.Close()
		if err != nil {
			lastError = err

			components.baseLogger.Error("error closing component", "error", err)
		}
	}

	return lastError
}

// MultiversXRelayerAddress returns the MultiversX's address associated to this relayer
func (components *solanaMultiversXBridgeComponents) MultiversXRelayerAddress() sdkCore.AddressHandler {
	return components.multiversXRelayerAddress
}
