package amm

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"route_master/aggregator/quoter"
	"route_master/commons"
	"route_master/contracts/pancakeswap_v3"
)

type PancakeSwapV3 struct {
	client        *ethclient.Client
	quoter        *pancakeswap_v3.Quoter
	router        *pancakeswap_v3.Router
	quoterAddress common.Address
	routerAddress common.Address
	factory       common.Address
	fee           *big.Int
}

func NewPancakeSwapV3(client *ethclient.Client, fee *big.Int) *PancakeSwapV3 {
	quoterAddress := common.HexToAddress("0xB048Bbc1Ee6b733FFfCFb9e9CeF7375518e25997")
	routerAddress := common.HexToAddress("0x1b81D678ffb9C0263b24A97847620C99d213eB14")

	quoter, _ := pancakeswap_v3.NewQuoter(quoterAddress, client)
	router, _ := pancakeswap_v3.NewRouter(routerAddress, client)

	return &PancakeSwapV3{
		client:        client,
		quoter:        quoter,
		router:        router,
		quoterAddress: quoterAddress,
		routerAddress: routerAddress,
		fee:           fee,
	}
}

func (p *PancakeSwapV3) String() string {
	return "pancakeSwapV3"
}

func (p *PancakeSwapV3) ExactIn(
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

	amountOut, err := p.querySwapExactIn(fromToken.Address, toToken.Address, amount, p.fee)
	if err != nil {
		return nil, fmt.Errorf("failed to query swap: %w", err)
	}

	if amountOut.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("insufficient liquidity for swap")
	}

	calldata, err := p.buildSwapCalldata(fromToken.Address, toToken.Address, amount, amountOut, p.fee, filler)
	if err != nil {
		return nil, fmt.Errorf("failed to build swap calldata: %w", err)
	}

	routes := map[string]interface{}{
		"fee":       p.fee.String(),
		"tokenIn":   fromToken.Address.String(),
		"tokenOut":  toToken.Address.String(),
		"amountIn":  amount.String(),
		"amountOut": amountOut.String(),
	}

	timestamp := time.Now().Unix()
	return &quoter.Quote{
		RequestId: requestId,
		To:        p.routerAddress,
		AmountIn:  amount,
		AmountOut: amountOut,
		Prices:    map[string]string{fromToken.Symbol: toToken.Symbol},
		Routes:    routes,
		Protocol:  p.String(),
		Path:      []common.Address{fromToken.Address, toToken.Address},
		Calldata:  calldata,
		Timestamp: timestamp,
	}, nil
}

func (p *PancakeSwapV3) ExactOut(
	fromToken commons.Token,
	toToken commons.Token,
	amountOut string,
	requestId string,
	route interface{},
	filler common.Address,
) (*quoter.Quote, error) {
	return nil, nil
}

func (p *PancakeSwapV3) querySwapExactIn(fromAddress, toAddress common.Address, amountIn, fee *big.Int) (*big.Int, error) {
	params := pancakeswap_v3.IQuoterV2QuoteExactInputSingleParams{
		TokenIn:           fromAddress,
		TokenOut:          toAddress,
		AmountIn:          amountIn,
		Fee:               fee,
		SqrtPriceLimitX96: big.NewInt(0),
	}

	quoterAbi, err := abi.JSON(strings.NewReader(pancakeswap_v3.QuoterMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse quoter ABI: %w", err)
	}

	data, err := quoterAbi.Pack("quoteExactInputSingle", params)
	if err != nil {
		return nil, fmt.Errorf("failed to pack quoteExactInputSingle data: %w", err)
	}

	result, err := p.client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &p.quoterAddress,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call quoteExactInputSingle: %w", err)
	}

	var outputs struct {
		AmountOut               *big.Int
		SqrtPriceX96After       *big.Int
		InitializedTicksCrossed uint32
		GasEstimate             *big.Int
	}
	err = quoterAbi.UnpackIntoInterface(&outputs, "quoteExactInputSingle", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack quoteExactInputSingle result: %w", err)
	}

	return outputs.AmountOut, nil
}

func (p *PancakeSwapV3) buildSwapCalldata(fromAddress, toAddress common.Address, amountIn, amountOut, fee *big.Int, filler common.Address) ([]byte, error) {
	deadline := big.NewInt(time.Now().Unix() + 1800)

	minAmountOut := new(big.Int).Mul(amountOut, big.NewInt(995))
	minAmountOut = new(big.Int).Div(minAmountOut, big.NewInt(1000))

	params := pancakeswap_v3.ISwapRouterExactInputSingleParams{
		TokenIn:           fromAddress,
		TokenOut:          toAddress,
		Fee:               fee,
		Recipient:         filler,
		Deadline:          deadline,
		AmountIn:          amountIn,
		AmountOutMinimum:  minAmountOut,
		SqrtPriceLimitX96: big.NewInt(0),
	}

	routerAbi, err := abi.JSON(strings.NewReader(pancakeswap_v3.RouterMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse router ABI: %w", err)
	}

	calldata, err := routerAbi.Pack("exactInputSingle", params)
	if err != nil {
		return nil, fmt.Errorf("failed to pack exactInputSingle data: %w", err)
	}

	return calldata, nil
}
