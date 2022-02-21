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
  let swapTxWithGas = {};
  while (true) {
    try {
      console.log(`Initiating swapping the tokens`);
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
  const signedSwapTxWithGasRaw = await wallet.SignTransaction(swapTxWithGas);

  let swapTxHash = "";
  while (true) {
    try {
      swapTxHash = await wallet.BroadcastRawTransaction(signedSwapTxWithGasRaw);
      break;
    } catch (error) {
      console.error(error);
    } finally {
      await Wait(2);
    }
  }
  console.log(`Token Swap Transaction has been sent: ${swapTxHash}`);
  while (true) {
    console.log("Querying transaction status");
    try {
      const success = await wallet.GetTransactionReceipt(swapTxHash);
      if (success) {
        return swapTxHash;
      } else {
        return Promise.reject("Swap Transaction failed");
      }
      break;
    } catch (error) {
      console.error(error);
    }
    await Wait(2);
  }
  return Promise.reject("Swap Transaction failed");
};
