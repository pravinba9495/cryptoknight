import { Router } from "../api/oneinch";
import { Wallet } from "../api/wallet";

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
  console.log(`Initiating swapping the tokens`);
  const swapTx = await router.GetSwapTransactionData(params);
  const swapTxWithGas = {
    ...swapTx,
    gas: swapTx.gas + Math.ceil(0.25 * swapTx.gas),
  };
  const signedApproveTxWithGasRaw = await wallet.SignTransaction(swapTxWithGas);
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
      throw new Error("Swap Transaction failed");
    }
  }
  return swapTxHash;
};
