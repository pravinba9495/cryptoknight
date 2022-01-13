package coingecko

type Coin struct {
	ID     string `json:"id,omitempty"`
	Symbol string `json:"symbol,omitempty"`
	Name   string `json:"name,omitempty"`
}
