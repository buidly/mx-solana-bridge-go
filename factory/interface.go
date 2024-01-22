package factory

import (
	"context"

	sdkCore "github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-solana-bridge-go/core"
)

type dataGetter interface {
	GetTokenIdForSftAddress(ctx context.Context, sftAddress []byte) ([][]byte, error)
	GetSftAddressForTokenId(ctx context.Context, tokenId []byte) ([][]byte, error)
	GetAllStakedRelayers(ctx context.Context) ([][]byte, error)
	IsInterfaceNil() bool
}

// MultiversXRoleProvider defines the operations for the MultiversX role provider
type MultiversXRoleProvider interface {
	Execute(ctx context.Context) error
	IsWhitelisted(address sdkCore.AddressHandler) bool
	SortedPublicKeys() [][]byte
	IsInterfaceNil() bool
}

// SolanaRoleProvider defines the operations for the Solana role provider
type SolanaRoleProvider interface {
	Execute(ctx context.Context) error
	VerifyEthSignature(signature []byte, messageHash []byte) error
	IsInterfaceNil() bool
}

// Broadcaster defines a component able to communicate with other such instances and manage signatures and other state related data
type Broadcaster interface {
	BroadcastSignature([]*core.EthereumSignature)
	BroadcastJoinTopic()
	SortedPublicKeys() [][]byte
	RegisterOnTopics() error
	AddBroadcastClient(client core.BroadcastClient) error
	Close() error
	IsInterfaceNil() bool
}

// StateMachine defines a state machine component
type StateMachine interface {
	Execute(ctx context.Context) error
	IsInterfaceNil() bool
}

// PollingHandler defines a polling handler component
type PollingHandler interface {
	StartProcessingLoop() error
	IsInterfaceNil() bool
}
