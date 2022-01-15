package helpers

import (
	"github.com/pravinba9495/kryptonite/models"
)

// GetTokenAddress retrieves the token's contract address in the router's network
func GetTokenAddress(router *models.Router, token string) string {
	var address string
	for _, routerToken := range router.SupportedTokens {
		if routerToken.Symbol == token {
			address = routerToken.Address.Hex()
		}
	}
	return address
}
