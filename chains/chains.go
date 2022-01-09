package chains

import "errors"

// GetChainNameByID returns the chain name for the given chain id
func GetChainNameByID(chainID int) (string, error) {
	switch chainID {
	case 1:
		return "Ethereum", nil
	case 56:
		return "Binance Smart Chain", nil
	case 137:
		return "Polygon", nil
	case 10:
		return "Optimisim", nil
	case 42161:
		return "Arbitrum", nil
	default:
		return "", errors.New("unknown chain id provided")
	}
}
