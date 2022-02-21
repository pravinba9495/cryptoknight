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
  console.log(
    `${
      amount === "0" ? "Revoking" : "Approving"
    } the router access to the tokens`
  );
  let approveTx = {};
  let approveTxGas = 0;

  while (true) {
    try {
      approveTx = await router.GetApproveTransactionData(
        tokenContractAddress,
        amount
      );
      break;
    } catch (error) {
      console.error(error);
    } finally {
      await Wait(2);
    }
  }
  while (true) {
    try {
      approveTxGas = await wallet.EstimateGas(approveTx);
      break;
    } catch (error) {
      console.error(error);
    } finally {
      await Wait(2);
    }
  }
  const approveTxWithGas = {
    ...approveTx,
    gas: approveTxGas,
  };
  let approveTxHash = "";
  while (true) {
    try {
      const nonce = await wallet.GetNonce();
      const signedApproveTxWithGasRaw = await wallet.SignTransaction({
        ...approveTxWithGas,
        nonce: nonce.toString(),
      });
      approveTxHash = await wallet.BroadcastRawTransaction(
        signedApproveTxWithGasRaw
      );
      break;
    } catch (error) {
      console.error(error);
    } finally {
      await Wait(2);
    }
  }
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
    } finally {
      await Wait(2);
    }
  }
  return Promise.reject(
    `${amount === "0" ? "Revoke" : "Approve"}  Transaction failed`
  );
};
