import Axios from "axios";
import Web3 from "web3";
import BN from "bn.js";

export class Gas {
  static GetGasPrice = async (chainId: number): Promise<bigint | BN> => {
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
        if (response.error.length > 0) {
          return Promise.reject(response.error[0]);
        }
        if (chainId === 1) {
          return Web3.utils.toBN(
            `${Math.ceil((Number(response["fastest"]) * 1.5) / 10)}000000000`
          );
        } else if (chainId === 137) {
          return Web3.utils.toBN(
            `${Math.ceil(
              Number(response["fast"]["maxPriorityFee"]) * 1.5
            )}000000000`
          );
        } else {
          return Promise.reject("Unsupported chain");
        }
      })
      .catch((error) => {
        if (error.response) {
          return Promise.reject(error.response.data);
        } else {
          return Promise.reject(new Error("Request failed"));
        }
      });
  };
}
