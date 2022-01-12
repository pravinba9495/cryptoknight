package coingecko

type Coin struct {
	ID     string `json:"id,omitempty"`
	Symbol string `json:"symbol,omitempty"`
	Name   string `json:"name,omitempty"`
}

type ChartPoint []float64

type MarketChartResponseDto struct {
	Prices []ChartPoint `json:"prices,omitempty"`
}
