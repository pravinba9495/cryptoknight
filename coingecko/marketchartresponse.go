package coingecko

// MarketChartResponseDto refers to the coingecko market chart response
type MarketChartResponseDto struct {
	Prices []ChartPoint `json:"prices,omitempty"`
}
