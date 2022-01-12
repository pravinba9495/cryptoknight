package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/pravinba9495/kryptonite/coingecko"
	"github.com/pravinba9495/kryptonite/models"
	"github.com/pravinba9495/kryptonite/oneinch"
)

func main() {

	var address, privateKey, aggregator, stableToken, targetToken string
	var chainID, days, profit uint64
	flag.StringVar(&address, "public", "0x0000000000000000000000000000000000000000", "Your wallet public address")
	flag.StringVar(&privateKey, "private", "", "Your wallet private key")
	flag.StringVar(&stableToken, "stable", "USDC", "Stable token (ERC20) to use. Example: USDC, USDT, DAI")
	flag.StringVar(&targetToken, "target", "WETH", "Target ERC20 token to hold. Example: WETH, WMATIC, LINK.")
	flag.StringVar(&aggregator, "aggregator", "1INCH", "Aggregator to use. Allowed options: 1INCH, PARASWAP")
	flag.Uint64Var(&chainID, "chain", 1, "Chain to use. Allowed options: 1 (Ethereum), 10 (Optimism), 56 (Binance Smart Chain), 137 (Polygon/Matic), 42161 (Arbitrum)")
	flag.Uint64Var(&days, "days", 30, "No. of days to use to calculate moving average")
	flag.Uint64Var(&profit, "profit", 15, "No. of days to use to calculate moving average, support and resistance")
	flag.Parse()

	wallet, err := (&models.Wallet{}).New(address, privateKey, chainID)
	if err != nil {
		log.Fatalln(err)
	}

	var router *models.Router
	if aggregator == "1INCH" {
		r, err := oneinch.New(chainID)
		if err != nil {
			log.Fatalln(err)
		}
		router = r
	}

	var stableTokenContractAddress string
	var targetTokenContractAddress string

	for _, walletToken := range wallet.TokensWithBalance {
		for _, routerToken := range router.SupportedTokens {
			if routerToken.Symbol == targetToken && walletToken.Address.Hex() == routerToken.Address.Hex() {
				targetTokenContractAddress = routerToken.Address.Hex()
			}
			if routerToken.Symbol == stableToken && walletToken.Address.Hex() == routerToken.Address.Hex() {
				stableTokenContractAddress = routerToken.Address.Hex()
			}
		}
	}

	if stableTokenContractAddress == "" || targetTokenContractAddress == "" {
		log.Fatalln("Token pair not available to swap in this chain/router")
	}

	// rpc, err := networks.GetRpcURLByChainID(wallet.ChainID)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// client, err := ethclient.Dial(rpc)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer client.Close()

	// nonce, err := client.PendingNonceAt(context.Background(), *wallet.Address)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// tx := &types.LegacyTx{
	// 	Nonce:    nonce,
	// 	Gas:      21000,
	// 	GasPrice: big.NewInt(30000000000),
	// 	To:       router.Address,
	// 	Value:    big.NewInt(0),
	// 	Data:     []byte(""),
	// }

	// _, err = types.SignNewTx(wallet.PrivateKey, types.LatestSignerForChainID(big.NewInt(int64(wallet.ChainID))), tx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := client.SendTransaction(context.TODO(), signedTx); err != nil {
	// 	log.Fatalln(err)
	// }
	// log.Printf("Transaction sent: %s", signedTx.Hash().Hex())

	for {
		if err := wallet.RefreshWalletBalance(); err != nil {
			log.Println(err)
		} else {
			log.Printf("Wallet Main Balance: %s (wei)", wallet.MainAccountBalance)
			if err := wallet.RefreshTokenBalances(); err != nil {
				log.Println(err)
			} else {
				for _, tokenWithBalance := range wallet.TokensWithBalance {
					if tokenWithBalance.Balance.Uint64() > 0 || tokenWithBalance.Address.Hex() == stableTokenContractAddress || tokenWithBalance.Address.Hex() == targetTokenContractAddress {
						symbol := ""
						for _, t := range router.SupportedTokens {
							if tokenWithBalance.Address.Hex() == t.Address.Hex() {
								symbol = t.Symbol
								break
							}
						}
						log.Printf("%s [%s] => %s (wei)", symbol, tokenWithBalance.Address, tokenWithBalance.Balance)
					}
				}
				coins, err := coingecko.GetCoinsList()
				if err != nil {
					log.Println(err)
				} else {
					targetCoinID := ""
					for _, coin := range coins {
						if coin.Symbol == strings.ToLower(targetToken) {
							targetCoinID = coin.ID
							break
						}
					}
					if targetCoinID == "" {
						log.Println("coin not found")
					} else {

						currentTokenPrice, err := coingecko.GetCoinPrice(targetCoinID)
						if err != nil {
							log.Println(err)
						} else {
							data, err := coingecko.GetMarketChartByCoin(targetCoinID, days)
							if err != nil {
								log.Println(err)
							} else {

								var ma, sum, recentLow, recentHigh float64

								prices := make([]float64, 0)
								for i, point := range data.Prices {
									if i < len(data.Prices)-1 {
										if point[1] > 0 {
											sum += point[1]
											prices = append(prices, point[1])
										}
									}
								}

								days = uint64(len(prices))

								if days > 0 {
									recentLow = prices[0]
									recentHigh = prices[0]
									for _, price := range prices {
										if recentHigh < price {
											recentHigh = price
										}
										if recentLow > price {
											recentLow = price
										}
									}

									ma = sum / float64(days)
									upside := ((recentHigh - currentTokenPrice) * 100) / currentTokenPrice
									downside := ((currentTokenPrice - recentLow) * 100) / currentTokenPrice
									log.Printf("Current Price: %f $", currentTokenPrice)
									log.Printf("Recent Low (Last %d days): %f $", days, recentLow)
									log.Printf("Recent High (Last %d days): %f $", days, recentHigh)
									log.Printf("Average Price (Last %d days): %f $", days, ma)
									log.Printf("Upside Potential: +%.2f%s", upside, "%")
									log.Printf("Downside Possible: -%.2f%s", downside, "%")

								} else {
									log.Println("malformed prices data received")
								}

							}
						}
					}
				}
			}
		}
		time.Sleep(1 * time.Minute)
		log.Printf("")
		log.Printf("")
		log.Printf("")
	}
}
