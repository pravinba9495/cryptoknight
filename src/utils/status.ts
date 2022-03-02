import BN from "bn.js";

export const GetCurrentStatus = (
  stableTokenBalance: BN,
  targetTokenBalance: BN
): string => {
  let currentStatus = "UNKNOWN";
  if (stableTokenBalance.gte(targetTokenBalance)) {
    currentStatus = "WAITING_TO_BUY";
  } else if (targetTokenBalance.gte(stableTokenBalance)) {
    currentStatus = "WAITING_TO_SELL";
  } else {
    currentStatus = "UNKNOWN";
  }
  return currentStatus;
};
