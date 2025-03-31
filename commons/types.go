package commons

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Token struct {
	Address  common.Address
	Decimals uint8
	Symbol   string
	Name     string
}

type PairInfo struct {
	PairAddress common.Address
	Token0      common.Address
	ReserveIn   *big.Int
	ReserveOut  *big.Int
	IsToken0In  bool
}

type SegmentDetail struct {
	Protocol       string
	AmountIn       string
	AmountOut      string
	TotalAmountIn  string
	TotalAmountOut string
	Calldata       []byte
	Routes         map[string]interface{}
	To             common.Address
}

type Settlement struct {
	Settler common.Address
	Value   string
	Data    []byte
}

type Route struct {
	Protocol       string
	TotalAmountIn  string
	TotalAmountOut string
}
