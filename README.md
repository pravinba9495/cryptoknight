<p align="center">
  <img width="100" src="https://www.maxpixel.net/static/photo/1x/Growth-Green-Cartoon-Crystal-Green-Cartoon-307264.png" alt="Kryptonite Logo">
  <br/>
  <h1>Kryptonite</h1>
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
