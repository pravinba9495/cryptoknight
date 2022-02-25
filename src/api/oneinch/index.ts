import Axios from "axios";

class Token {
  id: string = "";
  name: string = "";
  decimals: number = 0;
  symbol: string = "";
  address: string = "";
  constructor(token: Token) {
    this.id = token.id || "";
    this.name = token.name || "";
    this.symbol = token.symbol || "";
    this.address = token.address || "";
    this.decimals = token.decimals || 0;
  }
}

export class Router {
  ChainID: number;

  constructor(chainId: number) {
    this.ChainID = chainId;
  }

  async GetSwapTransactionData(params: any): Promise<any> {
    return Axios.get<any>(`https://api.1inch.io/v4.0/${this.ChainID}/swap`, {
      params,
      timeout: 5000,
    })
      .then((response) => response.data)
      .then((response) => response.tx as any)
      .catch((error) => {
        if (error.response) {
          return Promise.reject(error.response.data);
        } else {
          return Promise.reject(new Error("Request failed"));
        }
      });
  }

  async GetQuote(params: any): Promise<any> {
    return Axios.get<any>(`https://api.1inch.io/v4.0/${this.ChainID}/quote`, {
      params,
      timeout: 5000,
    })
      .then((response) => response.data)
      .catch((error) => {
        if (error.response) {
          return Promise.reject(error.response.data);
        } else {
          return Promise.reject(new Error("Request failed"));
        }
      });
  }

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
          return Promise.reject(new Error("Request failed"));
        }
      });
  }

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
          return Promise.reject(new Error("Request failed"));
        }
      });
  }

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
          return Promise.reject(new Error("Request failed"));
        }
      });
  }

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
        timeout: 5000,
      }
    )
      .then((response) => response.data)
      .then((response: any) => BigInt(response.allowance))
      .catch((error) => {
        if (error.response) {
          return Promise.reject(error.response.data);
        } else {
          return Promise.reject(new Error("Request failed"));
        }
      });
  }

  async GetApproveTransactionData(
    tokenAddress: string,
    amount: string
  ): Promise<any> {
    return Axios.get<any>(
      `https://api.1inch.io/v4.0/${this.ChainID}/approve/transaction`,
      {
        params: {
          tokenAddress,
          amount: amount === "-1" ? undefined : amount,
        },
        timeout: 5000,
      }
    )
      .then((response) => response.data as any)
      .catch((error) => {
        if (error.response) {
          return Promise.reject(error.response.data);
        } else {
          return Promise.reject(new Error("Request failed"));
        }
      });
  }

  async BroadcastRawTransaction(rawTransaction: any): Promise<string> {
    return Axios.post<string>(
      `https://tx-gateway.1inch.io/v1.1/${this.ChainID}/broadcast`,
      {
        rawTransaction,
        timeout: 5000,
      }
    )
      .then((response) => response.data)
      .then((response: any) => (response.transactionHash || "") as string)
      .catch((error) => {
        if (error.response) {
          return Promise.reject(error.response.data);
        } else {
          return Promise.reject(new Error("Request failed"));
        }
      });
  }
}
