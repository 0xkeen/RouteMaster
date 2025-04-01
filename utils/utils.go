package utils

import (
	"github.com/inconshreveable/log15"
	"github.com/shopspring/decimal"
	"math/big"
	"regexp"
)

var logger = log15.New("module", "Util")

func IsValidAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	return re.MatchString(address)
}

func IsValidAddressList(addresses []string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	for _, addr := range addresses {
		if !re.MatchString(addr) {
			return false
		}
	}

	return true
}

func ToDecimal(iValue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := iValue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	default:
		logger.Info("type is not supported", "type", v)
		return decimal.Decimal{}
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

func ToWei(iAmount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iAmount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	default:
		logger.Info("type is not supported", "type", v)
		return &big.Int{}
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}
