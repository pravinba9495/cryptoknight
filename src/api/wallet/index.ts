import Web3 from "web3";
import { Wait } from "../../utils/wait";
import { Web3Client, GetWeb3ContractClient } from "../web3";
import BN from "bn.js";

const ERC20Abi: any = [
  {
    constant: true,
    inputs: [{ name: "_owner", type: "address" }],
    name: "balanceOf",
    outputs: [{ name: "balance", type: "uint256" }],
    type: "function",
  },
];

export class Wallet {
  Address: string;
  MaskedAddress: string;
  Key: string;
  ChainID: number;

  constructor(address: string, key: string, chainId: number) {
    this.Address = address;
    this.MaskedAddress = address
      .split("")
      .map((c, idx) => {
        if (idx < 5 || idx > address.length - 6) {
          return c;
        } else {
          return "*";
        }
      })
      .join("");
    this.Key = key;
    this.ChainID = chainId;
  }

  async GetBalance(): Promise<BN> {
    const balance = await Web3Client.eth.getBalance(this.Address);
    return Web3.utils.toBN(balance);
  }

  async GetTokenBalance(tokenContractAddress: string): Promise<BN> {
    const contract = GetWeb3ContractClient(ERC20Abi, tokenContractAddress);
    const balance = await contract.methods.balanceOf(this.Address).call();
    return Web3.utils.toBN(balance);
  }

  async SignTransaction(transaction: any): Promise<any> {
    const tx: any = await Web3Client.eth.accounts.signTransaction(
      transaction,
      this.Key
    );
    const { rawTransaction, transactionHash } = tx;
    return { rawTransaction, transactionHash };
  }

  async EstimateGas(transaction: any): Promise<number> {
    const gas: number = await Web3Client.eth.estimateGas({
      ...transaction,
      from: this.Address,
    });
    return gas;
  }

  async GetTransactionReceipt(txHash: string): Promise<boolean> {
    let retries = 0;
    while (retries < 50) {
      const receipt = await Web3Client.eth.getTransactionReceipt(txHash);
      if (receipt != null) {
        console.log(`Fetched transaction receipt for ${txHash}`);
        return receipt.status;
      } else {
        retries += 1;
        console.log(
          `Retrying (${retries}) to fetch transaction receipt for ${txHash}`
        );
        await Wait(5);
      }
    }
    return false;
  }

  async GetNonce() {
    const nonce = await Web3Client.eth.getTransactionCount(
      this.Address,
      "pending"
    );
    return nonce;
  }
}
