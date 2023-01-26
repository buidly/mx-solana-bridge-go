package factory

import (
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-eth-bridge/clients/chain"
	"github.com/ElrondNetwork/elrond-eth-bridge/config"
	"github.com/ElrondNetwork/elrond-eth-bridge/core"
	"github.com/ElrondNetwork/elrond-eth-bridge/status"
	"github.com/ElrondNetwork/elrond-eth-bridge/testsCommon"
	bridgeTests "github.com/ElrondNetwork/elrond-eth-bridge/testsCommon/bridge"
	p2pMocks "github.com/ElrondNetwork/elrond-eth-bridge/testsCommon/p2p"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/testscommon/statusHandler"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	erdgoCore "github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createMockEthElrondBridgeArgs() ArgsEthereumToElrondBridge {
	stateMachineConfig := config.ConfigStateMachine{
		StepDurationInMillis:       1000,
		IntervalForLeaderInSeconds: 60,
	}

	cfg := config.Config{
		Eth: config.EthereumConfig{
			Chain:                        chain.Ethereum,
			NetworkAddress:               "http://127.0.0.1:8545",
			SafeContractAddress:          "5DdDe022a65F8063eE9adaC54F359CBF46166068",
			PrivateKeyFile:               "testdata/grace.sk",
			IntervalToResendTxsInSeconds: 0,
			GasLimitBase:                 200000,
			GasLimitForEach:              30000,
			GasStation: config.GasStationConfig{
				Enabled:                    true,
				URL:                        "",
				PollingIntervalInSeconds:   1,
				RequestRetryDelayInSeconds: 1,
				MaxFetchRetries:            3,
				RequestTimeInSeconds:       1,
				MaximumAllowedGasPrice:     100,
				GasPriceSelector:           "FastGasPrice",
				GasPriceMultiplier:         1,
			},
			MaxRetriesOnQuorumReached:          1,
			IntervalToWaitForTransferInSeconds: 1,
			MaxBlocksDelta:                     10,
		},
		Elrond: config.ElrondConfig{
			PrivateKeyFile:                  "testdata/grace.pem",
			IntervalToResendTxsInSeconds:    60,
			NetworkAddress:                  "http://127.0.0.1:8079",
			MultisigContractAddress:         "erd1qqqqqqqqqqqqqpgqgftcwj09u0nhmskrw7xxqcqh8qmzwyexd8ss7ftcxx",
			GasMap:                          testsCommon.CreateTestElrondGasMap(),
			MaxRetriesOnQuorumReached:       1,
			MaxRetriesOnWasTransferProposed: 1,
			ProxyMaxNoncesDelta:             5,
		},
		Relayer: config.ConfigRelayer{
			RoleProvider: config.RoleProviderConfig{
				PollingIntervalInMillis: 1000,
			},
		},
		StateMachine: map[string]config.ConfigStateMachine{
			"EthereumToElrond": stateMachineConfig,
			"ElrondToEthereum": stateMachineConfig,
		},
	}
	configs := config.Configs{
		GeneralConfig:   cfg,
		ApiRoutesConfig: config.ApiRoutesConfig{},
		FlagsConfig: config.ContextFlagsConfig{
			RestApiInterface: core.WebServerOffString,
		},
	}

	argsProxy := blockchain.ArgsElrondProxy{
		ProxyURL:            cfg.Elrond.NetworkAddress,
		CacheExpirationTime: time.Minute,
		EntityType:          erdgoCore.ObserverNode,
	}
	proxy, _ := blockchain.NewElrondProxy(argsProxy)
	return ArgsEthereumToElrondBridge{
		Configs:                   configs,
		Messenger:                 &p2pMocks.MessengerStub{},
		StatusStorer:              testsCommon.NewStorerMock(),
		Proxy:                     proxy,
		ElrondClientStatusHandler: &testsCommon.StatusHandlerStub{},
		Erc20ContractsHolder:      &bridgeTests.ERC20ContractsHolderStub{},
		ClientWrapper:             &bridgeTests.EthereumClientWrapperStub{},
		TimeForBootstrap:          minTimeForBootstrap,
		TimeBeforeRepeatJoin:      minTimeBeforeRepeatJoin,
		MetricsHolder:             status.NewMetricsHolder(),
		AppStatusHandler:          &statusHandler.AppStatusHandlerStub{},
	}
}

func TestNewEthElrondBridgeComponents(t *testing.T) {
	t.Parallel()

	t.Run("nil Proxy", func(t *testing.T) {
		t.Parallel()
		args := createMockEthElrondBridgeArgs()
		args.Proxy = nil

		components, err := NewEthElrondBridgeComponents(args)
		assert.Equal(t, errNilProxy, err)
		assert.Nil(t, components)
	})
	t.Run("nil Messenger", func(t *testing.T) {
		t.Parallel()
		args := createMockEthElrondBridgeArgs()
		args.Messenger = nil

		components, err := NewEthElrondBridgeComponents(args)
		assert.Equal(t, errNilMessenger, err)
		assert.Nil(t, components)
	})
	t.Run("nil ClientWrapper", func(t *testing.T) {
		t.Parallel()
		args := createMockEthElrondBridgeArgs()
		args.ClientWrapper = nil

		components, err := NewEthElrondBridgeComponents(args)
		assert.Equal(t, errNilEthClient, err)
		assert.Nil(t, components)
	})
	t.Run("nil StatusStorer", func(t *testing.T) {
		t.Parallel()
		args := createMockEthElrondBridgeArgs()
		args.StatusStorer = nil

		components, err := NewEthElrondBridgeComponents(args)
		assert.Equal(t, errNilStatusStorer, err)
		assert.Nil(t, components)
	})
	t.Run("nil Erc20ContractsHolder", func(t *testing.T) {
		t.Parallel()
		args := createMockEthElrondBridgeArgs()
		args.Erc20ContractsHolder = nil

		components, err := NewEthElrondBridgeComponents(args)
		assert.Equal(t, errNilErc20ContractsHolder, err)
		assert.Nil(t, components)
	})
	t.Run("err on createElrondKeysAndAddresses, empty pk file", func(t *testing.T) {
		t.Parallel()
		args := createMockEthElrondBridgeArgs()
		args.Configs.GeneralConfig.Elrond.PrivateKeyFile = ""

		components, err := NewEthElrondBridgeComponents(args)
		assert.NotNil(t, err)
		assert.Nil(t, components)
	})
	t.Run("err on createElrondKeysAndAddresses, empty multisig address", func(t *testing.T) {
		t.Parallel()
		args := createMockEthElrondBridgeArgs()
		args.Configs.GeneralConfig.Elrond.MultisigContractAddress = ""

		components, err := NewEthElrondBridgeComponents(args)
		assert.NotNil(t, err)
		assert.Nil(t, components)
	})
	t.Run("err on createElrondClient", func(t *testing.T) {
		t.Parallel()
		args := createMockEthElrondBridgeArgs()
		args.Configs.GeneralConfig.Elrond.GasMap = config.ElrondGasMapConfig{}

		components, err := NewEthElrondBridgeComponents(args)
		assert.NotNil(t, err)
		assert.Nil(t, components)
	})
	t.Run("err on createElrondRoleProvider", func(t *testing.T) {
		t.Parallel()
		args := createMockEthElrondBridgeArgs()
		args.Configs.GeneralConfig.Relayer.RoleProvider.PollingIntervalInMillis = 0

		components, err := NewEthElrondBridgeComponents(args)
		assert.NotNil(t, err)
		assert.Nil(t, components)
	})
	t.Run("err on createEthereumClient, empty eth config", func(t *testing.T) {
		t.Parallel()
		args := createMockEthElrondBridgeArgs()
		args.Configs.GeneralConfig.Eth = config.EthereumConfig{}

		components, err := NewEthElrondBridgeComponents(args)
		assert.NotNil(t, err)
		assert.Nil(t, components)
	})
	t.Run("err on createEthereumClient, invalid gas price selector", func(t *testing.T) {
		t.Parallel()
		args := createMockEthElrondBridgeArgs()
		args.Configs.GeneralConfig.Eth.GasStation.GasPriceSelector = core.WebServerOffString

		components, err := NewEthElrondBridgeComponents(args)
		assert.NotNil(t, err)
		assert.Nil(t, components)
	})
	t.Run("err missing state machine config", func(t *testing.T) {
		t.Parallel()
		args := createMockEthElrondBridgeArgs()
		args.Configs.GeneralConfig.StateMachine = make(map[string]config.ConfigStateMachine)

		components, err := NewEthElrondBridgeComponents(args)
		assert.True(t, errors.Is(err, errMissingConfig))
		assert.True(t, strings.Contains(err.Error(), args.Configs.GeneralConfig.Eth.Chain.EvmCompatibleChainToElrondName()))
		assert.Nil(t, components)
	})
	t.Run("invalid time for bootstrap", func(t *testing.T) {
		t.Parallel()
		args := createMockEthElrondBridgeArgs()
		args.TimeForBootstrap = minTimeForBootstrap - 1

		components, err := NewEthElrondBridgeComponents(args)
		assert.True(t, errors.Is(err, errInvalidValue))
		assert.True(t, strings.Contains(err.Error(), "for TimeForBootstrap"))
		assert.Nil(t, components)
	})
	t.Run("invalid time before retry", func(t *testing.T) {
		t.Parallel()
		args := createMockEthElrondBridgeArgs()
		args.TimeBeforeRepeatJoin = minTimeBeforeRepeatJoin - 1

		components, err := NewEthElrondBridgeComponents(args)
		assert.True(t, errors.Is(err, errInvalidValue))
		assert.True(t, strings.Contains(err.Error(), "for TimeBeforeRepeatJoin"))
		assert.Nil(t, components)
	})
	t.Run("nil MetricsHolder", func(t *testing.T) {
		t.Parallel()
		args := createMockEthElrondBridgeArgs()
		args.MetricsHolder = nil

		components, err := NewEthElrondBridgeComponents(args)
		assert.Equal(t, errNilMetricsHolder, err)
		assert.Nil(t, components)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()
		args := createMockEthElrondBridgeArgs()

		components, err := NewEthElrondBridgeComponents(args)
		require.Nil(t, err)
		require.NotNil(t, components)
		require.Equal(t, 6, len(components.closableHandlers))
		require.False(t, check.IfNil(components.ethToElrondStatusHandler))
		require.False(t, check.IfNil(components.elrondToEthStatusHandler))
	})
}

func TestEthElrondBridgeComponents_StartAndCloseShouldWork(t *testing.T) {
	t.Parallel()

	args := createMockEthElrondBridgeArgs()
	components, err := NewEthElrondBridgeComponents(args)
	assert.Nil(t, err)

	err = components.Start()
	assert.Nil(t, err)
	assert.Equal(t, 6, len(components.closableHandlers))

	time.Sleep(time.Second * 2) // allow go routines to start

	err = components.Close()
	assert.Nil(t, err)
}

func TestEthElrondBridgeComponents_Start(t *testing.T) {
	t.Parallel()

	t.Run("messenger errors on bootstrap", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected error")
		args := createMockEthElrondBridgeArgs()
		args.Messenger = &p2pMocks.MessengerStub{
			BootstrapCalled: func() error {
				return expectedErr
			},
		}
		components, _ := NewEthElrondBridgeComponents(args)

		err := components.Start()
		assert.Equal(t, expectedErr, err)
	})
	t.Run("broadcaster errors on RegisterOnTopics", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected error")
		args := createMockEthElrondBridgeArgs()
		components, _ := NewEthElrondBridgeComponents(args)
		components.broadcaster = &testsCommon.BroadcasterStub{
			RegisterOnTopicsCalled: func() error {
				return expectedErr
			},
		}

		err := components.Start()
		assert.Equal(t, expectedErr, err)
	})
}

