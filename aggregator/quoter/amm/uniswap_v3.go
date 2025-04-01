package amm

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"route_master/aggregator/quoter"
	"route_master/commons"
	"route_master/config"
	"route_master/contracts/uniswap_v3"
	"route_master/networking"
)

type UniswapV3 struct {
	fee              *big.Int
	client           *networking.SolidEthClient
	quoterAddress    common.Address
	uniswapRouterABI *abi.ABI
	routerAddress    common.Address
}

func NewUniswapV3(client *networking.SolidEthClient, fee *big.Int) *UniswapV3 {
	quoterAddress := common.HexToAddress("0x61fFE014bA17989E743c5F6cB21bF9697530B21e")

	uniswapRouterABI, _ := abi.JSON(strings.NewReader(config.RouterABIs.V3RouterABI))

	return &UniswapV3{
		fee:              fee,
		client:           client,
		quoterAddress:    quoterAddress,
		uniswapRouterABI: &uniswapRouterABI,
		routerAddress:    common.HexToAddress("0xE592427A0AEce92De3Edee1F18E0157C05861564"),
	}
}

func (u *UniswapV3) ExactIn(
	fromToken commons.Token,
	toToken commons.Token,
	amountIn string,
	requestId string,
	route interface{},
	filler common.Address,
) (*quoter.Quote, error) {
	amount, ok := new(big.Int).SetString(amountIn, 10)
	if !ok {
		return nil, fmt.Errorf("invalid amount: %s", amountIn)
	}

	params := uniswap_v3.IQuoterQuoteExactInputSingleParams{
		TokenIn:           fromToken.Address,
		TokenOut:          toToken.Address,
		AmountIn:          amount,
		Fee:               u.fee,
		SqrtPriceLimitX96: big.NewInt(0),
	}

	//use solid client
	client, err := u.client.Client()
	if err != nil {
		return nil, fmt.Errorf("error getting eth client: %v", err)
	}
	quoterContract, _ := uniswap_v3.NewQuoterV3(u.quoterAddress, client)

	var result commons.QuoteExactInputResponse
	err = retry.Do(
		func() error {
			var err error
			result, err = quoterContract.QuoteExactInputSingle(&bind.CallOpts{}, params)
			return err
		},
		retry.Attempts(1),
		retry.Delay(100*time.Millisecond),
		retry.OnRetry(func(n uint, err error) {
			fmt.Printf("Retrying UniswapV3.QuoteExactInputSingle after error: %v\n", err)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote after retries: %w", err)
	}

	deadline := big.NewInt(time.Now().Add(30 * time.Second).Unix())
	inputStruct := struct {
		TokenIn           common.Address
		TokenOut          common.Address
		Fee               *big.Int
		Recipient         common.Address
		Deadline          *big.Int
		AmountIn          *big.Int
		AmountOutMinimum  *big.Int
		SqrtPriceLimitX96 *big.Int
	}{
		TokenIn:           fromToken.Address,
		TokenOut:          toToken.Address,
		Fee:               u.fee,
		Recipient:         filler,
		Deadline:          deadline,
		AmountIn:          amount,
		AmountOutMinimum:  big.NewInt(0),
		SqrtPriceLimitX96: big.NewInt(0),
	}

	calldata, err := u.uniswapRouterABI.Pack("exactInputSingle", inputStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function data: %w", err)
	}

	outputAmountWei := new(big.Int)
	outputAmountWei.SetString(result.AmountReceived.String(), 10)

	quote := &quoter.Quote{
		RequestId: requestId,
		Timestamp: time.Now().Unix(),
		Prices: map[string]string{
			amountIn: outputAmountWei.String(),
		},
		Protocol: u.String(),
		Path: []common.Address{
			fromToken.Address,
			toToken.Address,
		},
		Calldata: calldata,
		To:       u.routerAddress,
	}

	return quote, nil
}

func (u *UniswapV3) ExactOut(
	fromToken commons.Token,
	toToken commons.Token,
	amountOut string,
	requestId string,
	route interface{},
	filler common.Address,
) (*quoter.Quote, error) {
	return nil, nil
}

func (u *UniswapV3) String() string {
	return fmt.Sprintf("uniswap_v3_%d", u.fee.Int64())
}
