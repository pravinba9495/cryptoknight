package technical

func IsABuy(currentTokenPrice float64, buyLimit float64) bool {
	return buyLimit > currentTokenPrice
}

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
