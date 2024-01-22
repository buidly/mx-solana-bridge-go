package multiversxtosol

import (
	"fmt"
	solmultiversx "github.com/multiversx/mx-solana-bridge-go/bridges/solanaMultiversX"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-solana-bridge-go/bridges/solanaMultiversX/steps"
	"github.com/multiversx/mx-solana-bridge-go/core"
)

// CreateSteps creates all machine states providing the bridge executor
func CreateSteps(executor steps.Executor) (core.MachineStates, error) {
	if check.IfNil(executor) {
		return nil, solmultiversx.ErrNilExecutor
	}

	return createMachineStates(executor)
}

func createMachineStates(executor steps.Executor) (core.MachineStates, error) {
	machineStates := make(core.MachineStates)

	stepsSlice := []core.Step{
		&getPendingStep{
			bridge: executor,
		},
		&signProposedTransferStep{
			bridge: executor,
		},
		&waitForQuorumOnTransferStep{
			bridge: executor,
		},
		&performTransferStep{
			bridge: executor,
		},
		&waitTransferConfirmationStep{
			bridge: executor,
		},
		&resolveSetStatusStep{
			bridge: executor,
		},
		&proposeSetStatusStep{
			bridge: executor,
		},
		&signProposedSetStatusStep{
			bridge: executor,
		},
		&waitForQuorumOnSetStatusStep{
			bridge: executor,
		},
		&performSetStatusStep{
			bridge: executor,
		},
	}

	for _, s := range stepsSlice {
		_, found := machineStates[s.Identifier()]
		if found {
			return nil, fmt.Errorf("%w for identifier '%s'", solmultiversx.ErrDuplicatedStepIdentifier, s.Identifier())
		}

		machineStates[s.Identifier()] = s
	}

	return machineStates, nil
}
