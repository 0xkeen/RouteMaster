package types

import "github.com/ethereum/go-ethereum/common"

var (
	WETHAddress = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	ZeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
	ETHAddress  = common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE")
)

func IsETHAddress(addr common.Address) bool {
	return addr == ZeroAddress || addr == ETHAddress
}
