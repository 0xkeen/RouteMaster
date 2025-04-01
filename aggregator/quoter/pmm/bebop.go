package pmm

import (
	"fmt"
	"math/big"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"go.uber.org/ratelimit"

	"route_master/aggregator/quoter"
	"route_master/clients"
	"route_master/commons"
)

// Verify that BebopQuoter implements the Quoter interface

type BebopQuoter struct {
	limiter ratelimit.Limiter
	client  *clients.BebopClient
}

func NewBebopQuoter() *BebopQuoter {
	return &BebopQuoter{
		limiter: ratelimit.New(1),
		client:  clients.NewBebopClient(),
	}
}

func (b *BebopQuoter) ExactIn(
	fromToken commons.Token,
	toToken commons.Token,
	amountIn string,
	requestId string,
	route interface{},
	filler common.Address,
) (*quoter.Quote, error) {
	quoteRequest := clients.BebopQuoteRequest{
		SellTokens:     []string{fromToken.Address.String()},
		BuyTokens:      []string{toToken.Address.String()},
		SellAmounts:    []string{amountIn},
		TakerAddress:   filler.String(),
		SkipValidation: true,
		ApprovalType:   "Standard",
		Gasless:        false,
		Receiver:       filler.String(),
	}

	b.limiter.Take()
	var response *clients.BebopQuoteResponse
	err := retry.Do(
		func() error {
			var err error
			response, err = b.client.GetQuote(quoteRequest)
			return err
		},
		retry.Attempts(1),
		retry.Delay(100*time.Millisecond),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote after retries: %w", err)
	}

	outputAmountWei := new(big.Int)
	outputAmountWei.SetString(response.BuyTokens[toToken.Address.String()].Amount, 10)

	callData, _ := hexutil.Decode(response.Tx.Data)
	quote := quoter.Quote{
		RequestId: requestId,
		Timestamp: time.Now().Unix(),
		Prices: map[string]string{
			amountIn: outputAmountWei.String(),
		},
		Routes: map[string]interface{}{
			"gas":              response.Tx.Gas,
			"expiry":           response.Expiry,
			"tx":               response.Tx,
			"onchainOrderType": response.OnchainOrderType,
		},
		Tx:       response.Tx,
		Protocol: b.String(),
		Path:     []common.Address{fromToken.Address, toToken.Address},
		To:       common.HexToAddress(response.Tx.To),
		Calldata: callData,
	}

	return &quote, nil
}

func (b *BebopQuoter) ExactOut(
	fromToken commons.Token,
	toToken commons.Token,
	amountOut string,
	requestId string,
	route interface{},
	filler common.Address,
) (*quoter.Quote, error) {
	return nil, fmt.Errorf("bebop exactOut not implemented")
}

func (b *BebopQuoter) String() string {
	return "bebop"
}
