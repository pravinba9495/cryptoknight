package models

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pravinba9495/kryptonite/networks"
)

type TokenAddressWithBalance struct {
	Address *common.Address
	Balance *big.Int
}

type Wallet struct {
	ChainID            uint64
	Address            *common.Address
	PrivateKey         *ecdsa.PrivateKey
	MainAccountBalance *big.Int
	TokensWithBalance  []TokenAddressWithBalance
}

func (w *Wallet) New(address string, key string, chainID uint64) (*Wallet, error) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return nil, err
	}
	// publicKey := privateKey.Public()
	// publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	// if !ok {
	// 	return nil, errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	// }
	// addrPublic := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	// if addrPublic != address {
	// 	return nil, errors.New("public key and private key don't match with each other, please check")
	// }
	addr := common.HexToAddress(address)
	w = &Wallet{
		Address:            &addr,
		PrivateKey:         privateKey,
		ChainID:            chainID,
		TokensWithBalance:  nil,
		MainAccountBalance: big.NewInt(0),
	}

	tokensWithBalance, err := w.GetTokenBalances()
	if err != nil {
		return nil, err
	}
	w.TokensWithBalance = tokensWithBalance
	return w, nil
}

func (w *Wallet) RefreshWalletBalance() error {
	balance, err := w.GetMainAccountBalance()
	if err != nil {
		return err
	}
	w.MainAccountBalance = balance
	return nil
}

func (w *Wallet) RefreshTokenBalances() error {
	balances, err := w.GetTokenBalances()
	if err != nil {
		return err
	}
	w.TokensWithBalance = balances
	return nil
}

func (w *Wallet) GetMainAccountBalance() (*big.Int, error) {
	rpc, err := networks.GetRpcURLByChainID(w.ChainID)
	if err != nil {
		return nil, err
	}
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	balance, err := client.BalanceAt(context.Background(), *w.Address, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (w *Wallet) GetTokenBalances() ([]TokenAddressWithBalance, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://balances.1inch.io/v1.1/"+fmt.Sprint(w.ChainID)+"/balances/"+w.Address.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var dto map[string]string
		if err := json.Unmarshal(body, &dto); err != nil {
			return nil, err
		}
		tokensWithBalance := make([]TokenAddressWithBalance, 0)
		for token, balance := range dto {
			tokenAddr := common.HexToAddress(token)
			tokenBalance := big.NewInt(0)
			n, ok := tokenBalance.SetString(balance, 10)
			if !ok {
				return nil, errors.New("cannot convert string to *big.Int")
			}
			tokensWithBalance = append(tokensWithBalance, TokenAddressWithBalance{
				Address: &tokenAddr,
				Balance: n,
			})
		}
		return tokensWithBalance, nil
	} else {
		return nil, errors.New(resp.Status)
	}
}
