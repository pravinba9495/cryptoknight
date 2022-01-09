package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pravinba9495/kryptonite/networks"
)

func main() {

	address := flag.String("address", "0x0000000000000000000000000000000000000000", "Your wallet address")
	chainID := flag.Int("chain", 1, "Your preferred chain")
	stableCoin := flag.String("stable", "USDC", "Your preferred stable coin")
	targetCoin := flag.String("target", "ETH", "Your high conviction coin")

	flag.Parse()

	log.Println("Your address: " + *address)
	log.Println("Your preferred chain: " + fmt.Sprint(*chainID))
	log.Println("Your preferred stable coin/token: " + *stableCoin)
	log.Println("Your high conviction coin/token: " + *targetCoin)

	rpc, err := networks.GetRpcURLByChainID(*chainID)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := ethclient.DialContext(context.TODO(), string(rpc))
	if err != nil {
		log.Fatalln(err)
	}

	balance, err := client.BalanceAt(context.TODO(), common.HexToAddress(*address), nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(balance)
}