func TestEthElrondBridgeComponents_Close(t *testing.T) {
	t.Parallel()

	t.Run("nil closable should not panic", func(t *testing.T) {
		t.Parallel()

		defer func() {
			r := recover()
			if r != nil {
				assert.Fail(t, fmt.Sprintf("should have not failed %v", r))
			}
		}()

		components := &ethElrondBridgeComponents{
			baseLogger: logger.GetOrCreate("test"),
		}
		components.addClosableComponent(nil)

		err := components.Close()
		assert.Nil(t, err)
	})
	t.Run("one component errors, should return error", func(t *testing.T) {
		t.Parallel()

		components := &ethElrondBridgeComponents{
			baseLogger: logger.GetOrCreate("test"),
		}

		expectedErr := errors.New("expected error")

		numCalls := 0
		components.addClosableComponent(&testsCommon.CloserStub{
			CloseCalled: func() error {
				numCalls++
				return nil
			},
		})
		components.addClosableComponent(&testsCommon.CloserStub{
			CloseCalled: func() error {
				numCalls++
				return expectedErr
			},
		})
		components.addClosableComponent(&testsCommon.CloserStub{
			CloseCalled: func() error {
				numCalls++
				return nil
			},
		})

		err := components.Close()
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, 3, numCalls)
	})
}

