package amm

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	
	"route_master/commons"
	"route_master/networking"
	"route_master/utils"
)

func TestQuoteSushiSwapV3(t *testing.T) {
	privateNodeClient, _ := networking.Dial(os.Getenv("PRIVATE_NODE_URL"))
	solidClient := networking.NewSolidEthClient(privateNodeClient, privateNodeClient)
	sushiV3 := NewSushiSwapV3(solidClient, big.NewInt(3000))
	quote, err := sushiV3.ExactIn(
		commons.Token{
			Address:  common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
			Decimals: 18,
			Symbol:   "WETH",
			Name:     "WETH",
		},
		commons.Token{
			Address:  common.HexToAddress("0x6b175474e89094c44da98b954eedeac495271d0f"),
			Decimals: 18,
			Symbol:   "DAI",
			Name:     "DAI",
		},
		utils.ToWei("0.5", 18).String(),
		"",
		nil,
		common.HexToAddress("0xC871d5AEafe08d43eFAe3B81351B42740402f442"),
	)
	assert.Nil(t, err)
	fmt.Println(hexutil.Encode(quote.Calldata))
}
