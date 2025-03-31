package aggregator

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"route_master/aggregator/quoter"
	"route_master/aggregator/quoter/amm"
	"route_master/aggregator/quoter/pmm"
	"route_master/commons"
	"route_master/networking"
)

func TestAggQuote(t *testing.T) {
	privateNodeClient, _ := networking.Dial(os.Getenv("PRIVATE_NODE_URL"))
	solidClient := networking.NewSolidEthClient(privateNodeClient, privateNodeClient)

	bebopQuoter := pmm.NewBebopQuoter()
	quoters := []quoter.Quoter{bebopQuoter}
	quoters = []quoter.Quoter{}
	fees := []int64{100, 500, 1000, 3000, 10000}
	for _, fee := range fees {
		quoters = append(quoters, amm.NewUniswapV3(solidClient, big.NewInt(fee)))
		quoters = append(quoters, amm.NewSushiSwapV3(solidClient, big.NewInt(fee)))
	}

	agg := NewAggregator(quoters, 50)

	result := agg.GetQuote(
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
		big.NewInt(1000000000000000000),
		uuid.NewString(),
		common.Address{},
		100,
	)
	j, _ := json.Marshal(result)
	fmt.Println(string(j))
	assert.NotEmpty(t, string(j))
}
