import { Gas } from "../api/gas";
import { Router } from "../api/oneinch";
import { Wallet } from "../api/wallet";
import { Args } from "./flags";
import { Forever } from "./forever";

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
  let approveTx: any = {};
  let signedApproveTxWithGasRaw: any = {};
  let approveTxGas = 0;
  let approveTxHash = "";
  let success = false;

  await Forever(async () => {
    approveTx = await router.GetApproveTransactionData(
      tokenContractAddress,
      amount
    );
  }, 2);

  await Forever(async () => {
    approveTxGas = await wallet.EstimateGas(approveTx);
  }, 2);

  const gasPrice = await Gas.GetGasPrice(Args.chainId);

  const approveTxWithGas = {
    ...approveTx,
    gasPrice: undefined,
    maxPriorityFeePerGas: gasPrice,
    maxFeePerGas: gasPrice,
    gas: approveTxGas,
  };

  await Forever(async () => {
    console.log(`Fetching nonce`);
    const nonce = await wallet.GetNonce();
    console.log(`Signing transaction`);
    const { rawTransaction, transactionHash } = await wallet.SignTransaction({
      ...approveTxWithGas,
      nonce: nonce.toString(),
    });
    signedApproveTxWithGasRaw = rawTransaction;
    approveTxHash = transactionHash;
  }, 2);

  console.log(
    `Attempting approve/reject transaction ${approveTxHash} with Gas: ${approveTxWithGas.gas}, MaxPriorityFeePerGas: ${approveTxWithGas.maxPriorityFeePerGas}, MaxFeePerGas: ${approveTxWithGas.maxFeePerGas} (wei)`
  );

  await Forever(
    async () => {
      console.log(`Broadcasting transaction ${approveTxHash}`);
      await router.BroadcastRawTransaction(signedApproveTxWithGasRaw);
    },
    2,
    3
  );

  console.log(
    `Token ${
      amount === "0" ? "Revoke" : "Approve"
    } Transaction has been sent: ${approveTxHash}`
  );

  await Forever(async () => {
    console.log("Querying transaction status");
    success = await wallet.GetTransactionReceipt(approveTxHash);
  }, 2);

  if (success) {
    console.log(
      `${amount === "0" ? "Revoke" : "Approve"}  Transaction Succeeded`
    );
    return approveTxHash;
  } else {
    console.error(
      `${amount === "0" ? "Revoke" : "Approve"}  Transaction Failed`
    );
    return "";
  }
};
