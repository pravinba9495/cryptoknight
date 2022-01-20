import yargs from "yargs/yargs";

/**
 * Args is a map of all the command line arguments received
 */
export const Args: any = yargs(process.argv.slice(2)).options({
  publicKey: { type: "string", default: "", demandOption: true },
  privateKey: { type: "string", default: "", demandOption: true },
  mode: { type: "string", default: "", demandOption: true },
  chatId: { type: "string", default: "", demandOption: true },
  stableToken: { type: "string", default: "USDC", demandOption: true },
  targetToken: { type: "string", default: "WETH", demandOption: true },
  botToken: { type: "string", default: "", demandOption: true },
  password: { type: "string", default: "kryptonite", demandOption: true },
  redisAddress: { type: "string", default: "", demandOption: true },
  chainId: { type: "number", default: 1, demandOption: true },
  profitPercent: { type: "number", default: 50, demandOption: true },
  stopLossPercent: { type: "number", default: 25, demandOption: true },
  slippagePercent: { type: "number", default: 1, demandOption: true },
}).argv;
