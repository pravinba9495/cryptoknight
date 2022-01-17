package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pravinba9495/kryptonite/bot"
	"github.com/pravinba9495/kryptonite/coingecko"
	"github.com/pravinba9495/kryptonite/helpers"
	"github.com/pravinba9495/kryptonite/models"
	"github.com/pravinba9495/kryptonite/technical"
	"github.com/pravinba9495/kryptonite/variables"
)

func main() {

	// Parse input arguments
	var address, privateKey, stableToken, targetToken, botToken, password, redisAddress, mode string
	var chainID, profitPercent, stopLossPercent uint64
	flag.StringVar(&address, "publicKey", "", "Your wallet public address")
	flag.StringVar(&privateKey, "privateKey", "", "Your wallet private key")
	flag.StringVar(&mode, "mode", "MANUAL", "Mode of operation. Allowed options: 'MANUAL' (requires you to authorize a swap via telgram bot), 'AUTO' (does not require any confirmation)")
	flag.StringVar(&bot.ChatID, "chatId", "", "Your telegram chat id. You will receive this when you authorize yourself with the bot for the first time")
	flag.StringVar(&stableToken, "stableToken", "USDC", "Stable token (ERC20) to use. Example: USDC, USDT, DAI")
	flag.StringVar(&targetToken, "targetToken", "WETH", "Target ERC20 token to hold. Example: WETH, WMATIC, LINK.")
	flag.StringVar(&botToken, "botToken", "", "Telegram bot token used to send and receive transaction confirmations")
	flag.StringVar(&password, "password", "kryptonite", "Password to share with the bot to authorize yourself as the admin")
	flag.StringVar(&redisAddress, "redisAddress", "", "Redis server host. Example: 192.168.1.100:6379")
	flag.Uint64Var(&chainID, "chainId", 1, "Chain to use. Allowed options: 1 (Ethereum), 10 (Optimism), 56 (Binance Smart Chain), 137 (Polygon/Matic), 42161 (Arbitrum)")
	flag.Uint64Var(&profitPercent, "profitPercent", 50, "Profit percent at which the bot will execute a sell order")
	flag.Uint64Var(&stopLossPercent, "stopLossPercent", 25, "Loss percent at which the bot will execute a stop loss order")
	flag.Parse()

	if botToken == "" {
		log.Fatalln("botToken is not provided")
	}

	go bot.Run(botToken, password)
	time.Sleep(2 * time.Second)

	if privateKey == "" {
		err := "privateKey is not provided"
		log.Fatalln(err)
		bot.OutboundChannel <- err
	}

	if redisAddress == "" {
		err := "redis address is not provided"
		bot.OutboundChannel <- err
		log.Fatalln(err)
	}

	if (mode != "MANUAL") && (mode != "AUTO") {
		err := "unsupported mode"
		bot.OutboundChannel <- err
		log.Fatalln(err)
	}

	// Connect to redis
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	// Check redis ping
	if err := rdb.Ping(context.TODO()).Err(); err != nil {
		bot.OutboundChannel <- err.Error()
		log.Fatalln(err)
	}

	// Initialize wallet
	wallet, err := (&models.Wallet{}).New(address, privateKey, chainID)
	if err != nil {
		bot.OutboundChannel <- err.Error()
		log.Fatalln(err)
	}

	// Initialize swap router
	router, err := (&models.Router{}).New(chainID)
	if err != nil {
		bot.OutboundChannel <- err.Error()
		log.Fatalln(err)
	}

	// Get coin list from CoinGecko
	coins, err := coingecko.GetCoinsList()
	if err != nil {
		bot.OutboundChannel <- err.Error()
		log.Fatalln(err)
	}

	targetCoinID := ""
	for _, coin := range coins {
		if coin.Symbol == strings.ToLower(targetToken) {
			targetCoinID = coin.ID
			break
		}
	}
	if targetCoinID == "" {
		err := "coin could not be found"
		bot.OutboundChannel <- err
		log.Fatalln(err)
	}

	stableTokenContractAddress, targetTokenContractAddress := helpers.GetTokenAddress(router, stableToken), helpers.GetTokenAddress(router, targetToken)
	if stableTokenContractAddress == "" || targetTokenContractAddress == "" {
		err := "token pair does not exist on this chain/network/router"
		bot.OutboundChannel <- err
		log.Fatalln(err)
	} else {
		currentStatus := "UNKNOWN"

		// Real process starts here
		for {

			// Refresh wallet balance periodically
			if err := wallet.RefreshWalletBalance(); err != nil {
				bot.OutboundChannel <- err.Error()
				log.Println(err)
			} else {

				// Print wallet address, wallet balance and router contract address on this chain/network
				log.Print(fmt.Sprintf("Wallet Main Balance: %s (wei)", wallet.MainAccountBalance.String()))
				log.Print(fmt.Sprintf("Router Contract Address: %s", router.Address.Hex()))

				// Refresh ERC20 token balances
				if err := wallet.RefreshTokenBalances(stableTokenContractAddress, targetTokenContractAddress); err != nil {
					bot.OutboundChannel <- err.Error()
					log.Println(err)
				} else {

					// Check current status
					if wallet.StableCoinBalance.Cmp(big.NewInt(0)) == 0 && wallet.TargetCoinBalance.Cmp(big.NewInt(0)) != 0 {
						currentStatus = "WAITING_TO_SELL"
					}

					if wallet.StableCoinBalance.Cmp(big.NewInt(0)) != 0 && wallet.TargetCoinBalance.Cmp(big.NewInt(0)) == 0 {
						currentStatus = "WAITING_TO_BUY"
					}

					// Print token addresses and balances
					log.Print(fmt.Sprintf("%s [%s] => %s (wei)", targetToken, targetTokenContractAddress, wallet.TargetCoinBalance.String()))
					log.Print(fmt.Sprintf("%s [%s] => %s (wei)", stableToken, stableTokenContractAddress, wallet.StableCoinBalance.String()))

					// Get current token price
					currentTokenPrice, err := coingecko.GetCoinPrice(targetCoinID)
					if err != nil {
						bot.OutboundChannel <- err.Error()
						log.Println(err)
					} else {

						// Get market history of the token since the last 1 year
						data, err := coingecko.GetMarketChartByCoin(targetCoinID, 365)
						if err != nil {
							bot.OutboundChannel <- err.Error()
							log.Println(err)
						} else {

							oneYearPoints := make([]float64, 0)
							for _, point := range data.Prices {
								oneYearPoints = append(oneYearPoints, point[1])
							}

							// Calculate pivot points
							supports, resistances, err := technical.GetSupportsAndResistances(oneYearPoints)
							if err != nil {
								bot.OutboundChannel <- err.Error()
								log.Println(err)
							} else {

								// Analyze the period of interest
								pricesShort := oneYearPoints[(len(oneYearPoints) - int(10)):]
								pricesLong := oneYearPoints[(len(oneYearPoints) - int(30)):]

								if len(pricesShort) > 0 && len(pricesLong) > 0 {

									// Calculate simple moving averages
									movingAverageShort := technical.GetMovingAverage(pricesShort)
									movingAverageLong := technical.GetMovingAverage(pricesLong)
									// Get nearest support and resistance levels
									recentSupport, recentResistance := technical.GetRecentSupportAndResistance(currentTokenPrice, supports, resistances)
									log.Print(fmt.Sprintf("Current Status: %s", currentStatus))

									variables.CurrentStatus = currentStatus
									variables.Verdict = "Verdict => Nothing to do"
									str := ""

									// Check if current status is WAITING_TO_BUY/WAITING_TO_SELL/UNKNOWN
									if currentStatus == "WAITING_TO_BUY" && wallet.StableCoinBalance.Cmp(big.NewInt(0)) > 0 {

										// Check if currentTokenPrice is a BUY
										bool, upside, downside := technical.IsABuy(currentTokenPrice, movingAverageShort, movingAverageLong, recentSupport, recentResistance, int64(profitPercent), int64(stopLossPercent))

										// If currentTokenPrice is a BUY
										if bool {
											// Print stats
											str = fmt.Sprintf("Verdict => BUY (Upside: +%.2f%s, Downside: %.2f%s)", upside, "%", downside, "%")
											if err := router.DoSwap(wallet, stableTokenContractAddress, wallet.StableCoinBalance, targetTokenContractAddress, mode); err != nil {
												if err.Error() == "REQUEST_EXPIRED_OR_DECLINED" {
													err = errors.New("Request expired/declined")
													bot.OutboundChannel <- err.Error()
													log.Println(err)
												} else {
													Die(err)
												}
											} else {
												_, err := rdb.HSet(context.TODO(), strings.ToUpper(stableToken)+"_"+strings.ToUpper(targetToken), "PreviousTokenPrice", currentTokenPrice, 0).Result()
												if err != nil {
													Die(err)
												}
											}
										} else {
											// currentTokenPrice is not a BUY, HODL
											str = "Verdict => HODL"
										}
										str = fmt.Sprintf("%s\n\nCurrent Price: $%f\nAverage Price (10 days): $%f\nRecent Support: $%f\nRecent Resistance: $%f\nUpside: +%f%s\nDownside: %f%s\n", str, currentTokenPrice, movingAverageShort, recentSupport, recentResistance, upside, "%", downside, "%")
									} else if currentStatus == "WAITING_TO_SELL" && wallet.TargetCoinBalance.Cmp(big.NewInt(0)) == 1 {
										// Get the price at which the token was last bought
										v, err := rdb.HGet(context.TODO(), strings.ToUpper(stableToken)+"_"+strings.ToUpper(targetToken), "PreviousTokenPrice").Result()
										if err != nil {
											bot.OutboundChannel <- err.Error()
											log.Println(err)
										} else {
											// Convert price from string to float64
											previousTokenPrice, err := strconv.ParseFloat(v, 64)
											if err != nil {
												bot.OutboundChannel <- err.Error()
												log.Println(err)
											} else {
												// Check if currentTokenPrice is a SELL based on technical analysis, previousTokenPrice, profitPercent and stopLossPercent
												isASell, typ, value := technical.IsASell(previousTokenPrice, currentTokenPrice, movingAverageShort, movingAverageLong, recentSupport, recentResistance, int64(profitPercent), int64(stopLossPercent))

												// If currentTokenPrice is a SELL
												if isASell {
													// Print stats
													str = fmt.Sprintf("Verdict => SELL (%s: %.2f%s)", typ, value, "%")
													if err := router.DoSwap(wallet, targetTokenContractAddress, wallet.TargetCoinBalance, stableTokenContractAddress, mode); err != nil {
														if err.Error() == "REQUEST_EXPIRED_OR_DECLINED" {
															err = errors.New("Request expired/declined")
															bot.OutboundChannel <- err.Error()
															log.Println(err)
														} else {
															Die(err)
														}
													} else {
														_, err := rdb.HSet(context.TODO(), strings.ToUpper(stableToken)+"_"+strings.ToUpper(targetToken), "PreviousTokenPrice", 0, 0).Result()
														if err != nil {
															Die(err)
														}
													}
												} else {
													// currentTokenPrice is not a SELL, HODL
													str = "Verdict => HODL"
												}
												str = fmt.Sprintf("%s\n\nCurrent Price: $%f\nAverage Price (10 days): $%f\nRecent Support: $%f\nRecent Resistance: $%f\n%s: %f%s\n", str, currentTokenPrice, movingAverageShort, recentSupport, recentResistance, typ, value, "%")
											}
										}
									}
									variables.Verdict = str
									log.Println(variables.Verdict)
								}
							}
						}
					}
				}
			}
			// Sleep for 1 minute before going back to work
			time.Sleep(1 * time.Minute)
		}
	}
}

func Die(err error) {
	bot.OutboundChannel <- err.Error()
	log.Println(err)
	err = errors.New("Bot is now going to die.")
	bot.OutboundChannel <- err.Error()
	log.Println(err)
	time.Sleep(5 * time.Second)
	os.Exit(1)
}
