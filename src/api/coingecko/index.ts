import Axios from "axios";
import { Token } from "../../models/token";

/**
 * CoinGecko API
 */
export class CoinGecko {
  /**
   * GetCoinPrice returns the coin price in USD
   * @param coinID Coin id from the coins list
   * @returns Promise<AxiosResponse<any, any>>
   */
  static GetCoinPrice = async (coinID: string): Promise<number> => {
    return Axios.get<number>("https://api.coingecko.com/api/v3/simple/price", {
      params: {
        vs_currencies: "usd",
        ids: coinID,
      },
    })
      .then((response) => response.data)
      .then((response: any) => {
        return Number(response[coinID]["usd"]);
      })
      .catch((error) => {
        if (error.response) {
          return Promise.reject(error.response.data);
        } else {
          return Promise.reject(error.request);
        }
      });
  };

  /**
   * GetCoinsList returns list of coins supported by Coingecko
   * @param coinID Coin id from the coins list
   * @returns Promise<AxiosResponse<any, any>>
   */
  static GetCoinsList = (): Promise<Token[]> => {
    return Axios.get<Token[]>("https://api.coingecko.com/api/v3/coins/list")
      .then((response) => response.data)
      .then((response: Token[]) => {
        return response;
      })
      .catch((error) => {
        if (error.response) {
          return Promise.reject(error.response.data);
        } else {
          return Promise.reject(error.request);
        }
      });
  };

  /**
   * GetCoinID returns the coin id from the symbol
   * @param symbol Symbol of the coin
   * @returns Promise<string>
   */
  static GetCoinID = async (symbol: string): Promise<string> => {
    try {
      const coins = await CoinGecko.GetCoinsList();
      let coinId = "";
      for (let coin of coins) {
        if (coin.symbol === symbol.toLowerCase()) {
          coinId = coin.id;
          break;
        }
      }
      return Promise.resolve(coinId);
    } catch (error) {
      return Promise.reject(error);
    }
  };
}
