<p align="center">
  <img width="100" src="https://www.maxpixel.net/static/photo/1x/Growth-Green-Cartoon-Crystal-Green-Cartoon-307264.png" alt="Kryptonite Logo">
  <br/>
  <h1>kryptonite</h1>
</p>

[![Go Reference](https://pkg.go.dev/badge/github.com/pravinba9495/kryptonite.svg)](https://pkg.go.dev/github.com/pravinba9495/kryptonite) ![Go Report Card](https://goreportcard.com/badge/github.com/pravinba9495/kryptonite) ![Issues](https://img.shields.io/github/issues-raw/pravinba9495/kryptonite) ![License](https://img.shields.io/github/license/pravinba9495/kryptonite) ![Release](https://img.shields.io/github/v/release/pravinba9495/kryptonite?include_prereleases)

Automated crypto swapping bot, written in Go. Supports swapping on Ethereum, BSC, Polygon, Optimisim and Arbitrum networks using 1inch AggregatorV4 router. Under active development.

## Table of Contents
- [Introduction](#introduction)
- [Setup](#setup)
  - [Requirements](#requirements)
  - [Parameters](#parameters)
  - [Usage](#usage)
- [Documentation](#documentation)
- [Development](#development)
- [Maintainers](#maintainers)
- [License](#license)

## Introduction

## Setup

### Requirements

### Parameters

The following command line parameters are supported.

<div align="center">

<table border="1">
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
<td>--privateKey</td>
<td>Your wallet private key</td>
<td>string</td>
<td></td>
</tr>

<tr>
<td>--publicKey</td>
<td>Your wallet public address</td>
<td>string</td>
<td></td>
</tr>

<tr>
<td>--chainId</td>
<td>Chain to use. Allowed options: 1 (Ethereum), 10 (Optimism), 56 (Binance Smart Chain), 137 (Polygon/Matic), 42161 (Arbitrum)")</td>
<td>integer</td>
<td>1</td>
</tr>

<tr>
<td>--stableToken</td>
<td>Stable token (ERC20) to use. Example: USDC, USDT, DAI</td>
<td>string</td>
<td>USDC</td>
</tr>

<tr>
<td>--targetToken</td>
<td>Target ERC20 token to hold. Example: WETH, WMATIC, LINK.</td>
<td>string</td>
<td>WETH</td>
</tr>

<tr>
<td>--redisAddress</td>
<td>Redis server host. Example: 192.168.1.100:6379</td>
<td>string</td>
<td></td>
</tr>

<tr>
<td>--botToken</td>
<td>Telegram Bot Token used to send and receive transaction confirmations</td>
<td>string</td>
<td></td>
</tr>

<tr>
<td>--chatId</td>
<td>Your telegram chat id. You will receive this when you authorize yourself with the bot for the first time.</td>
<td>string</td>
<td></td>
</tr>

<tr>
<td>--password</td>
<td>Password to share with the bot to authorize yourself as the admin</td>
<td>string</td>
<td>kryptonite</td>
</tr>

<tr>
<td>--days</td>
<td>No. of days to use to calculate moving average</td>
<td>integer</td>
<td>30</td>
</tr>

<tr>
<td>--profitPercent</td>
<td>Profit percent at which the bot will execute a sell order</td>
<td>integer</td>
<td>50</td>
</tr>

<tr>
<td>--stopLossPercent</td>
<td>Loss percent at which the bot will execute a stop loss order</td>
<td>integer</td>
<td>25</td>
</tr>

</tbody>
</table>

</div>

### Usage
```shell
docker run pravinba9495/kryptonite:latest kryptonite \
                --privateKey=<PRIVATE_KEY> \
                --publicKey=<PUBLIC_ADDRESS> \
                --chainId=<CHAIN_ID> \
                --stableToken=<STABLE_TOKEN> \
                --targetToken=<TARGET_TOKEN> \
                --redisAddress=<REDIS_ADDRESS> \
                --botToken=<BOT_TOKEN> \
                --chatId=<CHAT_ID> \
                --password=<PASSWORD> \
                --days=<DAYS> \
                --profitPercent=<PROFIT_PERCENT> \
                --stopLossPercent=<STOP_LOSS_PERCENT>
```

## Documentation
Kryptonite documentation is hosted at [Read the docs](https://pkg.go.dev/github.com/pravinba9495/kryptonite).

## Development
Kryptonite is still under development. Contributions are always welcome!

## Maintainers
* [@pravinba9495](https://github.com/pravinba9495)
## License
MIT
