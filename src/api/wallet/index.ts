import { timeout } from "../../utils/timeout";
import { Wait } from "../../utils/wait";
import { GetWeb3Client, GetWeb3ContractClient } from "../web3";

/**
 * ERC20 minimal ABI
 */
const ERC20Abi: any = [
  {
    constant: true,
    inputs: [{ name: "_owner", type: "address" }],
    name: "balanceOf",
    outputs: [{ name: "balance", type: "uint256" }],
    type: "function",
  },
];

/**
 * Wallet API
 */
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

  /**
   * GetBalance returns the wallet balance
   * @returns Promise<bigint> Wallet balance
   */
  async GetBalance(): Promise<bigint> {
    const balance = await GetWeb3Client().eth.getBalance(this.Address);
    return BigInt(balance);
  }

  /**
   * GetTokenBalance returns the token balance at the given contract address
   * @param tokenContractAddress
   * @returns Promise<bigint> Token balance
   */
  async GetTokenBalance(tokenContractAddress: string): Promise<bigint> {
    const fn = async () => {
      const contract = GetWeb3ContractClient(ERC20Abi, tokenContractAddress);
      const balance = await contract.methods.balanceOf(this.Address).call();
      return BigInt(balance);
    };
    const balance = await timeout(fn(), 5000);
    return balance;
  }

  /**
   * SignTransaction signs the given transaction
   * @param transaction Transaction to sign
   * @returns Promise<string> Raw transaction
   */
  async SignTransaction(transaction: any): Promise<string> {
    const tx: any = await GetWeb3Client().eth.accounts.signTransaction(
      transaction,
      this.Key
    );
    return tx.rawTransaction || "";
  }

  /**
   * EstimateGas estimates the gas required for the transaction
   * @param transaction
   * @returns Estimated gas amount
   */
  async EstimateGas(transaction: any): Promise<number> {
    const gas: number = await GetWeb3Client().eth.estimateGas({
      ...transaction,
      from: this.Address,
    });
    return gas;
  }

  /**
   * GetTransactionReceipt gets the transaction receipt
   * @param txHash Transaction Hash
   * @returns Promise<boolean>
   */
  async GetTransactionReceipt(txHash: string): Promise<boolean> {
    while (true) {
      const receipt = await GetWeb3Client().eth.getTransactionReceipt(txHash);
      if (receipt != null) {
        return receipt.status;
      } else {
        await Wait(5);
      }
    }
  }

  /**
   * GetNonce returns a nonce
   * @returns Promise<number> Number of transactions on the wallet
   */
  async GetNonce() {
    const nonce = await GetWeb3Client().eth.getTransactionCount(
      this.Address,
      "pending"
    );
    return nonce;
  }

  /**
   * SuggestGasPrice suggests gas price based on the current state of the network
   * @returns Promise<string> Gas Price
   */
  async SuggestGasPrice() {
    const gasPrice = await GetWeb3Client().eth.getGasPrice();
    return gasPrice;
  }

  /**
   * BroadcastRawTransaction sends the signed transaction to the node
   * @param transaction Signed Transaction
   * @returns Promise<string> Hash of the transaction
   */
  async BroadcastRawTransaction(transaction: any) {
    const { transactionHash } = await GetWeb3Client().eth.sendSignedTransaction(
      transaction
    );
    return transactionHash;
  }
}
