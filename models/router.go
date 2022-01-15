package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/go-querystring/query"
	"github.com/pravinba9495/kryptonite/bot"
	"github.com/pravinba9495/kryptonite/constants"
)

// Router represents the model for the swap router
type Router struct {
	Address         *common.Address
	ChainID         uint64
	SupportedTokens []Token
}

// ApproveSpenderResponseDto schema
type ApproveSpenderResponseDto struct {
	// Address of the 1inch router that must be trusted to spend funds for the exchange
	Address string `json:"address,omitempty"`
}

// ApproveAllowanceParamsDto schema
type ApproveAllowanceParamsDto struct {
	TokenAddress  string `json:"tokenAddress,omitempty" url:"tokenAddress,omitempty"`
	WalletAddress string `json:"walletAddress,omitempty" url:"walletAddress,omitempty"`
}

// ApproveAllowanceResponseDto schema
type ApproveAllowanceResponseDto struct {
	Allowance string `json:"allowance,omitempty"`
}

// ApproveCalldataParamsDto schema
type ApproveCalldataParamsDto struct {
	TokenAddress string `json:"tokenAddress,omitempty" url:"tokenAddress,omitempty"`
	Amount       string `json:"amount,omitempty" url:"amount,omitempty"`
}

// ApproveCalldataResponseDto schema
type ApproveCalldataResponseDto struct {
	// The encoded data to call the approve method on the swapped token contract
	Data string `json:"data,omitempty"`
	// Gas price for fast transaction processing
	GasPrice string `json:"gasPrice,omitempty"`
	// Token address that will be allowed to exchange through 1inch router
	To string `json:"to,omitempty"`
	// Native token value in WEI (for approve is always 0)
	Value string `json:"value,omitempty"`
}

// TokensResponseDto schema
type TokensResponseDto struct {
	// List of supported tokens
	Tokens map[string]Token `json:"tokens,omitempty"`
}

