package roleproviders

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

// DataGetter defines the interface able to handle get requests for MultiversX blockchain
type DataGetter interface {
	GetAllStakedRelayers(ctx context.Context) ([][]byte, error)
	IsInterfaceNil() bool
}

// SolanaChainInteractor defines an Solana client able to respond to requests
type SolanaChainInteractor interface {
	GetRelayers(ctx context.Context) ([]common.Address, error)
	IsInterfaceNil() bool
}
