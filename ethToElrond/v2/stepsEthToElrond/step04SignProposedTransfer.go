package stepsEthToElrond

import (
	"context"

	"github.com/ElrondNetwork/elrond-eth-bridge/core"
)

type signProposedTransferStep struct {
	bridge EthToElrondBridge
}

// Execute will execute this step returning the next step to be executed
func (step *signProposedTransferStep) Execute(ctx context.Context) (core.StepIdentifier, error) {
	batch := step.bridge.GetStoredBatch()

	wasSigned, err := step.bridge.WasProposedTransferSigned(ctx)
	if err != nil {
		step.bridge.GetLogger().Error("error determining if the proposed transfer was signed or not",
			"batch ID", batch.ID, "error", err)
		return GetPendingBatchFromEthereum, nil
	}

	if wasSigned {
		return WaitForQuorum, nil
	}

	err = step.bridge.SignProposedTransfer(ctx)
	if err != nil {
		step.bridge.GetLogger().Error("error signing the proposed transfer",
			"batch ID", batch.ID, "error", err)
		return GetPendingBatchFromEthereum, nil
	}

	return WaitForQuorum, nil
}

// Identifier returns the step's identifier
func (step *signProposedTransferStep) Identifier() core.StepIdentifier {
	return SignProposedTransferOnElrond
}

// IsInterfaceNil returns true if there is no value under the interface
func (step *signProposedTransferStep) IsInterfaceNil() bool {
	return step == nil
}
