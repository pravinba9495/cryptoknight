package oneinch

// Endpoint represents a 1inch API endpoint
type Endpoint string

const (
	// Approve endpoints
	SpenderEndpoint     Endpoint = "/approve/spender"
	TransactionEndpoint Endpoint = "/approve/transaction"
	AllowanceEndpoint   Endpoint = "/approve/allowance"
)

const (
	// Healthcheck endpoints
	HealthcheckEndpoint Endpoint = "/healthcheck"
)

const (
	// Info endpoints
	TokensEndpoint Endpoint = "/tokens"
)

const (
	// Swap endpoints
	QuoteEndpoint Endpoint = "/quote"
	SwapEndpoint  Endpoint = "/swap"
)
