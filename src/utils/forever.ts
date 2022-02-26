import { Args } from "./flags";
import { Wait } from "./wait";
import isOnline from "is-online";

export const Forever = async (
  callback: any,
  retryIntervalSeconds: number,
  maxRetries: number = 100
): Promise<any> => {
  let result;
  let retries = 0;
  while (true) {
    try {
      result = await callback();
      break;
    } catch (error) {
      console.error(error);
    }
    try {
      const online = await isOnline();
      if (online) {
        retries += 1;
        if (retries <= maxRetries) {
          await Wait(retryIntervalSeconds);
        } else {
          console.error(`Max retries of ${maxRetries} exhausted`);
          process.exit(1);
        }
      } else {
        await Wait(retryIntervalSeconds);
      }
    } catch (error) {
      console.error(error);
      await Wait(retryIntervalSeconds);
    }
  }
  if (retries > 0 && Args.trace) {
    console.log(`Retries: ${retries}, ${callback.toString()}`);
  }
  return result;
};
