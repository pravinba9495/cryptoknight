package models

import "github.com/ethereum/go-ethereum/common"

type Token struct {
	Name     string
	Symbol   string
	Decimals uint64
	Address  *common.Address
}
