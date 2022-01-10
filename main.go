package main

import (
	"context"
	"errors"
	"flag"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pravinba9495/kryptonite/api/oneinch"
	"github.com/pravinba9495/kryptonite/chains"
	"github.com/pravinba9495/kryptonite/networks"
)

func main() {

	address := flag.String("wallet", "0x0000000000000000000000000000000000000000", "Your wallet address")
	chainID := flag.Uint64("chain", 1, "Your preferred chain")
	stableCoin := flag.String("stable", "USDC", "Your preferred stable coin")
	targetCoin := flag.String("target", "ETH", "Your high conviction coin")
	router := flag.String("router", "1INCH", "Your preferred router")

	flag.Parse()

	chain, err := chains.GetChainNameByID(*chainID)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Your wallet address: " + *address)
	log.Println("Your preferred chain: " + chain)
	log.Println("Your preferred stable coin/token: " + *stableCoin)
	log.Println("Your high conviction target coin/token: " + *targetCoin)
	log.Println("Your preferred router: " + *router)

	rpc, err := networks.GetRpcURLByChainID(*chainID)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = ethclient.DialContext(context.TODO(), string(rpc))
	if err != nil {
		log.Fatalln(err)
	}

	routerAddress, err := oneinch.GetRouterAddressByChainID(*chainID)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Router Contract Address: " + routerAddress + " (Verify at https://blockscan.com/address/" + routerAddress + ")")

	if err := oneinch.GetRouterHealthStatus(*chainID); err != nil {
		log.Fatalln(err)
	}

	_, err = oneinch.GetRouterLiquiditySources(*chainID)
	if err != nil {
		log.Fatalln(err)
	}

	tokens, err := oneinch.GetRouterSupportedTokens(*chainID)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = oneinch.GetRouterPresets(*chainID)
	if err != nil {
		log.Fatalln(err)
	}

	stableCoinContractAddress := ""
	targetCoinContractAddress := ""

	for _, token := range tokens {
		if token.Symbol == *stableCoin {
			stableCoinContractAddress = token.Address
			log.Println(*stableCoin + " Token Address: " + stableCoinContractAddress + " (Verify at https://blockscan.com/address/" + stableCoinContractAddress + ")")
		}
		if token.Symbol == *targetCoin {
			targetCoinContractAddress = token.Address
			log.Println(*targetCoin + " Token Address: " + targetCoinContractAddress + " (Verify at https://blockscan.com/address/" + targetCoinContractAddress + ")")
		}
	}

	if stableCoinContractAddress == "" || targetCoinContractAddress == "" {
		log.Fatalln(errors.New("token pair not possible in this chain"))
	}

	tokenMap, err := oneinch.GetTokenBalances(*chainID, *address)
	if err != nil {
		log.Fatalln(err)
	}

	balances := "Current Token Balances:"
	for tokenAddress, balance := range tokenMap {
		for _, token := range tokens {
			if (token.Address == string(tokenAddress)) && (token.Symbol == *stableCoin || token.Symbol == *targetCoin) {
				balances += ("	" + token.Symbol + " => " + string(balance))
			}
		}
	}
	log.Println(balances)

	params := &oneinch.ApproveAllowanceParamsDto{
		TokenAddress:  stableCoinContractAddress,
		WalletAddress: *address,
	}
	routerAllowance, err := oneinch.GetRouterAllowance(*chainID, params)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(routerAllowance)

	txParams := &oneinch.ApproveCalldataParamsDto{
		TokenAddress: stableCoinContractAddress,
		Amount:       "0",
	}
	txData, err := oneinch.GetRouterTransactionData(*chainID, txParams)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(txData.Data)

	quoteParams := &oneinch.QuoteParamsDto{
		FromTokenAddress: stableCoinContractAddress,
		ToTokenAddress:   targetCoinContractAddress,
		Amount:           "51030000",
	}
	quote, err := oneinch.GetRouterQuote(*chainID, quoteParams)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(quote)
}
