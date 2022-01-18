<p align="center">
  <img width="100" src="https://www.maxpixel.net/static/photo/1x/Growth-Green-Cartoon-Crystal-Green-Cartoon-307264.png" alt="Kryptonite Logo">
  <br/>
  <h1>kryptonite</h1>
</p>

[![Go Reference](https://pkg.go.dev/badge/github.com/pravinba9495/kryptonite.svg)](https://pkg.go.dev/github.com/pravinba9495/kryptonite) ![Go Report Card](https://goreportcard.com/badge/github.com/pravinba9495/kryptonite) ![Issues](https://img.shields.io/github/issues-raw/pravinba9495/kryptonite) ![License](https://img.shields.io/github/license/pravinba9495/kryptonite) ![Release](https://img.shields.io/github/v/release/pravinba9495/kryptonite?include_prereleases)

Automated crypto swapping bot, written in Go. Supports multiple protocols like Uniswap, Curve, SushiSwap, Balancer, etc using 1inch AggregatorV4 router. Under active development.

**Not production ready yet**

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Benchmarks](#benchmarks)
- [Setup](#setup)
  - [Requirements](#requirements)
  - [Parameters](#parameters)
  - [Usage](#usage)
- [Instructions](#instructions)
- [Documentation](#documentation)
- [Development](#development)
- [Maintainers](#maintainers)
- [License](#license)

## Introduction
This project started as a hobby to figure out a way to keep an eye on the crypto market while juggling my day job. The Crypto market is highly volatile. Cryptocurrencies can fluctuate in price drastically within seconds you have your eyes off the screen. Unless you are a trader by profession, you cannot actively manage your portfolio, or make any meaningful and profitable moves or prevent a loss.

For people with assets on a hosted wallet or an exchange like Coinbase or Binance, you can set limit and stop loss orders. However, for people with crypto assets in a self-custodial wallet (not on exchanges like Coinbase or Binance), there are limited options.

With Kryptonite, you can set limit buy/sell, stop-loss orders and execute them automatically on your behalf, without manually confirming each time (although you can choose to confirm each one manually). Kryptonite can react to a market crash quickly and execute your stop loss order to protect your investment/profits. Importantly, Kryptonite strives to reduce your anxiety levels in an uncertain market.

Think of Kryptonite as an "Autopilot for Crypto Swapping". Kryptonite does not execute swaps on exchanges like Coinbase and Binance, simply because you cannot access your wallet's private key. Private keys are required by Kryptonite to execute an order on your behalf.

## Features
- Bundled with a webserver UI to view past swaps, current metrics and portfolio performance so far *(under development)*
- Executes Buy/Sell/Stop Loss orders manually using confirmations from you via the telegram bot
- Executes Buy/Sell/Stop Loss orders automatically based on set limits
- Supports multiple protocols like Uniswap, Curve, SushiSwap, Balancer, etc using the 1inch Aggregation protocol API

## Benchmarks

## Setup

### Requirements

### Parameters

The following command line parameters are supported.

<div align="center">

<table>
<thead>
<tr>
<th>Parameter</th>
<th>Description</th>
<th>Type</th>
<th>Default</th>
</tr>
</thead>
<tbody>

<tr>
<td>privateKey</td>
<td>Your wallet private key</td>
<td>string</td>
<td></td>
</tr>

<tr>
<td>publicKey</td>
<td>Your wallet public address</td>
<td>string</td>
<td></td>
</tr>

<tr>
<td>chainId</td>
<td>Chain to use. Allowed options: 1 (Ethereum), 10 (Optimism), 56 (Binance Smart Chain), 137 (Polygon/Matic), 42161 (Arbitrum)</td>
<td>integer</td>
<td>1</td>
</tr>

<tr>
<td>stableToken</td>
<td>Stable token (ERC20) to use. Example: USDC, USDT, DAI</td>
<td>string</td>
<td>USDC</td>
</tr>

<tr>
<td>targetToken</td>
<td>Target ERC20 token to hold. Example: WETH, WMATIC, LINK.</td>
<td>string</td>
<td>WETH</td>
</tr>

<tr>
<td>slippagePercent</td>
<td>Slipage tolerance during swap. Example: 1, 2, ..50</td>
<td>integer</td>
<td>1</td>
</tr>

<tr>
<td>redisAddress</td>
<td>Redis server host. Example: 192.168.1.100:6379</td>
<td>string</td>
<td></td>
</tr>

<tr>
<td>botToken</td>
<td>Telegram bot token used to send and receive transaction confirmations</td>
<td>string</td>
<td></td>
</tr>

<tr>
<td>chatId</td>
<td>Your telegram chat id. You will receive this when you authorize yourself with the bot for the first time.</td>
<td>string</td>
<td></td>
</tr>

<tr>
<td>password</td>
<td>Password to share with the bot to authorize yourself as the admin</td>
<td>string</td>
<td>kryptonite</td>
</tr>

<tr>
<td>profitPercent</td>
<td>Profit percent at which the bot will execute a sell order</td>
<td>integer</td>
<td>50</td>
</tr>

<tr>
<td>stopLossPercent</td>
<td>Loss percent at which the bot will execute a stop loss order</td>
<td>integer</td>
<td>25</td>
</tr>

<tr>
<td>mode</td>
<td>Mode of operation. Allowed options: MANUAL, AUTO</td>
<td>string</td>
<td>MANUAL</td>
</tr>

</tbody>
</table>

</div>

### Usage
```shell
docker run -d --name kryptonite pravinba9495/kryptonite:latest kryptonite \
                --privateKey=<PRIVATE_KEY> \
                --publicKey=<PUBLIC_ADDRESS> \
                --chainId=<CHAIN_ID> \
                --stableToken=<STABLE_TOKEN> \
                --targetToken=<TARGET_TOKEN> \
                --slippagePercent=<SLIPPAGE_PERCENT> \
                --redisAddress=<REDIS_ADDRESS> \
                --botToken=<BOT_TOKEN> \
                --chatId=<CHAT_ID> \
                --password=<PASSWORD> \
                --profitPercent=<PROFIT_PERCENT> \
                --stopLossPercent=<STOP_LOSS_PERCENT> \
                --mode=<MODE_OF_OPERATION>
```

## Disclaimer
By using Kryptonite, you understand and agree that there might be bugs in the code that could potentially lead to financial loss. You and only you are responsible for the use of this tool and any financial loss that it may cause. Nobody else shall be held accountable in this regard. Stop using this tool if you do not agree to this. Do conduct your own due diligence and consult your financial advisor before making any investment decisions.

## Documentation
Kryptonite documentation is hosted at [Read the docs](https://pkg.go.dev/github.com/pravinba9495/kryptonite).

## Development
Kryptonite is still under active development. Contributions are always welcome!

## Contributors
<a href="https://github.com/pravinba9495/kryptonite/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=pravinba9495/kryptonite" />
</a>

## License
MIT
