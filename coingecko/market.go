package coingecko

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// GetMarketChartByCoin returns the market chart data from the Coingecko API
func GetMarketChartByCoin(coinID string, days uint64) (*MarketChartResponseDto, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/coins/"+strings.ToLower(coinID)+"/market_chart?days="+fmt.Sprint(days)+"&vs_currency=usd&interval=daily", nil)
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
		var dto *MarketChartResponseDto
		if err := json.Unmarshal(body, &dto); err != nil {
			return nil, err
		}
		return dto, nil
	} else {
		return nil, errors.New(resp.Status + ":" + req.RequestURI)
	}
}
