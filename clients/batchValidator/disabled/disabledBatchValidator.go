package disabled

import (
	"context"

	"github.com/ElrondNetwork/elrond-eth-bridge/clients"
)

type DisabledBatchValidator struct{}

// ValidateBatch returns true,nil and will result in skipping batch validation
func (dbv *DisabledBatchValidator) ValidateBatch(_ context.Context, _ *clients.TransferBatch) (bool, error) {
	return true, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (dbv *DisabledBatchValidator) IsInterfaceNil() bool {
	return dbv == nil
}
