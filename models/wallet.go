package models

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"os"
	"os/exec"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pravinba9495/kryptonite/networks"
)

type Wallet struct {
	ChainID            uint64
	Address            *common.Address
	PrivateKey         *ecdsa.PrivateKey
	MainAccountBalance big.Int
	StableCoinBalance  big.Int
	TargetCoinBalance  big.Int
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
		StableCoinBalance:  *big.NewInt(0),
		TargetCoinBalance:  *big.NewInt(0),
		MainAccountBalance: *big.NewInt(0),
	}
	return w, nil
}

func (w *Wallet) RefreshWalletBalance() error {
	balance, err := w.GetMainAccountBalance()
	if err != nil {
		return err
	}
	w.MainAccountBalance = *balance
	return nil
}

func (w *Wallet) GetTokenBalances(stableTokenContractAddress string, targetTokenContractAddress string) error {
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

	n := new(big.Int)
	n, ok := n.SetString(balance, 10)
	if !ok {
		return errors.New("SetString: error")
	}

	w.TargetCoinBalance = *n

	_, err = exec.Command(path+"/node-eth", "--address="+w.Address.String(), "--token="+stableTokenContractAddress, "--provider="+rpc).Output()
	if err != nil {
		return err
	}
	out, err = exec.Command("cat", stableTokenContractAddress+".txt").Output()
	if err != nil {
		return err
	}
	balance = string(out)

	n, ok = n.SetString(balance, 10)
	if !ok {
		return errors.New("SetString: error")
	}

	w.StableCoinBalance = *n

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
