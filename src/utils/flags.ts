import yargs from "yargs/yargs";

/**
 * Args is a map of all the command line arguments received
 */
export const Args: any = yargs(process.argv.slice(2)).options({
  publicKey: { type: "string", demandOption: true },
  privateKey: { type: "string", demandOption: true },
  mode: { type: "string", default: "MANUAL", demandOption: true },
  chatId: { type: "string", demandOption: false },
  stableToken: { type: "string", default: "USDC", demandOption: true },
  targetToken: { type: "string", default: "WMATIC", demandOption: true },
  stableTokenTickerKraken: {
    type: "string",
    default: "USDCUSD",
    demandOption: true,
  },
  targetTokenTickerKraken: {
    type: "string",
    default: "MATICUSD",
    demandOption: true,
  },
  botToken: { type: "string", demandOption: true },
  password: { type: "string", default: "kryptonite", demandOption: true },
  redisAddress: { type: "string", demandOption: true },
  chainId: { type: "number", default: 137, demandOption: true },
  profitPercent: { type: "number", default: 50, demandOption: true },
  stopLossPercent: { type: "number", default: 25, demandOption: true },
  slippagePercent: { type: "number", default: 1, demandOption: true },
}).argv;
