package oneinch

// RouterEndpoint represents a 1inch API endpoint
type RouterEndpoint string

const (
	// Approve endpoints
	SpenderEndpoint     RouterEndpoint = "/approve/spender"
	TransactionEndpoint RouterEndpoint = "/approve/transaction"
	AllowanceEndpoint   RouterEndpoint = "/approve/allowance"
)

const (
	// Healthcheck endpoints
	HealthcheckEndpoint RouterEndpoint = "/healthcheck"
)

const (
	// Info endpoints
	LiquiditySourcesEndpoint RouterEndpoint = "/liquidity-sources"
	TokensEndpoint           RouterEndpoint = "/tokens"
	PresetsEndpoint          RouterEndpoint = "/presets"
)

const (
	// Swap endpoints
	QuoteEndpoint RouterEndpoint = "/quote"
	SwapEndpoint  RouterEndpoint = "/swap"
)
