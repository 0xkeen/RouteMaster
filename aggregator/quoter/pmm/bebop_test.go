package pmm

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	common2 "route_master/commons"
)

func TestBebopQuoter_ExactIn(t *testing.T) {
	quoter := NewBebopQuoter()

	fromToken := common2.Token{
		Address:  common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), // USDC
		Decimals: 6,
		Symbol:   "USDC",
		Name:     "USD Coin",
	}

	toToken := common2.Token{
		Address:  common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"), // WETH
		Decimals: 18,
		Symbol:   "WETH",
		Name:     "Wrapped Ether",
	}

	quote, err := quoter.ExactIn(
		fromToken,
		toToken,
		"10000000",
		"test-request-id",
		nil,
		common.HexToAddress("0xf38Ee8620305D064C42713A02348657d143D4995"),
	)
	j, _ := json.Marshal(quote)
	fmt.Println(string(j))

	assert.Nil(t, err)
}
