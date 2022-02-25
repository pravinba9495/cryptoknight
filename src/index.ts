import { Args } from "./utils/flags";
import Telegram from "./api/telegram";
import { Wait } from "./utils/wait";
import { NewClient } from "./redis";
import { Wallet } from "./api/wallet";
import { Router } from "./api/oneinch";
import {
  GetTradeSignal,
  InitTradingViewTechnicals,
  IsPuppeteerReady,
} from "./utils/puppet";
import { InitNgRok } from "./utils/ngrok";
import { GetCurrentStatus } from "./utils/status";
import { Approve } from "./utils/approve";
import { Alternative } from "./api/alternative";
import { Kraken } from "./api/kraken";
import { Forever } from "./utils/forever";
import { PrepareForSwap } from "./utils/prepare";

process.on("uncaughtException", async (error) => {
  console.error(error);
  try {
    await Forever(async () => {
      await Telegram.SendMessage(
        Args.botToken,
        Args.chatId,
        "Uncaught Exception"
      );
    }, 2);
  } catch (error) {
    console.error(error);
  } finally {
    process.exit(1);
  }
});

process.on("unhandledRejection", async (error) => {
  console.error(error);
  try {
    await Forever(async () => {
      await Telegram.SendMessage(
        Args.botToken,
        Args.chatId,
        "Unhandled Rejection"
      );
    }, 2);
  } catch (error) {
    console.error(error);
  } finally {
    process.exit(1);
  }
});

let LAST_TELEGRAM_SIGNAL = "";
const ORIGINAL_MODE = Args.mode;

