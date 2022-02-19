let ENV: any = process.env;
export const Args = {
  publicKey: ENV.PUBLIC_KEY || "",
  privateKey: ENV.PRIVATE_KEY || "",
  mode: ENV.MODE || "MANUAL",
  chatId: ENV.CHAT_ID || "",
  chartInterval: ENV.CHART_INTERVAL || "1h",
  stableToken: ENV.STABLE_TOKEN || "USDC",
  targetToken: ENV.TARGET_TOKEN || "WMATIC",
  stableTokenTickerKraken: ENV.STABLE_TOKEN_TICKER_KRAKEN || "USDCUSD",
  targetTokenTickerKraken: ENV.TARGET_TOKEN_TICKER_KRAKEN || "MATICUSD",
  botToken: ENV.BOT_TOKEN || "",
  password: ENV.PASSWORD || "kryptonite",
  redisAddress: ENV.REDIS_ADDRESS || "",
  chainId: Number(ENV.CHAIN_ID || 137),
  minProfitPercent: Number(ENV.MIN_PROFIT_PERCENT || 1),
  slippagePercent: Number(ENV.SLIPPAGE_PERCENT || 1),
};
