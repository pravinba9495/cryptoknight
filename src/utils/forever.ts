import { Wait } from "./wait";

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
    retries += 1;
    if (retries <= maxRetries) {
      await Wait(retryIntervalSeconds);
    } else {
      console.error(`Max retries of ${maxRetries} exhausted`);
      process.exit(1);
    }
  }
  if (retries > 0) {
    console.log(`Retries: ${retries}, ${callback.toString()}`);
  }
  return result;
};
