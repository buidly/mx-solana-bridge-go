package steps

import (
	"context"

	"github.com/ElrondNetwork/elrond-eth-bridge/core"
	"github.com/ElrondNetwork/elrond-eth-bridge/ethToElrond/v2/ethToElrond"
)

type performActionIDStep struct {
	bridge ethToElrond.EthToElrondBridge
}

// Execute will execute this step returning the next step to be executed
func (step *performActionIDStep) Execute(ctx context.Context) (core.StepIdentifier, error) {
	wasPerformed, err := step.bridge.WasActionIDPerformed(ctx)
	if err != nil {
		step.bridge.GetLogger().Error("error determining if the action ID was proposed or not",
			"action ID", step.bridge.GetStoredActionID(), "error", err)
		return ethToElrond.GettingPendingBatchFromEthereum, nil
	}

	if wasPerformed {
		step.bridge.GetLogger().Info("action ID performed",
			"action ID", step.bridge.GetStoredActionID())
		return ethToElrond.GettingPendingBatchFromEthereum, nil
	}

	if !step.bridge.MyTurnAsLeader() {
		return step.Identifier(), nil
	}

	err = step.bridge.PerformActionID(ctx)
	if err != nil {
		step.bridge.GetLogger().Info("error performing action ID",
			"action ID", step.bridge.GetStoredActionID(), "error", err)
		return ethToElrond.GettingPendingBatchFromEthereum, nil
	}

	return step.Identifier(), nil
}

// Identifier returns the step's identifier
func (step *performActionIDStep) Identifier() core.StepIdentifier {
	return ethToElrond.PerformingActionID
}

// IsInterfaceNil returns true if there is no value under the interface
func (step *performActionIDStep) IsInterfaceNil() bool {
	return step == nil
}
