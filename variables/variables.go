package variables

import (
	"errors"
	"flag"
)

// Shared variables
var (
	Address, PrivateKey, StableToken, TargetToken, BotToken, Password, RedisAddress, Mode, ChatID, CurrentStatus, Verdict string
	ChainID, ProfitPercent, StopLossPercent, Slippage                                                                     uint64
)

func ParseFlags() error {
	// Parse input arguments
	flag.StringVar(&Address, "publicKey", "", "Your wallet public address")
	flag.StringVar(&PrivateKey, "privateKey", "", "Your wallet private key")
	flag.StringVar(&Mode, "mode", "MANUAL", "Mode of operation. Allowed options: 'MANUAL' (requires you to authorize a swap via telgram bot), 'AUTO' (does not require any confirmation)")
	flag.StringVar(&ChatID, "chatId", "", "Your telegram chat id. You will receive this when you authorize yourself with the bot for the first time")
	flag.StringVar(&StableToken, "stableToken", "USDC", "Stable token (ERC20) to use. Example: USDC, USDT, DAI")
	flag.StringVar(&TargetToken, "targetToken", "WETH", "Target ERC20 token to hold. Example: WETH, WMATIC, LINK.")
	flag.StringVar(&BotToken, "botToken", "", "Telegram bot token used to send and receive transaction confirmations")
	flag.StringVar(&Password, "password", "kryptonite", "Password to share with the bot to authorize yourself as the admin")
	flag.StringVar(&RedisAddress, "redisAddress", "", "Redis server host. Example: 192.168.1.100:6379")
	flag.Uint64Var(&ChainID, "chainId", 1, "Chain to use. Allowed options: 1 (Ethereum), 10 (Optimism), 56 (Binance Smart Chain), 137 (Polygon/Matic), 42161 (Arbitrum)")
	flag.Uint64Var(&ProfitPercent, "profitPercent", 50, "Profit percent at which the bot will execute a sell order")
	flag.Uint64Var(&StopLossPercent, "stopLossPercent", 25, "Loss percent at which the bot will execute a stop loss order")
	flag.Uint64Var(&Slippage, "slippagePercent", 1, "Maximum slippage allowed during a swap")
	flag.Parse()

	if BotToken == "" {
		return errors.New("botToken is not provided")
	}

	if PrivateKey == "" {
		return errors.New("privateKey is not provided")
	}

	if RedisAddress == "" {
		return errors.New("redis address is not provided")
	}

	if (Mode != "MANUAL") && (Mode != "AUTO") {
		return errors.New("unsupported mode")
	}

	return nil
}
