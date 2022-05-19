import puppeteer from "puppeteer";
import { Args } from "./flags";
import { timeout } from "./timeout";
import { Wait } from "./wait";

let signal = "UNKNOWN";
let isPuppeteerReady = false;
let browser: puppeteer.Browser;

export const InitTradingViewTechnicals = async (
  ticker: string,
  interval: string
) => {
  try {
    signal = "UNKNOWN";
    browser = await puppeteer.launch({
      headless: true,
      defaultViewport: {
        width: 1920,
        height: 1080,
      },
      timeout: 5000,
      args: ["--no-sandbox", "--disable-setuid-sandbox"],
      pipe: true,
    });
    console.log("Browser Launch Successful");
    let isBuy = false;
    let isSell = false;
    let page = await browser.newPage();
    console.log("Page created");
    const fn = () => {
      return page.goto(
        `https://www.tradingview.com/symbols/${ticker}/technicals/`
      );
    };
    await timeout(fn(), 30000);
    console.log("Navigation To Page Successful");
    await page.waitForSelector(`button[id="${interval}"]`, {
      timeout: 10000,
    });
    while (true) {
      try {
        await page.click(`button[id="${interval}"]`);
        await Wait(2);
        const elements = await page.$$("." + Args.speedometerClass);
        if (elements.length !== 3) {
          throw "Puppeteer could not fetch trade signals from TradingView";
        }
        const promises: any[] = [];
        elements.forEach((element: any, index: number) => {
          if (index === 1) {
            promises.push(
              page.evaluate((e: any) => {
                return e.textContent;
              }, element)
            );
          }
        });
        const signals = await Promise.all(promises);
        isBuy =
          signals.filter((s) => s.includes("Buy")).length ===
          signals.length && signals.includes("Strong Buy");
        isSell = signals.filter((s) => s.includes("Sell")).length ===
        signals.length;
        signal = isBuy
          ? "STRONG BUY"
          : isSell
          ? "STRONG SELL"
          : `NEUTRAL`;
        isPuppeteerReady = true;
      } catch (error) {
        console.error(error);
        signal = "ERROR";
      }
      await Wait(5);
    }
  } catch (error) {
    signal = "ERROR";
    console.error(error);
  } finally {
    try {
      await browser.close();
    } catch (error) {
      console.error(error);
    }
  }
  InitTradingViewTechnicals(ticker, interval);
};

export const GetTradeSignal = (): string => {
  return signal;
};

export const IsPuppeteerReady = (): boolean => {
  return isPuppeteerReady;
};
