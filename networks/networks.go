package networks

// RpcURL represents the public RPC endpoint for the network
type RpcURL string

const (
	// For Ethereum
	Ethereum RpcURL = "https://cloudflare-eth.com/"

	// For Binance Smart Chain (BSC)
	BinanceSmartChain RpcURL = "https://bsc-dataseed.binance.org/"

	// For Polygon
	Polygon RpcURL = "https://polygon-rpc.com/"

	// For Optimism
	Optimisim RpcURL = "https://mainnet.optimism.io/"

	// For Arbitrum
	Arbitrum RpcURL = "https://arb1.arbitrum.io/rpc"
)
