package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type Router struct {
	Vendor          string
	Address         *common.Address
	ChainID         uint64
	SupportedTokens []Token
}
