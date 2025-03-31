package amm

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"route_master/aggregator/quoter"
	"route_master/commons"
	"route_master/config"
	"route_master/contracts/uniswap_v3"
	"route_master/networking"
)

type SushiSwapV3 struct {
	fee                *big.Int
	client             *networking.SolidEthClient
	quoterAddress      common.Address
	sushiswapRouterABI *abi.ABI
	routerAddress      common.Address
}

func NewSushiSwapV3(client *networking.SolidEthClient, fee *big.Int) *SushiSwapV3 {
	quoterAddress := common.HexToAddress("0x64e8802FE490fa7cc61d3463958199161Bb608A7") // SushiSwap V3 Quoter address
	sushiswapRouterABI, _ := abi.JSON(strings.NewReader(config.RouterABIs.V3RouterABI))

	return &SushiSwapV3{
		fee:                fee,
		client:             client,
		quoterAddress:      quoterAddress,
		sushiswapRouterABI: &sushiswapRouterABI,
		routerAddress:      common.HexToAddress("0xd9e1cE17f2641f24aE83637ab66a2cca9C378B9F"), // SushiSwap V3 Router address
	}
}

func (u *SushiSwapV3) ExactIn(
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

	result, err := quoterContract.QuoteExactInputSingle(&bind.CallOpts{}, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote: %w", err)
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

	calldata, err := u.sushiswapRouterABI.Pack("exactInputSingle", inputStruct)
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
		Routes: map[string]interface{}{
			"expiry": deadline.Int64(),
			"gas":    result.GasEstimate.Uint64(),
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

func (u *SushiSwapV3) ExactOut(
	fromToken commons.Token,
	toToken commons.Token,
	amountOut string,
	requestId string,
	route interface{},
	filler common.Address,
) (*quoter.Quote, error) {
	return nil, nil
}

func (u *SushiSwapV3) String() string {
	return fmt.Sprintf("sushiswap_v3_%d", u.fee.Int64())
}
