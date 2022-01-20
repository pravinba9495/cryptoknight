/**
 * Wait for the specified duration in seconds
 * @param durationSec Duration to wait in seconds
 * @returns
 */
export const Wait = async (durationSec: number): Promise<void> => {
  return new Promise((resolve, _) => {
    setTimeout(() => {
      resolve();
    }, durationSec * 1000);
  });
};
