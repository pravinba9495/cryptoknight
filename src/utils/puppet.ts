import puppeteer from "puppeteer";
import { timeout } from "./timeout";

export const GetTradeSignal = async (
  ticker: string,
  interval: string
): Promise<string> => {
  const fn = async () => {
    const browser = await puppeteer.launch({
      headless: true,
      defaultViewport: {
        width: 1920,
        height: 1080,
      },
      args: ["--no-sandbox", "--disable-setuid-sandbox"],
    });
    let isBuy = false;
    let isSell = false;

    try {
      const page = await browser.newPage();
      await page.goto(
        `https://www.tradingview.com/symbols/${ticker}/technicals/`,
        {
          timeout: 10000,
        }
      );
      await page.waitForSelector(`button[id="${interval}"]`, {
        timeout: 10000,
      });
      await page.click(`button[id="${interval}"]`);
      await page.waitForTimeout(10000);
      const elements = await page.$$(".speedometerSignal-DPgs-R4s");
      if (elements.length !== 3) {
        return Promise.reject(
          "Puppeteer could not fetch trade signals from TradingView"
        );
      }
      const promises: any[] = [];
      elements.forEach((element, index) => {
        if (index === 1) {
          promises.push(
            page.evaluate((e) => {
              return e.textContent;
            }, element)
          );
        }
      });
      const signals = await Promise.all(promises);
      isBuy =
        signals.filter((s) => s.includes("Buy")).length === signals.length;
      isSell =
        signals.filter((s) => s.includes("Sell")).length === signals.length;
      await browser.close();
    } catch (error) {
      console.error(error);
      await browser.close();
      return Promise.reject(error);
    }
    return isBuy ? "BUY" : isSell ? "SELL" : "HOLD";
  };
  const signal = await timeout(fn(), 60000);
  return signal;
};
