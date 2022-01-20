import { Router } from "../api/oneinch";
import { Wallet } from "../api/wallet";
import { Approve } from "./approve";
import { Args } from "./flags";
import { Revoke } from "./revoke";
import { Swap } from "./swap";
import { Wait } from "./wait";

/**
 * PrepareForSwap method prepares for a swap
 * @param router Router
 * @param wallet Wallet
 * @param fromTokenContractAddress Token Contract Address to swap from
 * @param fromTokenBalance  Amount of tokens to swap (from)
 * @param toTokenContractAddress Token Contract Address to swap to
 * @returns
 */
export const PrepareForSwap = async (
  router: Router,
  wallet: Wallet,
  fromTokenContractAddress: string,
  fromTokenBalance: number,
  toTokenContractAddress: string
): Promise<void> => {
  return new Promise(async (resolve, reject) => {
    try {
      let approveDone = false;
      let swapDone = false;
      let revokeDone = false;

      while (!swapDone || !approveDone || !revokeDone) {
        const fromTokenAllowance = await router.GetApprovedAllowance(
          fromTokenContractAddress,
          Args.publicKey
        );
        if (fromTokenAllowance >= fromTokenBalance) {
          console.log(
            `Router is approved to spend the required amount of tokens for a swap`
          );
          const params = {
            fromTokenAddress: fromTokenContractAddress,
            toTokenAddress: toTokenContractAddress,
            amount: fromTokenBalance,
            fromAddress: wallet.Address,
            slippage: Args.slippagePercent,
            disableEstimate: false,
            allowPartialFill: false,
          };
          await Swap(wallet, router, params);
          while (true) {
            console.log(`Refreshing new token balances`);
            const fromTokenBalance = await wallet.GetTokenBalance(
              fromTokenContractAddress
            );
            const toTokenBalance = await wallet.GetTokenBalance(
              toTokenContractAddress
            );
            if (fromTokenBalance === 0 && toTokenBalance !== 0) {
              swapDone = true;
              await Wait(5);
              break;
            }
            await Wait(5);
          }
          await Revoke(wallet, router, fromTokenContractAddress);
          while (true) {
            console.log(`Revoking router access to be on the safe side`);
            const fromTokenAllowance = await router.GetApprovedAllowance(
              fromTokenContractAddress,
              Args.publicKey
            );
            if (fromTokenAllowance === 0) {
              revokeDone = true;
              await Wait(5);
              break;
            }
            await Wait(5);
          }
        } else {
          await Approve(
            wallet,
            router,
            fromTokenContractAddress,
            fromTokenBalance
          );
          while (true) {
            console.log(`Refreshing router token allowance`);
            const fromTokenAllowance = await router.GetApprovedAllowance(
              fromTokenContractAddress,
              Args.publicKey
            );
            if (fromTokenAllowance >= fromTokenBalance) {
              approveDone = true;
              await Wait(5);
              break;
            }
            await Wait(5);
          }
        }
      }
      resolve();
    } catch (error) {
      reject(error);
    }
  });
};