// TransactionDto schema
type TransactionDto struct {
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Data     string `json:"data,omitempty"`
	Value    string `json:"value,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	Gas      uint64 `json:"gas,omitempty"`
}

// QuoteParamsDto schema
type QuoteParamsDto struct {
	FromTokenAddress string `json:"fromTokenAddress,omitempty" url:"fromTokenAddress,omitempty"`
	ToTokenAddress   string `json:"toTokenAddress,omitempty" url:"toTokenAddress,omitempty"`
	Amount           string `json:"amount,omitempty" url:"amount,omitempty"`
}

// QuoteResponseDto schema
type QuoteResponseDto struct {
	FromToken       Token  `json:"fromToken,omitempty"`
	ToToken         Token  `json:"toToken,omitempty"`
	FromTokenAmount string `json:"fromTokenAmount,omitempty"`
	ToTokenAmount   string `json:"toTokenAmount,omitempty"`
	EstimatedGas    uint64 `json:"estimatedGas,omitempty"`
}

// SwapParamsDto schema
type SwapParamsDto struct {
	FromTokenAddress string `json:"fromTokenAddress,omitempty" url:"fromTokenAddress,omitempty"`
	ToTokenAddress   string `json:"toTokenAddress,omitempty" url:"toTokenAddress,omitempty"`
	Amount           string `json:"amount,omitempty" url:"amount,omitempty"`
	FromAddress      string `json:"fromAddress,omitempty" url:"fromAddress,omitempty"`
	Slippage         string `json:"slippage,omitempty" url:"slippage,omitempty"`
	GasLimit         string `json:"gasLimit,omitempty" url:"gasLimit"`
	DisableEstimate  bool   `json:"disableEstimate,omitempty" url:"disableEstimate"`
}

// SwapResponseDto schema
type SwapResponseDto struct {
	FromToken       Token          `json:"fromToken,omitempty"`
	ToToken         Token          `json:"toToken,omitempty"`
	FromTokenAmount string         `json:"fromTokenAmount,omitempty"`
	ToTokenAmount   string         `json:"toTokenAmount,omitempty"`
	Tx              TransactionDto `json:"tx,omitempty"`
}

// New create a new router instance
func (r *Router) New(chainID uint64) (*Router, error) {
	router := &Router{
		ChainID: chainID,
	}
	addr, err := router.GetContractAddress()
	if err != nil {
		return nil, err
	}
	router.Address = addr
	tokens, err := router.GetSupportedTokens()
	if err != nil {
		return nil, err
	}
	router.SupportedTokens = tokens
	return router, nil
}

// GetHealthStatus checks the status of the 1inch router
func (r *Router) GetHealthStatus() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", constants.ApiBaseUrl+"/"+constants.ApiVersion+"/"+fmt.Sprint(r.ChainID)+string(constants.HealthcheckEndpoint), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status + ":" + req.URL.String() + ":" + string(body))
	}
	return nil
}

// GetContractAddress returns the address of the 1inch router that must be trusted to spend funds for the exchange
func (r *Router) GetContractAddress() (*common.Address, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", constants.ApiBaseUrl+"/"+constants.ApiVersion+"/"+fmt.Sprint(r.ChainID)+string(constants.SpenderEndpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		var dto *ApproveSpenderResponseDto
		if err := json.Unmarshal(body, &dto); err != nil {
			return nil, err
		}
		addr := common.HexToAddress(dto.Address)
		return &addr, nil
	} else {
		return nil, errors.New(resp.Status + ":" + req.URL.String() + ":" + string(body))
	}
}

// GetSupportedTokens returns the list of tokens that are available for swap in the 1inch Aggregation protocol
func (r *Router) GetSupportedTokens() ([]Token, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", constants.ApiBaseUrl+"/"+constants.ApiVersion+"/"+fmt.Sprint(r.ChainID)+string(constants.TokensEndpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		var dto *TokensResponseDto
		if err := json.Unmarshal(body, &dto); err != nil {
			return nil, err
		}
		tokens := make([]Token, 0)
		for _, token := range dto.Tokens {
			tokens = append(tokens, token)
		}
		r.SupportedTokens = tokens
		return tokens, nil
	} else {
		return nil, errors.New(resp.Status + ":" + req.URL.String() + ":" + string(body))
	}
}

// GetApprovedAllowance returns the number of tokens that the 1inch router is allowed to spend
func (r *Router) GetApprovedAllowance(chainID uint64, tokenAddress string, walletAddress string) (*ApproveAllowanceResponseDto, error) {
	params := &ApproveAllowanceParamsDto{
		TokenAddress:  tokenAddress,
		WalletAddress: walletAddress,
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", constants.ApiBaseUrl+"/"+constants.ApiVersion+"/"+fmt.Sprint(chainID)+string(constants.AllowanceEndpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = v.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		var dto *ApproveAllowanceResponseDto
		if err := json.Unmarshal(body, &dto); err != nil {
			return nil, err
		}
		return dto, nil
	} else {
		return nil, errors.New(resp.Status + ":" + req.URL.String() + ":" + string(body))
	}
}

// GetApproveAllowanceData generates transaction body to allow the exchange with the 1inch router
func (r *Router) GetApproveAllowanceData(tokenAddress string, amount string) (*ApproveCalldataResponseDto, error) {
	params := &ApproveCalldataParamsDto{
		TokenAddress: tokenAddress,
		Amount:       amount,
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", constants.ApiBaseUrl+"/"+constants.ApiVersion+"/"+fmt.Sprint(r.ChainID)+string(constants.TransactionEndpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = v.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		var dto *ApproveCalldataResponseDto
		if err := json.Unmarshal(body, &dto); err != nil {
			return nil, err
		}
		return dto, nil
	} else {
		return nil, errors.New(resp.Status + ":" + req.URL.String() + ":" + string(body))
	}
}

// GetQuote returns the best quote to exchange via 1inch router
func (r *Router) GetQuote(chainID uint64, params *QuoteParamsDto) (*QuoteResponseDto, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", constants.ApiBaseUrl+"/"+constants.ApiVersion+"/"+fmt.Sprint(chainID)+string(constants.QuoteEndpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = v.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		var dto *QuoteResponseDto
		if err := json.Unmarshal(body, &dto); err != nil {
			return nil, err
		}
		return dto, nil
	} else {
		return nil, errors.New(resp.Status + ":" + req.URL.String() + ":" + string(body))
	}
}

// GetSwapData returns the data for calling the 1inch router for exchange
func (r *Router) GetSwapData(chainID uint64, params *SwapParamsDto) (*SwapResponseDto, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", constants.ApiBaseUrl+"/"+constants.ApiVersion+"/"+fmt.Sprint(chainID)+string(constants.SwapEndpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = v.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		var dto *SwapResponseDto
		if err := json.Unmarshal(body, &dto); err != nil {
			return nil, err
		}
		return dto, nil
	} else {
		return nil, errors.New(resp.Status + ":" + req.URL.String() + ":" + string(body))
	}
}

// DoSwap prepares for a swap, waits for admin confirmation and proceeds for a swap
func (r *Router) DoSwap(w *Wallet, fromTokenContractAddress string, fromTokenBalance *big.Int, toTokenContractAddress string) error {

	accessRequestApproved := false

	for {
		allowanceDto, err := r.GetApprovedAllowance(r.ChainID, fromTokenContractAddress, w.Address.Hex())
		if err != nil {
			return err
		}
		allowance := new(big.Int)
		n, ok := allowance.SetString(allowanceDto.Allowance, 10)
		if !ok {
			err := "cannot convert string to big.Int"
			return errors.New(err)
		}
		allowance = n
		if allowance.Cmp(fromTokenBalance) > -1 {
			break
		} else {
			for !accessRequestApproved {
				bot.OutboundChannel <- "To proceed for a swap, the token contract must be allowed to access the required amount of tokens from your wallet.\n\nReply 'yes' to apprive this request\nReply 'no' to decline this request.\n\nThe approval request will auto expire in 30 seconds."
				bot.IsWaitingConfirmation = true
				go func() {
					time.Sleep(30 * time.Second)
					if bot.IsWaitingConfirmation {
						bot.ConfirmationChannel <- false
						bot.IsWaitingConfirmation = false
					}
				}()

				reply := <-bot.ConfirmationChannel
				bot.IsWaitingConfirmation = false
				if reply {
					accessRequestApproved = true
					dto, err := r.GetApproveAllowanceData(fromTokenContractAddress, fromTokenBalance.String())
					if err != nil {
						return err
					}
					toAddress := common.HexToAddress(fromTokenContractAddress)
					gasPrice := new(big.Int)
					gP, ok := gasPrice.SetString(dto.GasPrice, 10)
					if !ok {
						err := "cannot convert string to big.Int"
						return errors.New(err)
					}
					_, err = w.SendTransaction(&toAddress, &types.LegacyTx{
						GasPrice: gP,
						To:       &toAddress,
						Data:     []byte(dto.Data),
					})
					if err != nil {
						return errors.New("Token Access Approval Transaction Error: " + err.Error())
					}
				} else {
					return errors.New("REQUEST_EXPIRED_OR_DECLINED")
				}
			}
		}
		time.Sleep(5 * time.Second)
	}

	quoteParams := &QuoteParamsDto{
		FromTokenAddress: fromTokenContractAddress,
		ToTokenAddress:   toTokenContractAddress,
		Amount:           fromTokenBalance.String(),
	}
	quoteResDto, err := r.GetQuote(r.ChainID, quoteParams)
	if err != nil {
		return err
	}

	swapParams := &SwapParamsDto{
		FromTokenAddress: fromTokenContractAddress,
		ToTokenAddress:   toTokenContractAddress,
		Amount:           fromTokenBalance.String(),
		FromAddress:      w.Address.Hex(),
		Slippage:         "1",
		GasLimit:         fmt.Sprint(quoteResDto.EstimatedGas),
		DisableEstimate:  false,
	}
	swapResDto, err := r.GetSwapData(r.ChainID, swapParams)
	if err != nil {
		return err
	}

	fbalance1 := new(big.Float)
	fbalance1.SetString(swapResDto.FromTokenAmount)
	fromValue := new(big.Float).Quo(fbalance1, big.NewFloat(math.Pow10(int(swapResDto.FromToken.Decimals))))

	fbalance2 := new(big.Float)
	fbalance2.SetString(swapResDto.ToTokenAmount)
	toValue := new(big.Float).Quo(fbalance2, big.NewFloat(math.Pow10(int(swapResDto.ToToken.Decimals))))

	swapStr := fmt.Sprintf("%f %s => %f %s", fromValue, swapResDto.FromToken.Symbol, toValue, swapResDto.ToToken.Symbol)
	str := fmt.Sprintf("\n\nDate: %s\nOrder Type: %s\n\n%s\n\n", time.Now().Format(time.RFC822), "BUY", swapStr)

	bot.OutboundChannel <- str
	bot.OutboundChannel <- "Reply 'yes' to confirm\nReply 'no' to decline the swap.\n\nThe swap request will auto expire in 30 seconds."
	bot.IsWaitingConfirmation = true
	go func() {
		time.Sleep(30 * time.Second)
		if bot.IsWaitingConfirmation {
			bot.ConfirmationChannel <- false
			bot.IsWaitingConfirmation = false
		}
	}()

	reply := <-bot.ConfirmationChannel
	bot.IsWaitingConfirmation = false
	if reply {
		toAddress := common.HexToAddress(swapResDto.Tx.To)
		gasPrice := new(big.Int)
		gP, ok := gasPrice.SetString(swapResDto.Tx.GasPrice, 10)
		if !ok {
			err := "cannot convert string to big.Int"
			return errors.New(err)
		}
		_, err = w.SendTransaction(r.Address, &types.LegacyTx{
			GasPrice: gP,
			Gas:      swapResDto.Tx.Gas,
			To:       &toAddress,
			Data:     []byte(swapResDto.Tx.Data),
		})
		if err != nil {
			return errors.New("Swap Transaction Error: " + err.Error())
		}
		return nil
	} else {
		return errors.New("REQUEST_EXPIRED_OR_DECLINED")
	}
}
