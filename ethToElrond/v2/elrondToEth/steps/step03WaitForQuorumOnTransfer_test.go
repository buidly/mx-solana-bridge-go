package steps

import (
	"context"
	"testing"

	"github.com/ElrondNetwork/elrond-eth-bridge/core"
	"github.com/ElrondNetwork/elrond-eth-bridge/ethToElrond/v2/elrondToEth"
	"github.com/ElrondNetwork/elrond-eth-bridge/testsCommon/bridgeV2"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/stretchr/testify/assert"
)

func TestExecute_WaitForQuorumOnTransfer(t *testing.T) {
	t.Parallel()

	t.Run("error on IsQuorumReachedOnEthereum", func(t *testing.T) {
		t.Parallel()
		bridgeStub := createStubExecutorWaitForQuorumOnTransfer()
		bridgeStub.IsQuorumReachedOnEthereumCalled = func(ctx context.Context) (bool, error) {
			return false, expectedError
		}

		step := waitForQuorumOnTransferStep{
			bridge: bridgeStub,
		}

		stepIdentifier, err := step.Execute(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, initialStep, stepIdentifier)
	})

	t.Run("max retries reached", func(t *testing.T) {
		t.Parallel()
		bridgeStub := createStubExecutorWaitForQuorumOnTransfer()
		bridgeStub.ProcessMaxRetriesOnEthereumCalled = func() bool {
			return true
		}

		step := waitForQuorumOnTransferStep{
			bridge: bridgeStub,
		}

		stepIdentifier, err := step.Execute(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, initialStep, stepIdentifier)
	})

	t.Run("quorum not reached", func(t *testing.T) {
		t.Parallel()
		bridgeStub := createStubExecutorWaitForQuorumOnTransfer()
		bridgeStub.IsQuorumReachedOnEthereumCalled = func(ctx context.Context) (bool, error) {
			return false, nil
		}

		step := waitForQuorumOnTransferStep{
			bridge: bridgeStub,
		}

		assert.False(t, step.IsInterfaceNil())

		expectedStepIdentifier := step.Identifier()
		stepIdentifier, err := step.Execute(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, expectedStepIdentifier, stepIdentifier)
	})

	t.Run("quorum reached", func(t *testing.T) {
		t.Parallel()
		bridgeStub := createStubExecutorWaitForQuorumOnTransfer()
		bridgeStub.IsQuorumReachedOnEthereumCalled = func(ctx context.Context) (bool, error) {
			return true, nil
		}

		step := waitForQuorumOnTransferStep{
			bridge: bridgeStub,
		}

		assert.False(t, step.IsInterfaceNil())

		expectedStepIdentifier := core.StepIdentifier(elrondToEth.PerformingTransfer)
		stepIdentifier, err := step.Execute(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, expectedStepIdentifier, stepIdentifier)
	})
}

func createStubExecutorWaitForQuorumOnTransfer() *bridgeV2.ElrondToEthBridgeStub {
	stub := bridgeV2.NewElrondToEthBridgeStub()
	stub.GetLoggerCalled = func() logger.Logger {
		return testLogger
	}
	stub.ProcessMaxRetriesOnEthereumCalled = func() bool {
		return false
	}
	return stub
}
