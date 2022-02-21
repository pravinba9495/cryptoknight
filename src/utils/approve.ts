import { Router } from "../api/oneinch";
import { Wallet } from "../api/wallet";
import { Wait } from "./wait";

/**
 * Approve method will approve the required token access for the given router
 * @param wallet Wallet
 * @param router Router
 * @param tokenContractAddress Token Contract Address
 * @param amount Amount to approve
 * @returns Promise<string> Transaction Hash
 */
export const Approve = async (
  wallet: Wallet,
  router: Router,
  tokenContractAddress: string,
  amount: string
): Promise<string> => {
  let retries = 0;
  let nonce = 0;
  while (true) {
    try {
      nonce = await wallet.GetNonce();
      break;
    } catch (error) {
      console.error(error);
    }
    await Wait(2);
  }
  while (retries < 3) {
    try {
      console.log(
        `${
          amount === "0" ? "Revoking" : "Approving"
        } the router access to the tokens (Try: ${retries + 1})`
      );
      const approveTx = await router.GetApproveTransactionData(
        tokenContractAddress,
        amount
      );
      let approveTxGas = 0;
      while (true) {
        try {
          approveTxGas = await wallet.EstimateGas(approveTx);
          break;
        } catch (error) {
          console.error(error);
        }
        await Wait(2);
      }
      const approveTxWithGas = {
        ...approveTx,
        gas: approveTxGas,
      };
      const signedApproveTxWithGasRaw = await wallet.SignTransaction({
        ...approveTxWithGas,
        nonce: nonce.toString(),
      });
      const approveTxHash = await wallet.BroadcastRawTransaction(
        signedApproveTxWithGasRaw
      );
      nonce += 1;
      console.log(
        `Token ${
          amount === "0" ? "Revoke" : "Approve"
        } Transaction has been sent: ${approveTxHash}`
      );
      while (true) {
        console.log("Querying transaction status");
        try {
          const success = await wallet.GetTransactionReceipt(approveTxHash);
          if (success) {
            return approveTxHash;
          } else {
            return Promise.reject(
              `${amount === "0" ? "Revoke" : "Approve"} Transaction failed`
            );
          }
          break;
        } catch (error) {
          console.error(error);
        }
        await Wait(2);
      }
    } catch (error) {
      console.error(error);
    } finally {
      retries += 1;
    }
  }
  return Promise.reject(
    `${amount === "0" ? "Revoke" : "Approve"}  Transaction failed`
  );
};
