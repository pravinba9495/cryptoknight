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
	Protocols []RouterProtocol
}

// RouterProtocol schema
type RouterProtocol struct {
	ID    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	Img   string `json:"img,omitempty"`
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
