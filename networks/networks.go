package networks

import "errors"

const (
	// For Ethereum
	Ethereum string = "https://cloudflare-eth.com"

	// For Goerli Testnet
	Goerli string = "https://goerli.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"

	// For Binance Smart Chain (BSC)
	BinanceSmartChain string = "https://bsc-dataseed.binance.org"

	// For Polygon
	Polygon string = "https://polygon-rpc.com"

	// For Optimism
	Optimisim string = "https://mainnet.optimism.io"

	// For Arbitrum
	Arbitrum string = "https://arb1.arbitrum.io/rpc"
)

// GetRpcURLByChainID returns the RPC url for the given chain id
func GetRpcURLByChainID(chainID uint64) (string, error) {
	switch chainID {
	case 1:
		return Ethereum, nil
	case 5:
		return Goerli, nil
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
