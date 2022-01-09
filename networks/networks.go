package networks

import "errors"

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

// GetRpcURLByChainID returns the RPC url for the given chain id
func GetRpcURLByChainID(chainID int) (RpcURL, error) {
	switch chainID {
	case 1:
		return Ethereum, nil
	case 56:
		return BinanceSmartChain, nil
	case 137:
		return Polygon, nil
	case 10:
		return Optimisim, nil
	case 42161:
		return Arbitrum, nil
	default:
		return "", errors.New("unknown chain name provided")
	}
}
