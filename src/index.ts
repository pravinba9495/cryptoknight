import { Kraken } from "./api/kraken";
import { Router } from "./api/oneinch";
import { Wallet } from "./api/wallet";
import { Connect } from "./redis";
import { Args } from "./utils/flags";
import { PrepareForSwap } from "./utils/prepare";
import { GetTradeSignal, InitTradingViewTechnicals } from "./utils/puppet";
import { SendMessage } from "./utils/telegram";
import { Wait } from "./utils/wait";

process.on("uncaughtException", (error) => {
  console.error(error);
  process.exit(1);
});

process.on("unhandledRejection", (error) => {
  console.error(error);
  process.exit(1);
});

(async () => {
  try {
    const redis = await Connect(Args.redisAddress);
    const wallet = new Wallet(Args.publicKey, Args.privateKey, Args.chainId);
    const router = new Router(Args.chainId);
    let currentStatus = "UNKNOWN";

    const routerAddress = await router.GetContractAddress();
    const tokens = await router.GetSupportedTokens();

    const stableTokenContractAddress =
      tokens.find((token) => token.symbol === Args.stableToken)?.address || "";
    const targetTokenContractAddress =
      tokens.find((token) => token.symbol === Args.targetToken)?.address || "";
    if (
      stableTokenContractAddress === "" ||
      targetTokenContractAddress === ""
    ) {
      throw new Error("tokenContractAddress cannot be empty");
    }

    InitTradingViewTechnicals(Args.targetTokenTickerKraken, Args.chartInterval);

    while (true) {
      try {
        console.log(`\n\n${new Date()}\n`);
        console.log(`Wallet Address: ${wallet.MaskedAddress}`);
        console.log(`Chain ID: ${wallet.ChainID}`);
        console.log(`Router Contract Address: ${routerAddress}`);
        console.log(
          `Stable Token Contract Address (${Args.stableToken}): ${stableTokenContractAddress}`
        );
        console.log(
          `Target Token Contract Address (${Args.targetToken}): ${targetTokenContractAddress}`
        );
        const stableTokenBalance = await wallet.GetTokenBalance(
          stableTokenContractAddress
        );
        const targetTokenBalance = await wallet.GetTokenBalance(
          targetTokenContractAddress
        );

        if (
          stableTokenBalance !== BigInt(0) &&
          targetTokenBalance === BigInt(0)
        ) {
          currentStatus = "WAITING_TO_BUY";
        }
        if (
          stableTokenBalance === BigInt(0) &&
          targetTokenBalance !== BigInt(0)
        ) {
          currentStatus = "WAITING_TO_SELL";
        }

        let stableTokenCurrentPrice = await Kraken.GetCoinPrice(
          Args.stableTokenTickerKraken
        );
        let targetTokenCurrentPrice = await Kraken.GetCoinPrice(
          Args.targetTokenTickerKraken
        );

        const signal = await GetTradeSignal();

        if (currentStatus === "WAITING_TO_BUY") {
          const buyLimitPrice =
            Number(
              await redis.hGet(
                `${Args.stableToken}_${Args.targetToken}`,
                "BuyLimitPrice"
              )
            ) || 0;
          const params = {
            fromTokenAddress: stableTokenContractAddress,
            toTokenAddress: targetTokenContractAddress,
            amount: stableTokenBalance.toString(),
          };
          const quoteResponseDto = await router.GetQuote(params);
          const stableTokenAmnt =
            Number(stableTokenBalance) /
            Math.pow(10, quoteResponseDto.fromToken.decimals);
          const currentPortfolioValue =
            stableTokenAmnt * stableTokenCurrentPrice;
          const toTokenAmnt =
            Number(quoteResponseDto.toTokenAmount) /
            Math.pow(10, quoteResponseDto.toToken.decimals);
          const toTokenAmount = (1 - Args.slippagePercent / 100) * toTokenAmnt;
          const toTokenValue = toTokenAmount * targetTokenCurrentPrice;
          const actualSlippage =
            ((currentPortfolioValue - toTokenValue) * 100) /
            currentPortfolioValue;

          console.log(
            `Stable Token Balance (${Args.stableToken}): ${stableTokenAmnt} ${Args.stableToken}`
          );
          console.log(
            `Target Token Balance (${Args.targetToken}): ${
              Number(targetTokenBalance) /
              Math.pow(10, quoteResponseDto.toToken.decimals)
            } ${Args.targetToken}`
          );
          console.log(
            `Current Status: ${currentStatus}, Current Signal: ${signal}`
          );

          const buyLimitReached = buyLimitPrice >= targetTokenCurrentPrice;

          if (signal === "BUY" || buyLimitReached) {
            // Liquidity provider fee: 0.5% approx
            if (actualSlippage <= Args.slippagePercent + 0.5) {
              console.log(
                `BUY (Current Price: $${targetTokenCurrentPrice}, Slippage: ${actualSlippage.toFixed(
                  2
                )}%, Slippage Allowed: +${
                  Args.slippagePercent
                }%, Current Portfolio Value: $${currentPortfolioValue}, Minimum Return: ${toTokenAmount} ${
                  quoteResponseDto.toToken.symbol
                })`
              );
              try {
                await SendMessage(
                  Args.botToken,
                  Args.chatId,
                  `Signal Received: ${signal}`
                );
                await PrepareForSwap(
                  router,
                  wallet,
                  stableTokenContractAddress,
                  stableTokenBalance,
                  targetTokenContractAddress
                );
                await redis.hSet(
                  `${Args.stableToken}_${Args.targetToken}`,
                  "BuyLimitPrice",
                  0
                );
                await redis.hSet(
                  `${Args.stableToken}_${Args.targetToken}`,
                  "StopLimitPrice",
                  0
                );
                await redis.hSet(
                  `${Args.stableToken}_${Args.targetToken}`,
                  "SellLimitPrice",
                  9999999999
                );
                await redis.hSet(
                  `${Args.stableToken}_${Args.targetToken}`,
                  "LastBuyPrice",
                  targetTokenCurrentPrice
                );
                while (true) {
                  try {
                    const bal = await wallet.GetTokenBalance(
                      targetTokenContractAddress
                    );
                    const balAmnt =
                      Number(bal) /
                      Math.pow(10, quoteResponseDto.toToken.decimals);
                    stableTokenCurrentPrice = await Kraken.GetCoinPrice(
                      Args.stableTokenTickerKraken
                    );
                    targetTokenCurrentPrice = await Kraken.GetCoinPrice(
                      Args.targetTokenTickerKraken
                    );
                    const trade = {
                      date: new Date().getTime(),
                      sold: Args.stableToken,
                      soldAmount: stableTokenAmnt,
                      soldValue: stableTokenAmnt * stableTokenCurrentPrice,
                      bought: Args.targetToken,
                      boughtAmount: balAmnt,
                      boughtValue: balAmnt * targetTokenCurrentPrice,
                      tradeLossPercent:
                        ((balAmnt * targetTokenCurrentPrice -
                          stableTokenAmnt * stableTokenCurrentPrice) *
                          100) /
                        (stableTokenAmnt * stableTokenCurrentPrice),
                    };
                    await redis.lPush(
                      `${Args.stableToken}_${Args.targetToken}_SWAP_HISTORY`,
                      JSON.stringify(trade)
                    );
                    await SendMessage(
                      Args.botToken,
                      Args.chatId,
                      JSON.stringify(trade, null, 2)
                    );
                    break;
                  } catch (error) {
                    console.error(error);
                    await Wait(5);
                  }
                }
                currentStatus = "WAITING_TO_SELL";
              } catch (error) {
                console.error(error);
                process.exit(1);
              }
            } else {
              console.log(
                `HODL (Current Price: $${targetTokenCurrentPrice}, Slippage: ${actualSlippage.toFixed(
                  2
                )}%, Slippage Allowed: +${
                  Args.slippagePercent
                }%, Current Portfolio Value: $${currentPortfolioValue}, Minimum Return: ${toTokenAmount} ${
                  quoteResponseDto.toToken.symbol
                })`
              );
            }
          } else {
            console.log(
              `HODL (Current Price: $${targetTokenCurrentPrice}, Slippage Allowed: +${Args.slippagePercent}%, Current Portfolio Value: $${currentPortfolioValue}, Minimum Return: ${toTokenAmount} ${quoteResponseDto.toToken.symbol})`
            );
          }
        } else if (currentStatus === "WAITING_TO_SELL") {
          const lastBuyPrice =
            Number(
              await redis.hGet(
                `${Args.stableToken}_${Args.targetToken}`,
                "LastBuyPrice"
              )
            ) || targetTokenCurrentPrice;

          let stopLimitPrice =
            Number(
              await redis.hGet(
                `${Args.stableToken}_${Args.targetToken}`,
                "StopLimitPrice"
              )
            ) || 0;

          const sellLimitPrice =
            Number(
              await redis.hGet(
                `${Args.stableToken}_${Args.targetToken}`,
                "SellLimitPrice"
              )
            ) || 9999999999;

          const params = {
            fromTokenAddress: targetTokenContractAddress,
            toTokenAddress: stableTokenContractAddress,
            amount: targetTokenBalance.toString(),
          };
          const quoteResponseDto = await router.GetQuote(params);

          const targetTokenAmnt =
            Number(targetTokenBalance) /
            Math.pow(10, quoteResponseDto.fromToken.decimals);
          const currentPortfolioValue =
            targetTokenAmnt * targetTokenCurrentPrice;
          const toTokenAmount =
            (1 - Args.slippagePercent / 100) *
            (Number(quoteResponseDto.toTokenAmount) /
              Math.pow(10, quoteResponseDto.toToken.decimals));
          const toTokenValue = toTokenAmount * stableTokenCurrentPrice;
          const actualSlippage =
            ((currentPortfolioValue - toTokenValue) * 100) /
            currentPortfolioValue;
          const profitOrLossPercent =
            Number(
              (
                ((targetTokenCurrentPrice - lastBuyPrice) * 100) /
                lastBuyPrice
              ).toFixed(2)
            ) - actualSlippage;

          console.log(
            `Stable Token Balance (${Args.stableToken}): ${
              Number(stableTokenBalance) /
              Math.pow(10, quoteResponseDto.toToken.decimals)
            } ${Args.stableToken}`
          );
          console.log(
            `Target Token Balance (${Args.targetToken}): ${
              Number(targetTokenBalance) /
              Math.pow(10, quoteResponseDto.fromToken.decimals)
            } ${Args.targetToken}`
          );

          console.log(
            `Current Status: ${currentStatus}, Current Signal: ${signal}`
          );

          const sellLimitReached = targetTokenCurrentPrice >= sellLimitPrice;
          const stopLimitReached = stopLimitPrice >= targetTokenCurrentPrice;

          if (
            (signal === "SELL" &&
              profitOrLossPercent >= Args.minProfitPercent) ||
            sellLimitReached ||
            stopLimitReached
          ) {
            // Liquidity provider fee: 0.5% approx
            if (actualSlippage <= Args.slippagePercent + 0.5) {
              console.log(
                `SELL (Current Price: $${targetTokenCurrentPrice}, Last Bought Price: $${lastBuyPrice}, Slippage: ${actualSlippage.toFixed(
                  2
                )}%, Slippage Allowed: +${
                  Args.slippagePercent
                }%, Current Portfolio Value: $${currentPortfolioValue}, Minimum Return: ${toTokenAmount} ${
                  quoteResponseDto.toToken.symbol
                }, ${
                  profitOrLossPercent > 0 ? "Minimum Profit" : "Maximum Loss"
                }: ${
                  profitOrLossPercent > 0 ? "+" : ""
                }${profitOrLossPercent}%)`
              );
              try {
                await SendMessage(
                  Args.botToken,
                  Args.chatId,
                  `Signal Received: ${signal}`
                );
                await PrepareForSwap(
                  router,
                  wallet,
                  targetTokenContractAddress,
                  targetTokenBalance,
                  stableTokenContractAddress
                );
                await redis.hSet(
                  `${Args.stableToken}_${Args.targetToken}`,
                  "BuyLimitPrice",
                  0
                );
                await redis.hSet(
                  `${Args.stableToken}_${Args.targetToken}`,
                  "StopLimitPrice",
                  0
                );
                await redis.hSet(
                  `${Args.stableToken}_${Args.targetToken}`,
                  "SellLimitPrice",
                  9999999999
                );
                await redis.hSet(
                  `${Args.stableToken}_${Args.targetToken}`,
                  "LastBuyPrice",
                  0
                );
                while (true) {
                  try {
                    const bal = await wallet.GetTokenBalance(
                      stableTokenContractAddress
                    );
                    const balAmnt =
                      Number(bal) /
                      Math.pow(10, quoteResponseDto.toToken.decimals);

                    stableTokenCurrentPrice = await Kraken.GetCoinPrice(
                      Args.stableTokenTickerKraken
                    );
                    targetTokenCurrentPrice = await Kraken.GetCoinPrice(
                      Args.targetTokenTickerKraken
                    );
                    const trade = {
                      date: new Date().getTime(),
                      sold: Args.targetToken,
                      soldAmount: targetTokenAmnt,
                      soldValue: targetTokenAmnt * targetTokenCurrentPrice,
                      bought: Args.stableToken,
                      boughtAmount: balAmnt,
                      boughtValue: balAmnt * stableTokenCurrentPrice,
                      tradeLossPercent:
                        ((balAmnt * stableTokenCurrentPrice -
                          targetTokenAmnt * targetTokenCurrentPrice) *
                          100) /
                        (targetTokenAmnt * targetTokenCurrentPrice),
                    };
                    await redis.lPush(
                      `${Args.stableToken}_${Args.targetToken}_SWAP_HISTORY`,
                      JSON.stringify(trade)
                    );
                    await SendMessage(
                      Args.botToken,
                      Args.chatId,
                      JSON.stringify(trade, null, 2)
                    );
                    break;
                  } catch (error) {
                    console.error(error);
                    await Wait(5);
                  }
                }
                currentStatus = "WAITING_TO_BUY";
              } catch (error) {
                console.error(error);
                process.exit(1);
              }
            } else {
              console.log(
                `HODL (Current Price: $${targetTokenCurrentPrice}, Last Bought Price: $${lastBuyPrice}, Slippage: ${actualSlippage.toFixed(
                  2
                )}%, Slippage Allowed: +${
                  Args.slippagePercent
                }%, Current Portfolio Value: $${currentPortfolioValue}, Minimum Return: ${toTokenAmount} ${
                  quoteResponseDto.toToken.symbol
                }, ${
                  profitOrLossPercent > 0 ? "Minimum Profit" : "Maximum Loss"
                }: ${
                  profitOrLossPercent > 0 ? "+" : ""
                }${profitOrLossPercent}%)`
              );
            }
          } else {
            console.log(
              `HODL (Current Price: $${targetTokenCurrentPrice}, Last Bought Price: $${lastBuyPrice}, Slippage Allowed: +${
                Args.slippagePercent
              }%, Current Portfolio Value: $${currentPortfolioValue}, Minimum Return: ${toTokenAmount} ${
                quoteResponseDto.toToken.symbol
              }, ${
                profitOrLossPercent > 0 ? "Minimum Profit" : "Maximum Loss"
              }: ${profitOrLossPercent > 0 ? "+" : ""}${profitOrLossPercent}%)`
            );
          }
        } else {
          console.log(`Current Status: ${currentStatus}. Nothing to do`);
        }
      } catch (error) {
        console.error(error);
      }
      await Wait(10);
    }
    await redis.disconnect();
  } catch (error) {
    console.error(error);
    process.exit(1);
  }
})();
