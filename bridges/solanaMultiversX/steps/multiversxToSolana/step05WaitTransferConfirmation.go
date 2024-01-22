package multiversxtosol

import (
	"context"

	"github.com/multiversx/mx-solana-bridge-go/bridges/solanaMultiversX/steps"
	"github.com/multiversx/mx-solana-bridge-go/core"
)

type waitTransferConfirmationStep struct {
	bridge steps.Executor
}

// Execute will execute this step returning the next step to be executed
func (step *waitTransferConfirmationStep) Execute(ctx context.Context) core.StepIdentifier {
	wasPerformedAndFinal := step.bridge.WaitForTransferConfirmation(ctx)
	if wasPerformedAndFinal {
		return ResolvingSetStatusOnMultiversX
	}
	return PerformingTransfer
}

// Identifier returns the step's identifier
func (step *waitTransferConfirmationStep) Identifier() core.StepIdentifier {
	return WaitingTransferConfirmation
}

// IsInterfaceNil returns true if there is no value under the interface
func (step *waitTransferConfirmationStep) IsInterfaceNil() bool {
	return step == nil
}
