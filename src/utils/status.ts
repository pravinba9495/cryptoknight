export const GetCurrentStatus = (
  stableTokenBalance: bigint,
  targetTokenBalance: bigint
): string => {
  let currentStatus = "UNKNOWN";
  if (stableTokenBalance !== BigInt(0) && targetTokenBalance === BigInt(0)) {
    currentStatus = "WAITING_TO_BUY";
  }
  if (stableTokenBalance === BigInt(0) && targetTokenBalance !== BigInt(0)) {
    currentStatus = "WAITING_TO_SELL";
  }
  return currentStatus;
};
