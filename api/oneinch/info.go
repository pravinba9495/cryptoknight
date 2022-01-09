package oneinch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// RouterProtocolsData schema
type RouterProtocolsData struct {
	Protocols []RouterProtocol `json:"protocols,omitempty"`
}

// RouterProtocol schema
type RouterProtocol struct {
	ID    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	Img   string `json:"img,omitempty"`
}

// RouterTokensData schema
type RouterTokensData struct {
	Tokens map[string]RouterToken `json:"tokens,omitempty"`
}

// RouterToken schema
type RouterToken struct {
	Symbol   string `json:"symbol,omitempty"`
	Name     string `json:"name,omitempty"`
	Address  string `json:"address,omitempty"`
	Decimals int    `json:"decimals"`
	LogoURI  string `json:"logoURI,omitempty"`
}

// RouterPresets schema
type RouterPresets struct {
	MaxResult []RouterPreset `json:"MAX_RESULT,omitempty"`
	LowestGas []RouterPreset `json:"LOWEST_GAS,omitempty"`
}

// RouterPreset schema
type RouterPreset struct {
	ComplexityLevel int `json:"complexityLevel,omitempty"`
	MainRouterParts int `json:"mainRouteParts,omitempty"`
	Parts           int `json:"parts,omitempty"`
	VirtualParts    int `json:"virtualParts,omitempty"`
}

// GetLiquiditySources returns the list of liquidity sources that are available for routing in the 1inch Aggregation protocol
func GetLiquiditySources(chainID int) ([]RouterProtocol, error) {
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
		var protocolsData *RouterProtocolsData
		if err := json.Unmarshal(body, &protocolsData); err != nil {
			return nil, err
		}
		return protocolsData.Protocols, nil
	} else {
		return nil, errors.New(resp.Status)
	}
}

// GetRouterSupportedTokens returns the list of tokens that are available for swap in the 1inch Aggregation protocol
func GetRouterSupportedTokens(chainID int) ([]RouterToken, error) {
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
		var tokensData *RouterTokensData
		if err := json.Unmarshal(body, &tokensData); err != nil {
			return nil, err
		}

		tokens := make([]RouterToken, 0)
		for _, token := range tokensData.Tokens {
			tokens = append(tokens, token)
		}
		return tokens, nil
	} else {
		return nil, errors.New(resp.Status)
	}
}

// GetRouterPresets returns the list of preset configurations for the 1inch router
func GetRouterPresets(chainID int) (*RouterPresets, error) {
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
		var presets *RouterPresets
		if err := json.Unmarshal(body, &presets); err != nil {
			return nil, err
		}
		return presets, nil
	} else {
		return nil, errors.New(resp.Status)
	}
}