func TestEthElrondBridgeComponents_startBroadcastJoinRetriesLoop(t *testing.T) {
	t.Parallel()

	t.Run("close before minTimeBeforeRepeatJoin", func(t *testing.T) {
		t.Parallel()

		numberOfCalls := uint32(0)
		args := createMockEthElrondBridgeArgs()
		components, _ := NewEthElrondBridgeComponents(args)

		components.broadcaster = &testsCommon.BroadcasterStub{
			BroadcastJoinTopicCalled: func() {
				atomic.AddUint32(&numberOfCalls, 1)
			},
		}

		err := components.Start()
		assert.Nil(t, err)
		time.Sleep(time.Second * 3)

		err = components.Close()
		assert.Nil(t, err)
		assert.Equal(t, uint32(1), atomic.LoadUint32(&numberOfCalls)) // one call expected from Start
	})
	t.Run("broadcast should be called again", func(t *testing.T) {
		t.Parallel()

		numberOfCalls := uint32(0)
		args := createMockEthElrondBridgeArgs()
		components, _ := NewEthElrondBridgeComponents(args)
		components.timeBeforeRepeatJoin = time.Second * 3
		components.broadcaster = &testsCommon.BroadcasterStub{
			BroadcastJoinTopicCalled: func() {
				atomic.AddUint32(&numberOfCalls, 1)
			},
		}

		err := components.Start()
		assert.Nil(t, err)
		time.Sleep(time.Second * 7)

		err = components.Close()
		assert.Nil(t, err)
		assert.Equal(t, uint32(3), atomic.LoadUint32(&numberOfCalls)) // 3 calls expected: Start + 2 times from loop
	})
}

func TestEthElrondBridgeComponents_RelayerAddresses(t *testing.T) {
	t.Parallel()

	args := createMockEthElrondBridgeArgs()
	components, _ := NewEthElrondBridgeComponents(args)

	assert.Equal(t, "erd1r69gk66fmedhhcg24g2c5kn2f2a5k4kvpr6jfw67dn2lyydd8cfswy6ede", components.ElrondRelayerAddress().AddressAsBech32String())
	assert.Equal(t, "0x3FE464Ac5aa562F7948322F92020F2b668D543d8", components.EthereumRelayerAddress().String())
}
