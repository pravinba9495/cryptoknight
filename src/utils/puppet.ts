import puppeteer from "puppeteer";

export const GetTradeSignal = async (ticker: string) => {
  const browser = await puppeteer.launch({
    headless: true,
    defaultViewport: {
      width: 1920,
      height: 1080,
    },
    args: ["--no-sandbox", "--disable-setuid-sandbox"],
  });
  const page = await browser.newPage();
  await page.goto(`https://www.tradingview.com/symbols/${ticker}/technicals/`);
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
  const isBuy =
    signals.filter((s) => s.includes("Buy")).length === signals.length;
  const isSell =
    signals.filter((s) => s.includes("Sell")).length === signals.length;
  await browser.close();
  return isBuy ? "BUY" : isSell ? "SELL" : "HOLD";
};
