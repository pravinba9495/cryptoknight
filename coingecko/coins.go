package coingecko

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// GetCoinsList retrieves the list of coins from CoinGecko
func GetCoinsList() ([]Coin, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/coins/list", nil)
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
		dto := make([]Coin, 0)
		if err := json.Unmarshal(body, &dto); err != nil {
			return nil, err
		}
		return dto, nil
	} else {
		return nil, errors.New(resp.Status + ":" + req.RequestURI)
	}
}

// GetCoinPrice retrieves the price of the coin
func GetCoinPrice(coinID string) (float64, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/simple/price?vs_currencies=usd&ids="+coinID, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}
		dto := make(map[string]map[string]float64)
		if err := json.Unmarshal(body, &dto); err != nil {
			return 0, err
		}
		return dto[coinID]["usd"], nil
	} else {
		return 0, errors.New(resp.Status + ":" + req.RequestURI)
	}
}