(async () => {
  try {
    // ====== BEGIN (Main Steps) ======= //

    let routerAddress = "";
    let tokens: any[] = [];
    let ngRokURL = "";
    let stableTokenBalance = BigInt(0);
    let targetTokenBalance = BigInt(0);
    let stableTokenCurrentPrice = 0;
    let targetTokenCurrentPrice = 0;
    let currentStatus = "";
    let signal = "";

    const redis = await NewClient(Args.redisAddress);
    await redis.del(`${Args.stableToken}_${Args.targetToken}`);
    const wallet = new Wallet(Args.publicKey, Args.privateKey, Args.chainId);
    const router = new Router(Args.chainId);

    await Forever(async () => {
      routerAddress = await router.GetContractAddress();
      tokens = await router.GetSupportedTokens();
    }, 2);

    await Forever(async () => {
      ngRokURL = await InitNgRok(Args.port);
      console.log(`ngrok tunnel running at: ${ngRokURL}`);
    }, 2);

    await Forever(async () => {
      await Telegram.SetWebhook(Args.botToken, ngRokURL);
    }, 2);

    const stableTokenContractAddress =
      tokens.find((token) => token.symbol === Args.stableToken)?.address || "";
    const targetTokenContractAddress =
      tokens.find((token) => token.symbol === Args.targetToken)?.address || "";
    if (
      stableTokenContractAddress === "" ||
      targetTokenContractAddress === ""
    ) {
      throw "Token Contract Address cannot be empty";
    }

    await Forever(async () => {
      stableTokenBalance = await wallet.GetTokenBalance(
        stableTokenContractAddress
      );
      targetTokenBalance = await wallet.GetTokenBalance(
        targetTokenContractAddress
      );
    }, 2);

    InitTradingViewTechnicals(Args.targetTokenTickerKraken, Args.chartInterval);
    while (true) {
      if (IsPuppeteerReady()) {
        break;
      } else {
        console.log("Waiting for puppeteer to be ready");
        await Wait(2);
      }
    }

    const t = 5 * 60;
    await Forever(async () => {
      await redis.setEx(
        "LAST_SIGNAL_UPDATE",
        t,
        new Date().getTime().toString()
      );
    }, 2);

    if (Args.preAuth) {
      await Forever(async () => {
        const stableTokenAllowance = await router.GetApprovedAllowance(
          stableTokenContractAddress,
          wallet.Address
        );
        const targetTokenAllowance = await router.GetApprovedAllowance(
          targetTokenContractAddress,
          wallet.Address
        );
        if (stableTokenAllowance >= stableTokenBalance) {
          console.log(
            `Router already preauthorized to spend ${Args.stableToken}`
          );
        } else {
          console.log(`Preauthorizing router for ${Args.stableToken}`);
          let approveTxHash = "";
          await Forever(async () => {
            approveTxHash = await Approve(
              wallet,
              router,
              stableTokenContractAddress,
              "-1"
            );
          }, 2);
          if (approveTxHash === "") {
            console.error(`Approve transaction for ${Args.stableToken} failed`);
            await Forever(async () => {
              await Telegram.SendMessage(
                Args.botToken,
                Args.chatId,
                `[STOPPING] Approve transaction for ${Args.stableToken} failed`
              );
            }, 2);
            process.exit(1);
          } else {
            console.error(
              `Approve transaction ${approveTxHash} for ${Args.stableToken} succeeded`
            );
          }
        }
        if (targetTokenAllowance >= targetTokenBalance) {
          console.log(
            `Router already preauthorized to spend ${Args.targetToken}`
          );
        } else {
          console.log(`Preauthorizing router for ${Args.targetToken}`);
          let approveTxHash = "";
          await Forever(async () => {
            approveTxHash = await Approve(
              wallet,
              router,
              targetTokenContractAddress,
              "-1"
            );
          }, 2);
          if (approveTxHash === "") {
            console.error(`Approve transaction for ${Args.targetToken} failed`);
            await Forever(async () => {
              await Telegram.SendMessage(
                Args.botToken,
                Args.chatId,
                `[STOPPING] Approve transaction for ${Args.targetToken} failed`
              );
            }, 2);
            process.exit(1);
          } else {
            console.error(
              `Approve transaction ${approveTxHash} for ${Args.targetToken} succeeded`
            );
          }
        }
      }, 2);
    }

    // ====== END (Main Steps) ======= //

    // ====== BEGIN (Loop) ======= //
    while (true) {
      const start = new Date().getTime();
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

      await Forever(async () => {
        currentStatus = GetCurrentStatus(
          stableTokenBalance,
          targetTokenBalance
        );
        console.log(`Current Status: ${currentStatus}`);
      }, 2);

      await Forever(async () => {
        const gasPrice = await wallet.SuggestGasPrice();
        console.log(
          `Current Gas Price (Gwei): ${Number(gasPrice) / 1000000000}`
        );
      }, 2);

      await Forever(async () => {
        const { fearGreedIndex, fearGreedIndexClassification } =
          await Alternative.GetCryptoFearIndex();
        console.log(
          `Current Fear/Greed Level: ${fearGreedIndexClassification} (${fearGreedIndex})`
        );
        if (fearGreedIndex <= Args.switchModeLimit) {
          if (Args.mode !== "MANUAL") {
            Args.mode = "MANUAL";
            await Forever(async () => {
              const message = `Switching mode to MANUAL due to fear in the market`;
              console.log(message);
              await Telegram.SendMessage(Args.botToken, Args.chatId, message);
              await redis.del(`${Args.stableToken}_${Args.targetToken}`);
            }, 2);
          }
        } else {
          if (Args.mode !== ORIGINAL_MODE) {
            Args.mode = ORIGINAL_MODE;
            await Forever(async () => {
              const message = `Switching mode to ${Args.mode} due to hope in the market`;
              console.log(message);
              await Telegram.SendMessage(Args.botToken, Args.chatId, message);
              await redis.del(`${Args.stableToken}_${Args.targetToken}`);
            }, 2);
          }
        }
      }, 2);

      await Forever(async () => {
        signal = await GetTradeSignal();
      }, 2);

      if (
        signal.includes("STRONG BUY") ||
        signal.includes("STRONG SELL") ||
        signal.includes("WEAK")
      ) {
        await Forever(async () => {
          await redis.setEx(
            "LAST_SIGNAL_UPDATE",
            t,
            new Date().getTime().toString()
          );
        }, 2);
      }

      if (
        (signal.includes("STRONG BUY") || signal.includes("STRONG SELL")) &&
        LAST_TELEGRAM_SIGNAL !== signal
      ) {
        await Forever(async () => {
          await Telegram.SendMessage(Args.botToken, Args.chatId, signal);
        }, 2);
      }
      LAST_TELEGRAM_SIGNAL = signal;

      await Forever(async () => {
        const exists = await redis.exists("LAST_SIGNAL_UPDATE");
        if (exists !== 1) {
          await Telegram.SendMessage(
            Args.botToken,
            Args.chatId,
            `Did not receive valid signal for more than ${t} seconds.`
          );
        }
      }, 2);

      await Forever(async () => {
        stableTokenCurrentPrice = await Kraken.GetCoinPrice(
          Args.stableTokenTickerKraken
        );
        targetTokenCurrentPrice = await Kraken.GetCoinPrice(
          Args.targetTokenTickerKraken
        );
      }, 2);

      if (currentStatus === "WAITING_TO_BUY") {
        let buyLimitPrice = 0;
        let buyBackLimitPrice = 9999999999;
        await Forever(async () => {
          buyLimitPrice =
            Number(
              await redis.hGet(
                `${Args.stableToken}_${Args.targetToken}`,
                "BuyLimitPrice"
              )
            ) || 0;
          buyBackLimitPrice =
            Number(
              await redis.hGet(
                `${Args.stableToken}_${Args.targetToken}`,
                "BuyBackLimitPrice"
              )
            ) || 9999999999;
        }, 2);
        const params = {
          fromTokenAddress: stableTokenContractAddress,
          toTokenAddress: targetTokenContractAddress,
          amount: stableTokenBalance.toString(),
        };

        let quoteResponseDto: any = {};
        await Forever(async () => {
          quoteResponseDto = await router.GetQuote(params);
        }, 2);

        const stableTokenAmnt =
          Number(stableTokenBalance) /
          Math.pow(10, quoteResponseDto.fromToken.decimals);
        const currentPortfolioValue = stableTokenAmnt * stableTokenCurrentPrice;
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
          `Current Status: ${currentStatus}, Current Mode: ${Args.mode}, Current Signal: ${signal}`
        );

        const buyLimitReached = buyLimitPrice >= targetTokenCurrentPrice;
        const buyBackLimitReached =
          targetTokenCurrentPrice >= buyBackLimitPrice;

        if (
          (signal === "STRONG BUY" && Args.mode === "AUTO") ||
          ((buyLimitReached || buyBackLimitReached) && Args.mode === "MANUAL")
        ) {
          // Liquidity provider fee: 0.5% approx
          if (
            actualSlippage <= Args.slippagePercent + 0.5 ||
            buyLimitReached ||
            buyBackLimitReached
          ) {
            console.log(
              `BUY (Current Price: $${targetTokenCurrentPrice}, Slippage: ${actualSlippage.toFixed(
                2
              )}%, Slippage Allowed: +${
                Args.slippagePercent
              }%, Current Portfolio Value: $${currentPortfolioValue}, Minimum Return: ${toTokenAmount} ${
                quoteResponseDto.toToken.symbol
              })`
            );

            await Forever(async () => {
              await Telegram.SendMessage(
                Args.botToken,
                Args.chatId,
                `Signal Received: ${signal}`
              );
            }, 2);

            await PrepareForSwap(
              router,
              wallet,
              stableTokenContractAddress,
              stableTokenBalance,
              targetTokenContractAddress
            );

            await Forever(async () => {
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
                "BuyBackLimitPrice",
                9999999999
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
            }, 2);

            let bal = BigInt(0);
            let balAmnt = 0;

            await Forever(async () => {
              bal = await wallet.GetTokenBalance(targetTokenContractAddress);
              balAmnt =
                Number(bal) / Math.pow(10, quoteResponseDto.toToken.decimals);
              stableTokenCurrentPrice = await Kraken.GetCoinPrice(
                Args.stableTokenTickerKraken
              );
              targetTokenCurrentPrice = await Kraken.GetCoinPrice(
                Args.targetTokenTickerKraken
              );
            }, 2);

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

            await Forever(async () => {
              await redis.lPush(
                `${Args.stableToken}_${Args.targetToken}_SWAP_HISTORY`,
                JSON.stringify(trade)
              );
            }, 2);

            await Forever(async () => {
              await Telegram.SendMessage(
                Args.botToken,
                Args.chatId,
                JSON.stringify(trade, null, 2)
              );
            }, 2);

            if (Args.mode == "AUTO") {
              await Forever(async () => {
                await redis.hSet(
                  `${Args.stableToken}_${Args.targetToken}`,
                  "StopLimitPrice",
                  ((100 - Args.stopLimitPercent) / 100) *
                    targetTokenCurrentPrice
                );
              }, 2);
            }

            currentStatus = "WAITING_TO_SELL";

            await Forever(async () => {
              stableTokenBalance = await wallet.GetTokenBalance(
                stableTokenContractAddress
              );
              targetTokenBalance = await wallet.GetTokenBalance(
                targetTokenContractAddress
              );
            }, 2);
          } else {
            console.log(
              `HOLD (Current Price: $${targetTokenCurrentPrice}, Slippage: ${actualSlippage.toFixed(
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
            `HOLD (Current Price: $${targetTokenCurrentPrice}, Slippage Allowed: +${Args.slippagePercent}%, Current Portfolio Value: $${currentPortfolioValue}, Minimum Return: ${toTokenAmount} ${quoteResponseDto.toToken.symbol})`
          );
        }
      } else if (currentStatus === "WAITING_TO_SELL") {
        let lastBuyPrice = 0;
        let stopLimitPrice = 0;
        let sellLimitPrice = 9999999999;

        await Forever(async () => {
          lastBuyPrice =
            Number(
              await redis.hGet(
                `${Args.stableToken}_${Args.targetToken}`,
                "LastBuyPrice"
              )
            ) || 0;

          if (lastBuyPrice === 0) {
            const newLastBuyPrice = targetTokenCurrentPrice;
            await Forever(async () => {
              await redis.hSet(
                `${Args.stableToken}_${Args.targetToken}`,
                "LastBuyPrice",
                newLastBuyPrice
              );
            }, 2);
            lastBuyPrice = newLastBuyPrice;
          }

          stopLimitPrice =
            Number(
              await redis.hGet(
                `${Args.stableToken}_${Args.targetToken}`,
                "StopLimitPrice"
              )
            ) || 0;

          if (Args.mode == "AUTO" && stopLimitPrice === 0) {
            const newStopLimitPrice =
              ((100 - Args.stopLimitPercent) / 100) * lastBuyPrice;
            await Forever(async () => {
              await redis.hSet(
                `${Args.stableToken}_${Args.targetToken}`,
                "StopLimitPrice",
                newStopLimitPrice
              );
            }, 2);
            stopLimitPrice = newStopLimitPrice;
          }

          if (Args.trailingStop && Args.mode === "AUTO") {
            const newStopLimitPrice =
              ((100 - Args.stopLimitPercent) / 100) * targetTokenCurrentPrice;
            if (stopLimitPrice <= newStopLimitPrice) {
              await Forever(async () => {
                await redis.hSet(
                  `${Args.stableToken}_${Args.targetToken}`,
                  "StopLimitPrice",
                  newStopLimitPrice
                );
              }, 2);
              stopLimitPrice = newStopLimitPrice;
            }
          }

          sellLimitPrice =
            Number(
              await redis.hGet(
                `${Args.stableToken}_${Args.targetToken}`,
                "SellLimitPrice"
              )
            ) || 9999999999;
        }, 2);

        const params = {
          fromTokenAddress: targetTokenContractAddress,
          toTokenAddress: stableTokenContractAddress,
          amount: targetTokenBalance.toString(),
        };
        let quoteResponseDto: any = {};

        await Forever(async () => {
          quoteResponseDto = await router.GetQuote(params);
        }, 2);

        const targetTokenAmnt =
          Number(targetTokenBalance) /
          Math.pow(10, quoteResponseDto.fromToken.decimals);
        const currentPortfolioValue = targetTokenAmnt * targetTokenCurrentPrice;
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
          signal === "STRONG SELL" &&
          Args.mode === "AUTO" &&
          (profitOrLossPercent < 0 ||
            Math.abs(profitOrLossPercent) < Args.minProfitPercent)
        ) {
          await Forever(async () => {
            await Telegram.SendMessage(
              Args.botToken,
              Args.chatId,
              `Could Not Sell\nSignal Received: ${signal}, Profit/Loss: ${profitOrLossPercent}%`
            );
          }, 2);
        }

        if (
          (((signal === "STRONG SELL" &&
            profitOrLossPercent >= 0 &&
            Math.abs(profitOrLossPercent) >= Args.minProfitPercent) ||
            stopLimitReached) &&
            Args.mode === "AUTO") ||
          ((sellLimitReached || stopLimitReached) && Args.mode === "MANUAL")
        ) {
          // Liquidity provider fee: 0.5% approx
          if (
            actualSlippage <= Args.slippagePercent + 0.5 ||
            stopLimitReached ||
            sellLimitReached
          ) {
            console.log(
              `SELL (Current Price: $${targetTokenCurrentPrice}, Last Bought Price: $${lastBuyPrice}, Sell Limit Price: $${sellLimitPrice}, Stop Limit Price: $${stopLimitPrice}, Slippage: ${actualSlippage.toFixed(
                2
              )}%, Slippage Allowed: +${
                Args.slippagePercent
              }%, Current Portfolio Value: $${currentPortfolioValue}, Minimum Return: ${toTokenAmount} ${
                quoteResponseDto.toToken.symbol
              }, ${
                profitOrLossPercent > 0 ? "Minimum Profit" : "Maximum Loss"
              }: ${profitOrLossPercent > 0 ? "+" : ""}${profitOrLossPercent}%)`
            );
            await Forever(async () => {
              await Telegram.SendMessage(
                Args.botToken,
                Args.chatId,
                `Signal Received: ${signal}, Profit/Loss: ${profitOrLossPercent}%`
              );
            }, 2);

            await PrepareForSwap(
              router,
              wallet,
              targetTokenContractAddress,
              targetTokenBalance,
              stableTokenContractAddress
            );

            await Forever(async () => {
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
                "BuyBackLimitPrice",
                9999999999
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
            }, 2);

            let bal = BigInt(0);
            let balAmnt = 0;

            await Forever(async () => {
              bal = await wallet.GetTokenBalance(stableTokenContractAddress);
              balAmnt =
                Number(bal) / Math.pow(10, quoteResponseDto.toToken.decimals);

              stableTokenCurrentPrice = await Kraken.GetCoinPrice(
                Args.stableTokenTickerKraken
              );
              targetTokenCurrentPrice = await Kraken.GetCoinPrice(
                Args.targetTokenTickerKraken
              );
            }, 2);

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

            await Forever(async () => {
              await redis.lPush(
                `${Args.stableToken}_${Args.targetToken}_SWAP_HISTORY`,
                JSON.stringify(trade)
              );
            }, 2);

            await Forever(async () => {
              await Telegram.SendMessage(
                Args.botToken,
                Args.chatId,
                JSON.stringify(trade, null, 2)
              );
            }, 2);

            currentStatus = "WAITING_TO_BUY";

            await Forever(async () => {
              stableTokenBalance = await wallet.GetTokenBalance(
                stableTokenContractAddress
              );
              targetTokenBalance = await wallet.GetTokenBalance(
                targetTokenContractAddress
              );
            }, 2);
          } else {
            console.log(
              `HOLD (Current Price: $${targetTokenCurrentPrice}, Last Bought Price: $${lastBuyPrice}, Sell Limit Price: $${sellLimitPrice}, Stop Limit Price: $${stopLimitPrice}, Slippage: ${actualSlippage.toFixed(
                2
              )}%, Slippage Allowed: +${
                Args.slippagePercent
              }%, Current Portfolio Value: $${currentPortfolioValue}, Minimum Return: ${toTokenAmount} ${
                quoteResponseDto.toToken.symbol
              }, ${
                profitOrLossPercent > 0 ? "Minimum Profit" : "Maximum Loss"
              }: ${profitOrLossPercent > 0 ? "+" : ""}${profitOrLossPercent}%)`
            );
          }
        } else {
          console.log(
            `HOLD (Current Price: $${targetTokenCurrentPrice}, Last Bought Price: $${lastBuyPrice}, Sell Limit Price: $${sellLimitPrice}, Stop Limit Price: $${stopLimitPrice}, Slippage Allowed: +${
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

      // Here
      const end = new Date().getTime();
      console.log(`\nLoop Time: ${(end - start) / 1000} (sec)`);
      await Wait(10);
    }
    // ====== END (Loop) ======= //
  } catch (error) {
    console.error(error);
  } finally {
    process.exit(1);
  }
})();
