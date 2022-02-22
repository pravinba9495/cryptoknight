import Web3 from "web3";
import { GetRpcURLByChainID } from "../../networks";
import { Wait } from "../../utils/wait";
let web3Client: any;

export const InitWeb3Client = async (chainId: number) => {
  while (true) {
    try {
      const web3 = new Web3(
        new Web3.providers.HttpProvider(GetRpcURLByChainID(chainId))
      );
      await web3.eth.getGasPrice();
      web3Client = web3;
      return Promise.resolve(null);
    } catch (error) {
      console.error(error);
    } finally {
      await Wait(2);
    }
  }
};

export const GetWeb3ContractClient = (
  AbiItem: any,
  contractAddress: string
) => {
  return new (web3Client as Web3).eth.Contract(AbiItem, contractAddress);
};

export const GetWeb3Client = () => {
  return web3Client as Web3;
};
