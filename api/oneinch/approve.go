package oneinch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GetRouterAddressByChainID returns the address of the 1inch router that must be trusted to spend funds for the exchange
func GetRouterAddressByChainID(chainID int) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ApiBaseUrl+"/"+ApiVersion+"/"+fmt.Sprint(chainID)+string(SpenderEndpoint), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		var dto *ApproveSpenderResponseDto
		if err := json.Unmarshal(body, &dto); err != nil {
			return "", err
		}
		return dto.Address, nil
	} else {
		return "", errors.New(resp.Status)
	}
}

// GetRouterTransactionData generates transaction body to allow the exchange with the 1inch router
func GetRouterTransactionData(chainID int, tokenAddress string, amount int64) (*ApproveCalldataResponseDto, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ApiBaseUrl+"/"+ApiVersion+"/"+fmt.Sprint(chainID)+string(TransactionEndpoint), nil)
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
		var dto *ApproveCalldataResponseDto
		if err := json.Unmarshal(body, &dto); err != nil {
			return nil, err
		}
		return dto, nil
	} else {
		return nil, errors.New(resp.Status)
	}
}

// GetRouterAllowance returns the number of tokens that the 1inch router is allowed to spend
func GetRouterAllowance(chainID int, tokenAddress string, walletAddress string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ApiBaseUrl+"/"+ApiVersion+"/"+fmt.Sprint(chainID)+string(AllowanceEndpoint), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")

	query := req.URL.Query()

	query.Add("tokenAddress", tokenAddress)
	query.Add("walletAddress", fmt.Sprint(walletAddress))

	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		var dto *ApproveAllowanceResponseDto
		if err := json.Unmarshal(bytes, &dto); err != nil {
			return "", err
		}
		return dto.Allowance, nil
	} else {
		return "", errors.New(resp.Status)
	}
}
