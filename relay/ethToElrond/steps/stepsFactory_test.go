package steps

import (
	"testing"

	"github.com/ElrondNetwork/elrond-eth-bridge/relay/ethToElrond/steps/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSteps_Errors(t *testing.T) {
	t.Parallel()

	steps, err := CreateSteps(nil)

	assert.Nil(t, steps)
	assert.Equal(t, ErrNilBridgeExecutor, err)
}

func TestCreateSteps_ShouldWork(t *testing.T) {
	t.Parallel()

	steps, err := CreateSteps(mock.NewBridgeExecutorMock())

	require.NotNil(t, steps)
	require.Nil(t, err)
	require.Equal(t, 7, len(steps))
}
