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
  let retries = 0;
  let nonce = await wallet.GetNonce();
  while (retries < 3) {
    try {
      console.log(
        `Revoking the router access to the tokens (Try: ${retries + 1})`
      );
      const revokeTx = await router.GetApproveTransactionData(
        tokenContractAddress,
        "0"
      );
      const revokeTxGas = await wallet.EstimateGas(revokeTx);
      const revokeTxWithGas = {
        ...revokeTx,
        gas: revokeTxGas,
      };
      const signedRevokedTxWithGasRaw = await wallet.SignTransaction({
        ...revokeTxWithGas,
        nonce: nonce.toString(),
      });
      const revokeTxHash = await wallet.BroadcastRawTransaction(
        signedRevokedTxWithGasRaw
      );
      nonce += 1;
      console.log(
        `Token Access Revoke Transaction has been sent: ${revokeTxHash}`
      );
      while (true) {
        console.log("Querying transaction status");
        const success = await wallet.GetTransactionReceipt(revokeTxHash);
        if (success) {
          return revokeTxHash;
        } else {
          return Promise.reject("Revoke Transaction failed");
        }
      }
    } catch (error) {
      console.error(error);
    } finally {
      retries += 1;
    }
  }
  return Promise.reject("Revoke Transaction failed");
};
