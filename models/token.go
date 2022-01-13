package models

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Token struct {
	Name     string
	Symbol   string
	Decimals uint64
	Address  *common.Address
}

type TokenAddressWithBalance struct {
	Address *common.Address
	Balance *big.Int
}
