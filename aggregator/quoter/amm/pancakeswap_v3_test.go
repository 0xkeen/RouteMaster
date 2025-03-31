package amm

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"testing"

	"route_master/commons"
	"route_master/utils"
)

func TestQuotePancakeSwapV3(t *testing.T) {
	client, _ := ethclient.Dial(os.Getenv("PRIVATE_NODE_URL"))
	uniV3 := NewPancakeSwapV3(client, big.NewInt(3000))
	quote, err := uniV3.ExactIn(
		commons.Token{
			Address:  common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
			Decimals: 18,
			Symbol:   "WETH",
			Name:     "WETH",
		},
		commons.Token{
			Address:  common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7"),
			Decimals: 6,
			Symbol:   "USDT",
			Name:     "USDT",
		},
		utils.ToWei("0.1", 18).String(),
		"",
		nil,
		common.HexToAddress("0xC871d5AEafe08d43eFAe3B81351B42740402f442"),
	)
	assert.Error(t, err)
	fmt.Println(quote)
}
