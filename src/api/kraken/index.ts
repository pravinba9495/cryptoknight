import Axios from "axios";

export class Kraken {
  static GetCoinPrice = async (pair: string): Promise<number> => {
    return Axios.get<number>("https://api.kraken.com/0/public/Ticker", {
      params: {
        pair,
      },
      timeout: 5000,
    })
      .then((response) => response.data)
      .then((response: any) => {
        if (response.error.length > 0) {
          return Promise.reject(response.error[0]);
        }
        return Number(response["result"][pair.toUpperCase()]["c"][0]);
      })
      .catch((error) => {
        if (error.response && error.response.data) {
          return Promise.reject(error.response.data);
        } else {
          return Promise.reject(new Error("Request failed"));
        }
      });
  };
}
