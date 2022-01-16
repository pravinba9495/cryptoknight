package technical

import (
	"errors"
	"math"
)

func GetMovingAverage(points []float64) float64 {
	sum := 0.00
	for _, price := range points {
		sum += price
	}
	movingAverage := sum / float64(len(points))
	return movingAverage
}

func IsASell(previousTokenPrice float64, currentTokenPrice float64, movingAverage float64, recentSupport float64, recentResistance float64, profitPercent int64, stopLoss int64) (bool, string, float64) {
	upside := ((recentResistance - currentTokenPrice) * 100) / currentTokenPrice
	currentProfitOrLossPercent := ((currentTokenPrice - previousTokenPrice) * 100) / previousTokenPrice

	if currentProfitOrLossPercent > 0 {
		// Profit Taking
		if float64(profitPercent) < currentProfitOrLossPercent && math.Abs(upside) < 2 {
			return true, "Profit", currentProfitOrLossPercent
		} else {
			// HODL
			return false, "Current Profit (Unrealized)", currentProfitOrLossPercent
		}
	} else {
		// Stop Loss
		if float64(stopLoss) < currentProfitOrLossPercent {
			return true, "Stop Loss", currentProfitOrLossPercent
		} else {
			// HODL
			return true, "Current Loss (Unrealized)", currentProfitOrLossPercent
		}
	}
}

func IsABuy(currentTokenPrice float64, movingAverage float64, recentSupport float64, recentResistance float64, profitPercent int64, stopLossPercent int64) (bool, float64, float64) {
	upside := ((recentResistance - currentTokenPrice) * 100) / currentTokenPrice
	downside := ((recentSupport - currentTokenPrice) * 100) / currentTokenPrice

	cond1 := currentTokenPrice > recentSupport
	cond2 := currentTokenPrice < recentResistance
	cond3 := currentTokenPrice < (recentSupport + (0.05 * recentSupport))
	cond4 := movingAverage > recentSupport
	cond5 := movingAverage < recentResistance
	cond6 := math.Abs(upside) > math.Abs(downside)
	cond7 := float64(stopLossPercent) > math.Abs(downside)

	return cond1 &&
		cond2 &&
		cond3 &&
		cond4 &&
		cond5 &&
		cond6 &&
		cond7, upside, downside
}

func CalculateResistanceLevels(points []float64, candlesBefore uint64, candlesAfter uint64) []float64 {
	resistances := make([]float64, 0)
	for index, point := range points {
		bool := IsResistance(points, index, candlesBefore, candlesAfter)
		if bool {
			resistances = append(resistances, point)
		}
	}
	return resistances
}

func CalculateSupportLevels(points []float64, candlesBefore uint64, candlesAfter uint64) []float64 {
	supports := make([]float64, 0)
	for index, point := range points {
		bool := IsSupport(points, index, candlesBefore, candlesAfter)
		if bool {
			supports = append(supports, point)
		}
	}
	return supports
}

func IsSupport(points []float64, index int, candlesBefore uint64, candlesAfter uint64) bool {
	var aF, bF uint64
	for i := 0; i < int(candlesBefore) && index-1-i >= 0; i++ {
		if points[index-1-i] > points[index] {
			bF += 1
		}
	}
	for i := 0; i < int(candlesAfter) && index+1+i < len(points); i++ {
		if points[index+1+i] > points[index] {
			aF += 1
		}
	}
	return aF >= candlesAfter && bF >= candlesBefore
}

func IsResistance(points []float64, index int, candlesBefore uint64, candlesAfter uint64) bool {
	var aF, bF uint64
	for i := 0; i < int(candlesBefore) && index-1-i >= 0; i++ {
		if points[index-1-i] < points[index] {
			bF += 1
		}
	}
	for i := 0; i < int(candlesAfter) && index+1+i < len(points); i++ {
		if points[index+1+i] < points[index] {
			aF += 1
		}
	}
	return aF >= candlesAfter && bF >= candlesBefore
}

func GetRecentSupportAndResistance(currentTokenPrice float64, supports []float64, resistances []float64) (float64, float64) {
	recentResistance := currentTokenPrice
	recentSupport := currentTokenPrice
	if len(resistances) > 0 {
		for index := len(resistances) - 1; index >= 0; index-- {
			if currentTokenPrice < resistances[index] {
				recentResistance = resistances[index]
				break
			}
		}
	}
	if len(supports) > 0 {
		for index := len(supports) - 1; index >= 0; index-- {
			if currentTokenPrice > supports[index] {
				recentSupport = supports[index]
				break
			}
		}
	}
	return recentSupport, recentResistance
}

func GetSupportsAndResistances(points []float64) ([]float64, []float64, error) {
	var resistances, supports []float64
	var candlesAfter, candlesBefore uint64 = 4, 3
	for {
		resistances = CalculateResistanceLevels(points, candlesBefore, candlesAfter)
		supports = CalculateSupportLevels(points, candlesBefore, candlesAfter)
		if len(resistances) > 0 && len(supports) > 0 {
			break
		} else {
			candlesAfter -= 1
			if candlesAfter < 2 {
				return nil, nil, errors.New("could not calculate pivot points")
			}
		}
	}
	return supports, resistances, nil
}
