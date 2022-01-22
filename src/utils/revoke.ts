import { Router } from "../api/oneinch";
import { Wallet } from "../api/wallet";

/**
 * Revoke method will revoke the token access for the given router
 * @param wallet Wallet
 * @param router Router
 * @param tokenContractAddress Token Contract Address
 * @returns Promise<string> Transaction Hash
 */
export const Revoke = async (
  wallet: Wallet,
  router: Router,
  tokenContractAddress: string
): Promise<string> => {
  console.log(`Revoking the router access to the tokens`);
  const revokeTx = await router.GetApproveTransactionData(
    tokenContractAddress,
    "0"
  );
  const revokeTxGas = await wallet.EstimateGas(revokeTx);
  const revokeTxWithGas = {
    ...revokeTx,
    gas: revokeTxGas,
  };
  const signedRevokedTxWithGasRaw = await wallet.SignTransaction(
    revokeTxWithGas
  );
  const revokeTxHash = await router.BroadcastRawTransaction(
    signedRevokedTxWithGasRaw
  );
  console.log(`Token Access Revoke Transaction has been sent: ${revokeTxHash}`);
  while (true) {
    console.log("Querying transaction status");
    const success = await wallet.GetTransactionReceipt(revokeTxHash);
    if (success) {
      return revokeTxHash;
    } else {
      throw "Revoke Transaction failed";
    }
  }
  return revokeTxHash;
};
