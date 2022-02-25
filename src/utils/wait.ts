export const Wait = async (durationSec: number): Promise<void> => {
  return new Promise((resolve, _) => {
    setTimeout(() => {
      resolve();
    }, durationSec * 1000);
  });
};
