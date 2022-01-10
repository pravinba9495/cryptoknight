package oneinch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GetRouterLiquiditySources returns the list of liquidity sources that are available for routing in the 1inch Aggregation protocol
func GetRouterLiquiditySources(chainID uint64) ([]ProtocolImageDto, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ApiBaseUrl+"/"+ApiVersion+"/"+fmt.Sprint(chainID)+string(LiquiditySourcesEndpoint), nil)
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
		var dto *ProtocolsResponseDto
		if err := json.Unmarshal(body, &dto); err != nil {
			return nil, err
		}
		return dto.Protocols, nil
	} else {
		return nil, errors.New(resp.Status)
	}
}

// GetRouterSupportedTokens returns the list of tokens that are available for swap in the 1inch Aggregation protocol
func GetRouterSupportedTokens(chainID uint64) ([]TokenDto, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ApiBaseUrl+"/"+ApiVersion+"/"+fmt.Sprint(chainID)+string(TokensEndpoint), nil)
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
		var dto *TokensResponseDto
		if err := json.Unmarshal(body, &dto); err != nil {
			return nil, err
		}

		tokens := make([]TokenDto, 0)
		for _, token := range dto.Tokens {
			tokens = append(tokens, token)
		}
		return tokens, nil
	} else {
		return nil, errors.New(resp.Status)
	}
}

// GetRouterPresets returns the list of preset configurations for the 1inch router
func GetRouterPresets(chainID uint64) (*PresetsResponseDto, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ApiBaseUrl+"/"+ApiVersion+"/"+fmt.Sprint(chainID)+string(PresetsEndpoint), nil)
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
		var presets *PresetsResponseDto
		if err := json.Unmarshal(body, &presets); err != nil {
			return nil, err
		}
		return presets, nil
	} else {
		return nil, errors.New(resp.Status)
	}
}
