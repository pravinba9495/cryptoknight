package oneinch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// RouterAddress for the smart contract
type RouterAddress struct {
	Address string `json:"address,omitempty"`
}

// RouterAllowance schema
type RouterAllowance struct {
	Allowance string `json:"allowance,omitempty"`
}

// RouterTransactionData schema
type RouterTransactionData struct {
	Data     string `json:"data,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	To       string `json:"to,omitempty"`
	Value    string `json:"value,omitempty"`
}

// GetRouterAddressByChainID returns the address of the 1inch router that must be trusted to spend funds for the exchange
func GetRouterAddressByChainID(chainID int) (*RouterAddress, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ApiBaseUrl+"/"+ApiVersion+"/"+fmt.Sprint(chainID)+"/approve/spender", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var routerAddress *RouterAddress
		if err := json.Unmarshal(body, &routerAddress); err != nil {
			return nil, err
		}
		return routerAddress, nil
	} else {
		return nil, errors.New(resp.Status)
	}
}

// GenerateTransactionData generates transaction body to allow the exchange with the 1inch router
func GenerateTransactionData(chainID int, tokenAddress string, amount int64) (*RouterTransactionData, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ApiBaseUrl+"/"+ApiVersion+"/"+fmt.Sprint(chainID)+"/approve/transaction", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	query := req.URL.Query()

	query.Add("tokenAddress", tokenAddress)
	query.Add("amount", fmt.Sprint(amount))

	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var routerTransactionData *RouterTransactionData
		if err := json.Unmarshal(body, &routerTransactionData); err != nil {
			return nil, err
		}
		return routerTransactionData, nil
	} else {
		return nil, errors.New(resp.Status)
	}
}

// GetRouterAllowance returns the number of tokens that the 1inch router is allowed to spend
func GetRouterAllowance(chainID int, tokenAddress string, walletAddress string) (*RouterAllowance, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ApiBaseUrl+"/"+ApiVersion+"/"+fmt.Sprint(chainID)+"/approve/allowance", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	query := req.URL.Query()

	query.Add("tokenAddress", tokenAddress)
	query.Add("walletAddress", fmt.Sprint(walletAddress))

	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var routerAllowance *RouterAllowance
		if err := json.Unmarshal(bytes, &routerAllowance); err != nil {
			return nil, err
		}
		return routerAllowance, nil
	} else {
		return nil, errors.New(resp.Status)
	}
}
