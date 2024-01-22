package clients

import (
	"context"
	"encoding/json"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"net/http"
)

type DecimalDiffCalculatorArgs struct {
	SolanaRpcClient *rpc.Client
	MvxProxy        Proxy
}

type decimalDiffCalculator struct {
	solanaRpcClient *rpc.Client
	mvxProxy        Proxy
}

type TokenDetailsResponse struct {
	Identifier string `json:"identifier"`
	Decimals   uint8  `json:"decimals"`
	Ticker     string `json:"ticker"`
}

func NewDecimalDiffCalculator(args DecimalDiffCalculatorArgs) *decimalDiffCalculator {
	return &decimalDiffCalculator{
		solanaRpcClient: args.SolanaRpcClient,
		mvxProxy:        args.MvxProxy,
	}
}

func (c *decimalDiffCalculator) GetDecimalDifference(ctx context.Context, solanaTokenAddress []byte, mvxTokenAddress []byte) (int, error) {
	solanaTokenDetails, err := c.GetSolanaTokenDetails(ctx, solanaTokenAddress)
	if err != nil {
		return -1, err
	}
	mvxTokenDetails, err := c.GetMvxTokenDetails(ctx, mvxTokenAddress)
	if err != nil {
		return -1, err
	}
	if solanaTokenDetails == nil {
		return -1, SolanaTokenDataNilError
	}
	if mvxTokenDetails == nil {
		return -1, MvxTokenDataNilError
	}

	return int(solanaTokenDetails.Decimals) - int(mvxTokenDetails.Decimals), nil
}

func (c *decimalDiffCalculator) GetMvxTokenDetails(ctx context.Context, tokenAddressBytes []byte) (*TokenDetailsResponse, error) {

	tokenIdentifier := string(tokenAddressBytes)
	buff, code, err := c.mvxProxy.GetHTTP(ctx, "tokens/"+tokenIdentifier)
	if err != nil || code != http.StatusOK {
		return nil, err
	}

	response := &TokenDetailsResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *decimalDiffCalculator) GetSolanaTokenDetails(ctx context.Context, tokenAddressBytes []byte) (*token.Mint, error) {
	pubKey := solana.PublicKeyFromBytes(tokenAddressBytes)
	var mint token.Mint
	err := c.solanaRpcClient.GetAccountDataInto(
		ctx,
		pubKey,
		&mint,
	)
	if err != nil {
		return nil, err
	}
	return &mint, nil
}
