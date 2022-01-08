package chains

// ChainID represents the chain id of a blockchain
type ChainID int

const (
	// For Ethereum
	Ethereum ChainID = 1

	// For Binance Smart Chain (BSC)
	BinanceSmartChain ChainID = 56

	// For Polygon
	Polygon ChainID = 137

	// For Optimism
	Optimisim ChainID = 10

	// For Arbitrum
	Arbitrum ChainID = 42161
)
