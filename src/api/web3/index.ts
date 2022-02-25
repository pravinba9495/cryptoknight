import Web3 from "web3";
import { GetRpcURLByChainID } from "../../networks";
import { Args } from "../../utils/flags";

export const Web3Client = new Web3(
  new Web3.providers.HttpProvider(GetRpcURLByChainID(Args.chainId))
);

export const GetWeb3ContractClient = (
  AbiItem: any,
  contractAddress: string
) => {
  const web3 = new Web3(
    new Web3.providers.HttpProvider(GetRpcURLByChainID(Args.chainId))
  );
  return new web3.eth.Contract(AbiItem, contractAddress);
};
