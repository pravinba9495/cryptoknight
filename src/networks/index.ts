const NETWORKS = {
  // For Ethereum
  Ethereum: "https://cloudflare-eth.com",

  // For Goerli Testnet
  Goerli: "https://goerli.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161",

  // For Binance Smart Chain (BSC)
  BinanceSmartChain: "https://bsc-dataseed.binance.org",

  // For Polygon
  Polygon: "https://polygon-rpc.com",

  // For Optimism
  Optimisim: "https://mainnet.optimism.io",

  // For Arbitrum
  Arbitrum: "https://arb1.arbitrum.io/rpc",
};

export const GetRpcURLByChainID = (chainID: number): string => {
  switch (chainID) {
    case 1:
      return NETWORKS.Ethereum;
    case 5:
      return NETWORKS.Goerli;
    case 56:
      return NETWORKS.BinanceSmartChain;
    case 137:
      return NETWORKS.Polygon;
    case 10:
      return NETWORKS.Optimisim;
    case 42161:
      return NETWORKS.Arbitrum;
    default:
      return NETWORKS.Ethereum;
  }
};
