package oneinch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GetTokenBalances returns the token balances at the given wallet address
func GetTokenBalances(chainID uint64, walletAddress string) (map[TokenAddress]TokenBalance, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://balances.1inch.io/v1.1/"+fmt.Sprint(chainID)+"/balances/"+walletAddress, nil)
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
		var dto map[TokenAddress]TokenBalance
		if err := json.Unmarshal(body, &dto); err != nil {
			return nil, err
		}
		return dto, nil
	} else {
		return nil, errors.New(resp.Status)
	}
}
