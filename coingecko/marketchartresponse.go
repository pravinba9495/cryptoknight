package coingecko

type MarketChartResponseDto struct {
	Prices []ChartPoint `json:"prices,omitempty"`
}
