package coingecko

// Coin refers to the coingecko coin model
type Coin struct {
	ID     string `json:"id,omitempty"`
	Symbol string `json:"symbol,omitempty"`
	Name   string `json:"name,omitempty"`
}
