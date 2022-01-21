import { Router } from "../api/oneinch";
import { Wallet } from "../api/wallet";

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
    `Approving the router to access the required amount of tokens for a swap`
  );
  const approveTx = await router.GetApproveTransactionData(
    tokenContractAddress,
    amount
  );
  const approveTxGas = await wallet.EstimateGas(approveTx);
  const approveTxWithGas = {
    ...approveTx,
    gas: approveTxGas,
  };
  const signedApproveTxWithGasRaw = await wallet.SignTransaction(
    approveTxWithGas
  );
  const approveTxHash = await router.BroadcastRawTransaction(
    signedApproveTxWithGasRaw
  );
  console.log(`Token Approval Transaction has been sent: ${approveTxHash}`);
  return approveTxHash;
};
