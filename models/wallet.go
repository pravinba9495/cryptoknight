package models

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"os"
	"os/exec"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pravinba9495/kryptonite/networks"
)

// Wallet refers to the crypto wallet
type Wallet struct {
	ChainID            uint64
	Address            *common.Address
	PrivateKey         *ecdsa.PrivateKey
	MainAccountBalance *big.Int
	StableCoinBalance  *big.Int
	TargetCoinBalance  *big.Int
}

// New creates a new wallet instance
func (w *Wallet) New(address string, privateKey string, chainID uint64) (*Wallet, error) {
	pKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}
	publicKey := pKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	addrPublic := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	if addrPublic != address {
		return nil, errors.New("public key and private key don't match with each other")
	}
	addr := common.HexToAddress(addrPublic)
	w = &Wallet{
		Address:            &addr,
		PrivateKey:         pKey,
		ChainID:            chainID,
		StableCoinBalance:  big.NewInt(0),
		TargetCoinBalance:  big.NewInt(0),
		MainAccountBalance: big.NewInt(0),
	}
	return w, nil
}

// RefreshWalletBalance refreshes the wallet balance
func (w *Wallet) RefreshWalletBalance() error {
	balance, err := w.GetMainAccountBalance()
	if err != nil {
		return err
	}
	w.MainAccountBalance = balance
	return nil
}

// RefreshTokenBalance refreshes the specified ERC20 token balance
func (w *Wallet) RefreshTokenBalances(stableTokenContractAddress string, targetTokenContractAddress string) error {
	rpc, err := networks.GetRpcURLByChainID(w.ChainID)
	if err != nil {
		return err
	}
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return err
	}
	defer client.Close()

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	_, err = exec.Command(path+"/node-eth", "--address="+w.Address.String(), "--token="+targetTokenContractAddress, "--provider="+rpc).Output()
	if err != nil {
		return err
	}
	out, err := exec.Command("cat", targetTokenContractAddress+".txt").Output()
	if err != nil {
		return err
	}
	balance := string(out)

	n1 := new(big.Int)
	n1, ok := n1.SetString(balance, 10)
	if !ok {
		return errors.New("SetString: error")
	}

	w.TargetCoinBalance = n1

	_, err = exec.Command(path+"/node-eth", "--address="+w.Address.String(), "--token="+stableTokenContractAddress, "--provider="+rpc).Output()
	if err != nil {
		return err
	}
	out, err = exec.Command("cat", stableTokenContractAddress+".txt").Output()
	if err != nil {
		return err
	}
	balance = string(out)

	n2 := new(big.Int)
	n2, ok = n2.SetString(balance, 10)
	if !ok {
		return errors.New("SetString: error")
	}

	w.StableCoinBalance = n2

	return nil
}

// GetMainAccountBalance retrieves the main account balance
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

// SendTransaction creates and broadcasts a transaction to the blockchain after signing the payload
func (w *Wallet) SendTransaction(toAddress *common.Address, tx *types.LegacyTx) (*types.Transaction, error) {
	rpc, err := networks.GetRpcURLByChainID(w.ChainID)
	if err != nil {
		return nil, err
	}
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	nonce, err := client.PendingNonceAt(context.Background(), *w.Address)
	if err != nil {
		return nil, err
	}
	if tx.Gas == 0 {
		g, err := client.EstimateGas(context.TODO(), ethereum.CallMsg{
			Data: tx.Data,
		})
		if err != nil {
			return nil, err
		}
		tx.Gas = g
	}
	t := &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: tx.GasPrice,
		Gas:      uint64(1.25 * float64(tx.Gas)),
		To:       tx.To,
		Value:    big.NewInt(0),
		Data:     tx.Data,
	}
	signedTx, err := types.SignNewTx(w.PrivateKey, types.LatestSignerForChainID(big.NewInt(int64(w.ChainID))), t)
	if err != nil {
		return nil, err
	}
	if err := client.SendTransaction(context.TODO(), signedTx); err != nil {
		return nil, err
	}
	return signedTx, nil
}
