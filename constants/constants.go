package constants

const (
	ApiBaseUrl string = "https://api.1inch.io"
	ApiVersion string = "v4.0"

	// Approve endpoints
	SpenderEndpoint     string = "/approve/spender"
	TransactionEndpoint string = "/approve/transaction"
	AllowanceEndpoint   string = "/approve/allowance"

	// Healthcheck endpoints
	HealthcheckEndpoint string = "/healthcheck"

	// Info endpoints
	TokensEndpoint string = "/tokens"

	// Swap endpoints
	QuoteEndpoint string = "/quote"
	SwapEndpoint  string = "/swap"
)
