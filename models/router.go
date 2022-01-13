package models

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Router struct {
	Vendor          string
	Address         *common.Address
	ChainID         uint64
	SupportedTokens []Token
}

type ISwapRouter interface {
	GetContractAddress() (*common.Address, error)
	GetSupportedTokens() ([]Token, error)

	GetHealthStatus() error
	GetApprovedAllowance(tokenContractAddress *common.Address) (*big.Int, error)
	GenerateApprovalData() ([]byte, error)
	GenerateSwapQuotes()
	GenerateSwapData() ([]byte, error)
}
