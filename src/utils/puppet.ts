import puppeteer from "puppeteer";
import { Wait } from "./wait";

let signal = "HOLD";
let browser: puppeteer.Browser;

export const InitTradingViewTechnicals = async (
  ticker: string,
  interval: string
) => {
  while (true) {
    try {
      browser = await puppeteer.launch({
        headless: true,
        defaultViewport: {
          width: 1920,
          height: 1080,
        },
        args: ["--no-sandbox", "--disable-setuid-sandbox"],
      });
      let isBuy = false;
      let isSell = false;

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
      while (true) {
        try {
          await page.click(`button[id="${interval}"]`);
          await Wait(2);
          const elements = await page.$$(".speedometerSignal-DPgs-R4s");
          if (elements.length !== 3) {
            throw new Error(
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
            signals.filter((s) => s.includes("Strong Buy")).length ===
            signals.length;
          isSell =
            signals.filter((s) => s.includes("Strong Sell")).length ===
            signals.length;
          signal = isBuy ? "BUY" : isSell ? "SELL" : "HOLD";
        } catch (error) {
          console.error(error);
          signal = "HOLD";
          break;
        }
        await Wait(2);
      }
    } catch (error) {
      signal = "HOLD";
      console.error(error);
      try {
        await browser.close();
      } catch (error) {
        console.error(error);
      }
    }
    await Wait(2);
  }
};

export const GetTradeSignal = (): string => {
  return signal;
};
