package models

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"os"
	"os/exec"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pravinba9495/kryptonite/networks"
	"github.com/pravinba9495/kryptonite/store"
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

// ApproveSpender approves ERC20 access for 1inch router
func (w *Wallet) ApproveSpender(routerAddress *common.Address, tokenContractAddress *common.Address) (*types.Transaction, error) {
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
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	auth, _ := bind.NewKeyedTransactorWithChainID(w.PrivateKey, chainID)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = 11500000
	auth.GasPrice = gasPrice

	instance, err := store.NewStore(*tokenContractAddress, client)
	if err != nil {
		return nil, err
	}

	signedTx, err := instance.Approve(auth, *routerAddress, w.StableCoinBalance)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}
