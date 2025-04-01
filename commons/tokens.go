package commons

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/patrickmn/go-cache"

	"route_master/contracts/erc20"
)

var ETHAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
var WETHAddress = common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")

type Tokens struct {
	tokens  *cache.Cache
	backend *ethclient.Client
}

func NewTokens(backend *ethclient.Client) *Tokens {
	return &Tokens{
		tokens:  cache.New(cache.NoExpiration, 0),
		backend: backend,
	}
}

func (t *Tokens) GetToken(tokenAddress common.Address) Token {
	if token, ok := t.tokens.Get(tokenAddress.String()); ok {
		return token.(Token)
	}

	inputTokenInfo, _ := erc20.NewToken(tokenAddress, t.backend)
	decimals0, _ := inputTokenInfo.Decimals(nil)

	token := Token{
		Address:  tokenAddress,
		Decimals: decimals0,
	}

	t.tokens.Set(tokenAddress.String(), token, cache.NoExpiration)

	return token
}

func (t *Tokens) GetERC20Token(tokenAddress common.Address) Token {
	if tokenAddress == ETHAddress {
		tokenAddress = WETHAddress
	}
	return t.GetToken(tokenAddress)
}
