/**
 * Implements a timeouts
 * @param p A promise
 * @param timeout Timeout is ms
 * @returns Promise<any>
 */
export const timeout = (p: Promise<any>, timeout: number): Promise<any> => {
  let timer: NodeJS.Timeout;
  return Promise.race([
    p,
    new Promise(
      (_r, reject) =>
        (timer = setTimeout(() => {
          reject("Request timed out");
        }, timeout))
    ),
  ]).finally(() => clearTimeout(timer));
};
