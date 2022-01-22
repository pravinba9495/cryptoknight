import Axios from "axios";
import { string } from "yargs";
import { Token } from "../../models/token";

/**
 * Oneinch Router API
 */
export class Router {
  ChainID: number;

  constructor(chainId: number) {
    this.ChainID = chainId;
  }

  /**
   * GetSwapTransactionData returns the swap data required for the router
   * @param   params Swap parameters to pass to the router
   * @returns Promise<any>
   */
  async GetSwapTransactionData(params: any): Promise<any> {
    return Axios.get<any>(`https://api.1inch.io/v4.0/${this.ChainID}/swap`, {
      params,
    })
      .then((response) => response.data)
      .then((response) => response.tx as any)
      .catch((error) => {
        if (error.response) {
          return Promise.reject(error.response.data);
        } else {
          return Promise.reject(error.request);
        }
      });
  }

  /**
   * GetQuote gets the swap quote from 1inch
   * @param   params Quote parameters to pass to the router
   * @returns Promise<any>
   */
  async GetQuote(params: any): Promise<any> {
    return Axios.get<any>(`https://api.1inch.io/v4.0/${this.ChainID}/quote`, {
      params,
    })
      .then((response) => response.data)
      .catch((error) => {
        if (error.response) {
          return Promise.reject(error.response.data);
        } else {
          return Promise.reject(error.request);
        }
      });
  }

  /**
   * GetHealthStatus checks the health of the router
   * @returns Promise<void>
   */
  async GetHealthStatus(): Promise<boolean> {
    return Axios.get<boolean>(
      `https://api.1inch.io/v4.0/${this.ChainID}/healthcheck`
    )
      .then(() => {
        return Promise.resolve(true);
      })
      .catch((error) => {
        if (error.response) {
          return Promise.reject(error.response.data);
        } else {
          return Promise.reject(error.request);
        }
      });
  }

  /**
   * GetContractAddress returns the contract address of the router by the chain id
   * @returns Promise<string>
   */
  async GetContractAddress(): Promise<string> {
    return Axios.get<string>(
      `https://api.1inch.io/v4.0/${this.ChainID}/approve/spender`
    )
      .then((response) => response.data)
      .then((response: any) => response.address as string)
      .catch((error) => {
        if (error.response) {
          return Promise.reject(error.response.data);
        } else {
          return Promise.reject(error.request);
        }
      });
  }

  /**
   * GetSupportedTokens returns the contract address of the router by the chain id
   * @returns Promise<string>
   */
  async GetSupportedTokens(): Promise<Token[]> {
    return Axios.get<Token[]>(
      `https://api.1inch.io/v4.0/${this.ChainID}/tokens`
    )
      .then((response) => response.data)
      .then((response: any) => {
        let tokens: Token[] = [];
        for (let token of Object.keys(response.tokens as any)) {
          tokens.push(response.tokens[token] as Token);
        }
        return Promise.resolve(tokens.map((t) => new Token(t)));
      })
      .catch((error) => {
        if (error.response) {
          return Promise.reject(error.response.data);
        } else {
          return Promise.reject(error.request);
        }
      });
  }

  /**
   * GetApprovedAllowance returns the contract address of the router by the chain id
   * @param tokenAddress Token Contract Address
   * @param walletAddress Wallet Address
   * @returns Promise<bigint>
   */
  async GetApprovedAllowance(
    tokenAddress: string,
    walletAddress: string
  ): Promise<bigint> {
    return Axios.get<bigint>(
      `https://api.1inch.io/v4.0/${this.ChainID}/approve/allowance`,
      {
        params: {
          tokenAddress,
          walletAddress,
        },
      }
    )
      .then((response) => response.data)
      .then((response: any) => BigInt(response.allowance))
      .catch((error) => {
        if (error.response) {
          return Promise.reject(error.response.data);
        } else {
          return Promise.reject(error.request);
        }
      });
  }

  /**
   * GetApproveTransactionData gets the router data required to proceed with the token approval process
   * @param tokenAddress Token address to approve for the router
   * @param amount Amount to approve
   * @returns Promise<any> Transaction
   */
  async GetApproveTransactionData(
    tokenAddress: string,
    amount: string
  ): Promise<any> {
    return Axios.get<any>(
      `https://api.1inch.io/v4.0/${this.ChainID}/approve/transaction`,
      {
        params: {
          tokenAddress,
          amount,
        },
      }
    )
      .then((response) => response.data as any)
      .catch((error) => {
        if (error.response) {
          return Promise.reject(error.response.data);
        } else {
          return Promise.reject(error.request);
        }
      });
  }

  /**
   * BroadcastRawTransaction broadcasts transaction to the 1inch router
   * @params rawTransaction Raw Transaction
   * @returns Promise<string>
   */
  async BroadcastRawTransaction(rawTransaction: any): Promise<string> {
    return Axios.post<string>(
      `https://tx-gateway.1inch.io/v1.1/${this.ChainID}/broadcast`,
      {
        rawTransaction,
      }
    )
      .then((response) => response.data)
      .then((response: any) => response.transactionHash as string)
      .catch((error) => {
        if (error.response) {
          return Promise.reject(error.response.data);
        } else {
          return Promise.reject(error.request);
        }
      });
  }
}
