import { Router } from "../api/oneinch";
import Telegram from "../api/telegram";
import { Wallet } from "../api/wallet";
import { Approve } from "./approve";
import { Args } from "./flags";
import { Forever } from "./forever";
import { Swap } from "./swap";
import { Wait } from "./wait";

export const PrepareForSwap = async (
  router: Router,
  wallet: Wallet,
  fromTokenContractAddress: string,
  fromTokenBalance: bigint,
  toTokenContractAddress: string
): Promise<void> => {
  let fromTokenAllowance = BigInt(0);

  await Forever(async () => {
    fromTokenAllowance = await router.GetApprovedAllowance(
      fromTokenContractAddress,
      Args.publicKey
    );
  }, 2);

  if (fromTokenAllowance >= fromTokenBalance) {
    console.log(
      `Router is approved to spend the required amount of tokens for a swap`
    );
    const params = {
      fromTokenAddress: fromTokenContractAddress,
      toTokenAddress: toTokenContractAddress,
      amount: fromTokenBalance.toString(),
      fromAddress: wallet.Address,
      slippage: Args.slippagePercent,
      disableEstimate: false,
      allowPartialFill: false,
      gasLimit: 11500000,
    };
    let swapTxHash = "";
    await Forever(async () => {
      swapTxHash = await Swap(wallet, router, params);
    }, 2);
    if (swapTxHash === "") {
      await Forever(async () => {
        await Telegram.SendMessage(
          Args.botToken,
          Args.chatId,
          `[STOPPING] Swap transaction for failed`
        );
      }, 2);
      process.exit(1);
    }
    if (!Args.preAuth) {
      await Forever(async () => {
        await Approve(wallet, router, fromTokenContractAddress, "0");
      }, 2);
      while (true) {
        console.log(`Refreshing router token allowance`);
        let fromTokenAllowance = BigInt(0);
        await Forever(async () => {
          fromTokenAllowance = await router.GetApprovedAllowance(
            fromTokenContractAddress,
            Args.publicKey
          );
        }, 2);
        if (fromTokenAllowance === BigInt(0)) {
          console.log("Router approval is revoked successfully");
          break;
        }
        await Wait(2);
      }
    }
  } else {
    let approveTxHash = "";
    await Forever(async () => {
      approveTxHash = await Approve(
        wallet,
        router,
        fromTokenContractAddress,
        fromTokenBalance.toString()
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
      while (true) {
        let fromTokenAllowance = BigInt(0);
        await Forever(async () => {
          console.log(`Refreshing router token allowance`);
          fromTokenAllowance = await router.GetApprovedAllowance(
            fromTokenContractAddress,
            Args.publicKey
          );
        }, 2);
        if (fromTokenAllowance >= fromTokenBalance) {
          console.log(`Router is now approved to spend the required tokens`);
          break;
        } else {
          await Wait(2);
        }
      }
    }
  }
};
