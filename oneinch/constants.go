package oneinch

const (
	ApiBaseUrl string = "https://api.1inch.io"
	ApiVersion string = "v4.0"
)

const (
	// Approve endpoints
	SpenderEndpoint     string = "/approve/spender"
	TransactionEndpoint string = "/approve/transaction"
	AllowanceEndpoint   string = "/approve/allowance"
)

const (
	// Healthcheck endpoints
	HealthcheckEndpoint string = "/healthcheck"
)

const (
	// Info endpoints
	TokensEndpoint string = "/tokens"
)

const (
	// Swap endpoints
	QuoteEndpoint string = "/quote"
	SwapEndpoint  string = "/swap"
)
