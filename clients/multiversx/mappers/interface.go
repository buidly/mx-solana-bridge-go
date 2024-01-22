package mappers

import "context"

// DataGetter defines the interface able to handle get requests for MultiversX blockchain
type DataGetter interface {
	GetTokenIdForSftAddress(ctx context.Context, erc20Address []byte) ([][]byte, error)
	GetSftAddressForTokenId(ctx context.Context, tokenId []byte) ([][]byte, error)
	IsInterfaceNil() bool
}
