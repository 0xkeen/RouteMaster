package config

import "github.com/ethereum/go-ethereum/common"

var (
	WETH = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	ETH  = common.HexToAddress("0x0000000000000000000000000000000000000000")

	RedisKeyPrefixTokenInfo = "token_info"
)

var RouterABIs = struct {
	V3RouterABI string
}{
	V3RouterABI: `[{
		"inputs": [{
			"components": [
				{"name": "tokenIn", "type": "address"},
				{"name": "tokenOut", "type": "address"},
				{"name": "fee", "type": "uint24"},
				{"name": "recipient", "type": "address"},
				{"name": "deadline", "type": "uint256"},
				{"name": "amountIn", "type": "uint256"},
				{"name": "amountOutMinimum", "type": "uint256"},
				{"name": "sqrtPriceLimitX96", "type": "uint160"}
			],
			"name": "params",
			"type": "tuple"
		}],
		"name": "exactInputSingle",
		"outputs": [{"name": "amountOut", "type": "uint256"}],
		"stateMutability": "payable",
		"type": "function"
	}]`,
}
