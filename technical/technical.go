package technical

import (
	"math"

	"github.com/pravinba9495/kryptonite/variables"
)

// IsABuy checks if the currentTokenPrice is at a BUY
func IsABuy(currentTokenPrice float64, buyLimit float64, tokenAmountFromSwap float64, tokenAmountFromPrice float64) bool {
	slippage := math.Abs((tokenAmountFromPrice - tokenAmountFromSwap) * 100 / tokenAmountFromPrice)
	cond1 := float64(variables.Slippage) >= slippage
	cond2 := currentTokenPrice <= buyLimit
	return cond1 && cond2
}

// IsASell checks if the currentTokenPrice is at a SELL
func IsASell(currentTokenPrice float64, previousTokenPrice float64, sellLimit float64, stopLimit float64) (bool, string, float64) {
	isASell, typ, value := false, "", 0.0
	if currentTokenPrice > sellLimit {
		isASell = true
		typ = "Profit"
		value = ((currentTokenPrice - previousTokenPrice) * 100) / (previousTokenPrice)
	}
	if currentTokenPrice < stopLimit {
		isASell = true
		typ = "Stop Loss"
		value = ((currentTokenPrice - previousTokenPrice) * 100) / (previousTokenPrice)
	}
	return isASell, typ, value
}
