package models

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Token represents the model of a token
type Token struct {
	Name     string
	Symbol   string
	Decimals uint64
	Address  *common.Address
}

// TokenAddressWithBalance represents the model of a token address with its balance
type TokenAddressWithBalance struct {
	Address *common.Address
	Balance *big.Int
}
