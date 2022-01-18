package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
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

	if err := variables.ParseFlags(); err != nil {
		log.Fatalln(err)
	}

	go bot.Run(variables.BotToken, variables.Password)
	time.Sleep(2 * time.Second)

	// Connect to redis
	rdb := redis.NewClient(&redis.Options{
		Addr: variables.RedisAddress,
	})

	// Check redis ping
	if err := rdb.Ping(context.TODO()).Err(); err != nil {
		Die(err)
	}

	// Initialize wallet
	wallet, err := (&models.Wallet{}).New(variables.Address, variables.PrivateKey, variables.ChainID)
	if err != nil {
		Die(err)
	}

	// Initialize swap router
	router, err := (&models.Router{}).New(variables.ChainID)
	if err != nil {
		Die(err)
	}

	targetCoinID, err := coingecko.GetCoinID(variables.TargetToken)
	if err != nil {
		Die(err)
	}

	stableTokenContractAddress, targetTokenContractAddress := helpers.GetTokenAddress(router, variables.StableToken), helpers.GetTokenAddress(router, variables.TargetToken)
	if stableTokenContractAddress == "" || targetTokenContractAddress == "" {
		err := "token pair does not exist on this chain/network/router"
		Die(errors.New(err))
	} else {
		currentStatus := "UNKNOWN"

		// Real process starts here
		for {

			// Refresh wallet balance periodically
			if err := wallet.RefreshWalletBalance(); err != nil {
				bot.OutboundChannel <- err.Error()
			} else {

				// Print wallet address, wallet balance and router contract address on this chain/network
				log.Print(fmt.Sprintf("Wallet Main Balance: %s (wei)", wallet.MainAccountBalance.String()))
				log.Print(fmt.Sprintf("Router Contract Address: %s", router.Address.Hex()))

				// Refresh ERC20 token balances
				if err := wallet.RefreshTokenBalances(stableTokenContractAddress, targetTokenContractAddress); err != nil {
					bot.OutboundChannel <- err.Error()
				} else {

					// Check current status
					if wallet.StableCoinBalance.Cmp(big.NewInt(0)) == 0 && wallet.TargetCoinBalance.Cmp(big.NewInt(0)) != 0 {
						currentStatus = "WAITING_TO_SELL"
					}

					if wallet.StableCoinBalance.Cmp(big.NewInt(0)) != 0 && wallet.TargetCoinBalance.Cmp(big.NewInt(0)) == 0 {
						currentStatus = "WAITING_TO_BUY"
					}

					// Print token addresses and balances
					log.Print(fmt.Sprintf("%s [%s] => %s (wei)", variables.TargetToken, targetTokenContractAddress, wallet.TargetCoinBalance.String()))
					log.Print(fmt.Sprintf("%s [%s] => %s (wei)", variables.StableToken, stableTokenContractAddress, wallet.StableCoinBalance.String()))

					// Get current token price
					currentTokenPrice, err := coingecko.GetCoinPrice(targetCoinID)
					if err != nil {
						bot.OutboundChannel <- err.Error()
					} else {
						log.Print(fmt.Sprintf("Current Status: %s", currentStatus))
						variables.CurrentStatus = currentStatus
						variables.Verdict = "Verdict => Nothing to do"
						str := ""

						// Check if current status is WAITING_TO_BUY/WAITING_TO_SELL/UNKNOWN
						if currentStatus == "WAITING_TO_BUY" && wallet.StableCoinBalance.Cmp(big.NewInt(0)) > 0 {
							v, err := rdb.HGet(context.TODO(), strings.ToUpper(variables.StableToken)+"_"+strings.ToUpper(variables.TargetToken), "BuyLimit").Result()
							if err == redis.Nil {
								_, err := rdb.HSet(context.TODO(), strings.ToUpper(variables.StableToken)+"_"+strings.ToUpper(variables.TargetToken), "BuyLimit", 0).Result()
								if err != nil {
									Die(err)
								}
							} else if err != nil {
								bot.OutboundChannel <- err.Error()
							} else {
								// Convert price from string to float64
								buyLimit, err := strconv.ParseFloat(v, 64)
								if err != nil {
									bot.OutboundChannel <- err.Error()
								} else {
									// If currentTokenPrice is a BUY
									quoteParams := &models.QuoteParamsDto{
										FromTokenAddress: stableTokenContractAddress,
										ToTokenAddress:   targetTokenContractAddress,
										Amount:           wallet.StableCoinBalance.String(),
									}
									quoteResDto, err := router.GetQuote(router.ChainID, quoteParams)
									if err != nil {
										bot.OutboundChannel <- err.Error()
									} else {

										toTokenAmount := new(big.Float)
										toTokenAmount.SetString(quoteResDto.ToTokenAmount)
										s := new(big.Float).Quo(toTokenAmount, big.NewFloat(math.Pow10(int(quoteResDto.ToToken.Decimals))))

										fromTokenAmount := new(big.Float)
										fromTokenAmount.SetString(quoteResDto.FromTokenAmount)
										m := new(big.Float).Quo(fromTokenAmount, big.NewFloat(math.Pow10(int(quoteResDto.FromToken.Decimals))))

										tokenAmountFromSwap, err := strconv.ParseFloat(s.String(), 64)
										if err != nil {
											bot.OutboundChannel <- err.Error()
										} else {
											maxTokensReceivableAccordingToPrice, err := strconv.ParseFloat(m.String(), 64)
											if err != nil {
												bot.OutboundChannel <- err.Error()
											} else {
												tokenAmountFromPrice := maxTokensReceivableAccordingToPrice / currentTokenPrice
												if technical.IsABuy(currentTokenPrice, buyLimit, tokenAmountFromSwap, tokenAmountFromPrice) {
													// Print stats
													str = fmt.Sprintf("Verdict => BUY %s %s with %s %s\n\nCurrent Price: $%f\nBuy Limit: $%f\n\nCurrent Price has reached Buy Limit Price and the liquidity ratio is good in the pool\n", s.String(), quoteResDto.ToToken.Symbol, m.String(), quoteResDto.FromToken.Symbol, currentTokenPrice, buyLimit)
													variables.Verdict = str
													bot.OutboundChannel <- variables.Verdict
													if err := router.DoSwap(wallet, stableTokenContractAddress, wallet.StableCoinBalance, targetTokenContractAddress, variables.Mode); err != nil {
														if err.Error() == "REQUEST_EXPIRED_OR_DECLINED" {
															err = errors.New("Request expired/declined")
															bot.OutboundChannel <- err.Error()
														} else {
															Die(err)
														}
													} else {
														_, err := rdb.HSet(context.TODO(), strings.ToUpper(variables.StableToken)+"_"+strings.ToUpper(variables.TargetToken), "PreviousTokenPrice", currentTokenPrice).Result()
														if err != nil {
															Die(err)
														} else {
															_, err := rdb.HSet(context.TODO(), strings.ToUpper(variables.StableToken)+"_"+strings.ToUpper(variables.TargetToken), "SellLimit", currentTokenPrice*float64(1+(variables.ProfitPercent/100))).Result()
															if err != nil {
																Die(err)
															} else {
																_, err := rdb.HSet(context.TODO(), strings.ToUpper(variables.StableToken)+"_"+strings.ToUpper(variables.TargetToken), "StopLimit", currentTokenPrice*float64(1-(variables.StopLossPercent/100))).Result()
																if err != nil {
																	Die(err)
																}
															}
														}
													}
												} else {
													// currentTokenPrice is not a BUY, HODL
													str = fmt.Sprintf("Verdict => HODL\n\nCurrent Price: $%f\nBuy Limit: $%f\n\nCurrent Price is higher than Buy Limit Price or not enough liquidity in the pool\n", currentTokenPrice, buyLimit)
													variables.Verdict = str
												}
											}
										}
									}
								}
							}
						} else if currentStatus == "WAITING_TO_SELL" && wallet.TargetCoinBalance.Cmp(big.NewInt(0)) == 1 {
							v, err := rdb.HGet(context.TODO(), strings.ToUpper(variables.StableToken)+"_"+strings.ToUpper(variables.TargetToken), "SellLimit").Result()
							if err == redis.Nil {
								e := "No Sell limit is set. Please set a sell limit in the web interface or through the bot."
								bot.OutboundChannel <- e
							} else if err != nil {
								bot.OutboundChannel <- err.Error()
							} else {
								sellLimit, err := strconv.ParseFloat(v, 64)
								if err != nil {
									bot.OutboundChannel <- err.Error()
								} else {
									v, err := rdb.HGet(context.TODO(), strings.ToUpper(variables.StableToken)+"_"+strings.ToUpper(variables.TargetToken), "StopLimit").Result()
									if err == redis.Nil {
										e := "No stop limit is set. Please set a stop limit in the web interface or through the bot."
										bot.OutboundChannel <- e
									} else if err != nil {
										bot.OutboundChannel <- err.Error()
									} else {
										stopLimit, err := strconv.ParseFloat(v, 64)
										if err != nil {
											bot.OutboundChannel <- err.Error()
										} else {
											v, err := rdb.HGet(context.TODO(), strings.ToUpper(variables.StableToken)+"_"+strings.ToUpper(variables.TargetToken), "PreviousTokenPrice").Result()
											if err == redis.Nil {
												e := "Could not find previous token price."
												bot.OutboundChannel <- e
											} else if err != nil {
												bot.OutboundChannel <- err.Error()
											} else {
												previousTokenPrice, err := strconv.ParseFloat(v, 64)
												if err != nil {
													bot.OutboundChannel <- err.Error()
												} else {
													// Check if currentTokenPrice is a SELL
													isASell, typ, value := technical.IsASell(currentTokenPrice, previousTokenPrice, sellLimit, stopLimit)

													if isASell {
														// Print stats
														str = fmt.Sprintf("Verdict => SELL (%s: %.2f%s)", typ, value, "%")
														variables.Verdict = str
														bot.OutboundChannel <- variables.Verdict
														if err := router.DoSwap(wallet, targetTokenContractAddress, wallet.TargetCoinBalance, stableTokenContractAddress, variables.Mode); err != nil {
															if err.Error() == "REQUEST_EXPIRED_OR_DECLINED" {
																err = errors.New("Request expired/declined")
																bot.OutboundChannel <- err.Error()
															} else {
																Die(err)
															}
														} else {
															_, err := rdb.HSet(context.TODO(), strings.ToUpper(variables.StableToken)+"_"+strings.ToUpper(variables.TargetToken), "BuyLimit", 0).Result()
															if err != nil {
																Die(err)
															}
														}
													} else {
														// currentTokenPrice is not a SELL, HODL
														str = fmt.Sprintf("Verdict => HODL\n\nCurrent Price: $%f\nSell Limit: $%f\n\nCurrent Price is lower than Sell Limit Price\n", currentTokenPrice, sellLimit)
														variables.Verdict = str
													}
												}
											}
										}
									}
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
	err = errors.New("Kryptonite bot is now going to exit. Good bye.")
	bot.OutboundChannel <- err.Error()
	time.Sleep(5 * time.Second)
	os.Exit(1)
}
