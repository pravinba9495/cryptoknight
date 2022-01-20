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
  console.log(`Initiating swapping the tokens`);
  const swapTx = await router.GetSwapTransactionData(params);
  const swapTxGas = await wallet.EstimateGas(swapTx);
  const swapTxWithGas = {
    ...swapTx,
    gas: swapTxGas,
  };
  const signedApproveTxWithGasRaw = await wallet.SignTransaction(swapTxWithGas);
  const swapTxHash = await router.BroadcastRawTransaction(
    signedApproveTxWithGasRaw
  );
  console.log(`Token Swap Transaction has been sent: ${swapTxHash}`);
  return swapTxHash;
};
