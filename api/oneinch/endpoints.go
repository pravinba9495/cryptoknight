package oneinch

type RouterEndpoint string

const (
	// Approve endpoints
	SpenderEndpoint     RouterEndpoint = "/approve/spender"
	TransactionEndpoint RouterEndpoint = "/approve/transaction"
	AllowanceEndpoint   RouterEndpoint = "/approve/allowance"
)
