package oneinch

import (
	"errors"
	"fmt"
	"net/http"
)

// GetRouterHealthStatus checks the status of the 1inch router
func GetRouterHealthStatus(chainID uint64) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ApiBaseUrl+"/"+ApiVersion+"/"+fmt.Sprint(chainID)+string(HealthcheckEndpoint), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil
}
