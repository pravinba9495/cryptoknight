import Axios from "axios";
import Web3 from "web3";

export class Gas {
  static GetGasPrice = async (chainId: number): Promise<string> => {
    let URL = "";
    if (chainId === 1) {
      URL = "https://ethgasstation.info/json/ethgasAPI.json";
    } else if (chainId === 137) {
      URL = "https://gasstation-mainnet.matic.network/v2";
    } else {
      return Promise.reject("Unsupported chain");
    }
    return Axios.get<number>(URL, {
      timeout: 5000,
    })
      .then((response) => response.data)
      .then((response: any) => {
        if (chainId === 1) {
          return Web3.utils
            .toBN(
              `${Math.ceil((Number(response["fastest"]) * 1.5) / 10)}000000000`
            )
            .toString();
        } else if (chainId === 137) {
          return Web3.utils
            .toBN(
              `${Math.ceil(
                Number(response["fast"]["maxPriorityFee"]) * 1.5
              )}000000000`
            )
            .toString();
        } else {
          return Promise.reject("Unsupported chain");
        }
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
