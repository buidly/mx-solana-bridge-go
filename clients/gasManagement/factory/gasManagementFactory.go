package factory

import (
	"github.com/multiversx/mx-solana-bridge-go/clients"
	"github.com/multiversx/mx-solana-bridge-go/clients/gasManagement/disabled"
)

// CreateGasStation generates an implementation of GasHandler
func CreateGasStation() (clients.GasHandler, error) {
	return &disabled.DisabledGasStation{}, nil
}
