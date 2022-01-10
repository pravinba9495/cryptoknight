package oneinch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-querystring/query"
)

// GetRouterAddressByChainID returns the address of the 1inch router that must be trusted to spend funds for the exchange
func GetRouterAddressByChainID(chainID uint64) (string, error) {
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
func GetRouterTransactionData(chainID uint64, params *ApproveCalldataParamsDto) (*ApproveCalldataResponseDto, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ApiBaseUrl+"/"+ApiVersion+"/"+fmt.Sprint(chainID)+string(TransactionEndpoint), nil)
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
func GetRouterAllowance(chainID uint64, params *ApproveAllowanceParamsDto) (string, error) {
	zero := "0"

	client := &http.Client{}
	req, err := http.NewRequest("GET", ApiBaseUrl+"/"+ApiVersion+"/"+fmt.Sprint(chainID)+string(AllowanceEndpoint), nil)
	if err != nil {
		return zero, err
	}
	req.Header.Add("Accept", "application/json")

	v, err := query.Values(params)
	if err != nil {
		return "", err
	}

	req.URL.RawQuery = v.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return zero, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return zero, err
		}
		var dto *ApproveAllowanceResponseDto
		if err := json.Unmarshal(bytes, &dto); err != nil {
			return zero, err
		}
		return dto.Allowance, nil
	} else {
		return zero, errors.New(resp.Status)
	}
}
