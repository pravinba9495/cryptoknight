import { Router } from "../api/oneinch";
import { Wallet } from "../api/wallet";
import { Wait } from "./wait";

/**
 * Swap method initates the swap token process on the router
 * @param wallet Wallet
 * @param router Router
 * @param params Swap parameters
 * @returns Promise<string>
 */
export const Swap = async (
  wallet: Wallet,
  router: Router,
  params: any
): Promise<string> => {
  let retries = 0;
  while (retries < 3) {
    try {
      let swapTxWithGas = {};
      while (true) {
        try {
          console.log(`Initiating swapping the tokens (Try: ${retries + 1})`);
          const swapTx = await router.GetSwapTransactionData(params);
          swapTxWithGas = {
            ...swapTx,
            gas: swapTx.gas + Math.ceil(0.25 * swapTx.gas),
          };
          break;
        } catch (error) {
          console.error(error);
          await Wait(5);
        }
      }
      const signedApproveTxWithGasRaw = await wallet.SignTransaction(
        swapTxWithGas
      );
      const swapTxHash = await router.BroadcastRawTransaction(
        signedApproveTxWithGasRaw
      );
      console.log(`Token Swap Transaction has been sent: ${swapTxHash}`);
      while (true) {
        console.log("Querying transaction status");
        const success = await wallet.GetTransactionReceipt(swapTxHash);
        if (success) {
          return swapTxHash;
        } else {
          return Promise.reject("Swap Transaction failed");
        }
      }
    } catch (error) {
      console.error(error);
    } finally {
      retries += 1;
    }
  }
  return Promise.reject("Swap Transaction failed");
};
