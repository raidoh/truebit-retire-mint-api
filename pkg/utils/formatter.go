package utils

import (
	"fmt"
	"math/big"
)

var weiBase = new(big.Float).SetFloat64(1e18)

func WeiToEth(weiValue *big.Int) float64 {
	if weiValue == nil {
		return 0
	}
	ethFloat := new(big.Float).Quo(new(big.Float).SetInt(weiValue), weiBase)
	result, _ := ethFloat.Float64()
	return result
}

func FormatEthValue(weiValue *big.Int) string {
	if weiValue == nil {
		return "0"
	}
	return fmt.Sprintf("%.18f", WeiToEth(weiValue))
}

func FormatTokenSupply(supply *big.Int) string {
	if supply == nil {
		return "0"
	}
	supplyF64 := WeiToEth(supply)
	return FormatWithCommas(big.NewInt(int64(supplyF64)))
}

func FormatWithCommas(n *big.Int) string {
	if n == nil {
		return "0"
	}
	str := n.String()
	for i := len(str) - 3; i > 0; i -= 3 {
		str = str[:i] + "," + str[i:]
	}
	return str
}
