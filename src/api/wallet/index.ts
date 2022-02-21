import Web3 from "web3";
import { GetRpcURLByChainID } from "../../networks";
import { timeout } from "../../utils/timeout";
import { Wait } from "../../utils/wait";

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
    const web3 = new Web3(
      new Web3.providers.HttpProvider(GetRpcURLByChainID(this.ChainID))
    );
    const balance = await web3.eth.getBalance(this.Address);
    return BigInt(balance);
  }

  /**
   * GetTokenBalance returns the token balance at the given contract address
   * @param tokenContractAddress
   * @returns Promise<bigint> Token balance
   */
  async GetTokenBalance(tokenContractAddress: string): Promise<bigint> {
    const fn = async () => {
      const web3 = new Web3(
        new Web3.providers.HttpProvider(GetRpcURLByChainID(this.ChainID))
      );
      const contract = new web3.eth.Contract(ERC20Abi, tokenContractAddress);
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
    const web3 = new Web3(
      new Web3.providers.HttpProvider(GetRpcURLByChainID(this.ChainID))
    );
    const tx: any = await web3.eth.accounts.signTransaction(
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
    const web3 = new Web3(
      new Web3.providers.HttpProvider(GetRpcURLByChainID(this.ChainID))
    );
    const gas: number = await web3.eth.estimateGas({
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
    const web3 = new Web3(
      new Web3.providers.HttpProvider(GetRpcURLByChainID(this.ChainID))
    );
    while (true) {
      const receipt = await web3.eth.getTransactionReceipt(txHash);
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
    const web3 = new Web3(
      new Web3.providers.HttpProvider(GetRpcURLByChainID(this.ChainID))
    );
    const nonce = await web3.eth.getTransactionCount(this.Address, "latest");
    return nonce;
  }

  /**
   * BroadcastRawTransaction sends the signed transaction to the node
   * @param transaction Signed Transaction
   * @returns Promise<string> Hash of the transaction
   */
  async BroadcastRawTransaction(transaction: any) {
    const web3 = new Web3(
      new Web3.providers.HttpProvider(GetRpcURLByChainID(this.ChainID))
    );
    const { transactionHash } = await web3.eth.sendSignedTransaction(
      transaction
    );
    return transactionHash;
  }
}
