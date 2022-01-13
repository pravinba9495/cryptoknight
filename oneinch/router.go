package oneinch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-querystring/query"
	"github.com/pravinba9495/kryptonite/models"
)

type OneInchRouter models.Router

func New(chainID uint64) (*models.Router, error) {
	router := &OneInchRouter{}

	router.Vendor = "1INCH"
	router.ChainID = chainID

	addr, err := router.GetContractAddress()
	if err != nil {
		return nil, err
	}

	router.Address = addr

	tokens, err := router.GetSupportedTokens()
	if err != nil {
		return nil, err
	}

	router.SupportedTokens = tokens

	return (*models.Router)(router), nil
}

// GetHealthStatus checks the status of the 1inch router
func (r *OneInchRouter) GetHealthStatus() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ApiBaseUrl+"/"+ApiVersion+"/"+fmt.Sprint(r.ChainID)+string(HealthcheckEndpoint), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return errors.New(err.Error() + ":" + req.URL.String())
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status + ":" + req.URL.String())
	}
	return nil
}

// GetContractAddress returns the address of the 1inch router that must be trusted to spend funds for the exchange
func (r *OneInchRouter) GetContractAddress() (*common.Address, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ApiBaseUrl+"/"+ApiVersion+"/"+fmt.Sprint(r.ChainID)+string(SpenderEndpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error() + ":" + req.URL.String())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var dto *ApproveSpenderResponseDto
		if err := json.Unmarshal(body, &dto); err != nil {
			return nil, err
		}
		addr := common.HexToAddress(dto.Address)
		return &addr, nil
	} else {
		return nil, errors.New(resp.Status + ":" + req.URL.String())
	}
}

// GetSupportedTokens returns the list of tokens that are available for swap in the 1inch Aggregation protocol
func (r *OneInchRouter) GetSupportedTokens() ([]models.Token, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ApiBaseUrl+"/"+ApiVersion+"/"+fmt.Sprint(r.ChainID)+string(TokensEndpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(err.Error() + ":" + req.URL.String())
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

		tokens := make([]models.Token, 0)
		for _, token := range dto.Tokens {
			tokens = append(tokens, token)
		}
		r.SupportedTokens = tokens
		return tokens, nil
	} else {
		return nil, errors.New(resp.Status + ":" + req.URL.String())
	}
}

// GetApprovalTransactionData generates transaction body to allow the exchange with the 1inch router
func (r *OneInchRouter) GetApprovalTransactionData(params *ApproveCalldataParamsDto) (*ApproveCalldataResponseDto, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ApiBaseUrl+"/"+ApiVersion+"/"+fmt.Sprint(r.ChainID)+string(TransactionEndpoint), nil)
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
		return nil, errors.New(err.Error() + ":" + req.URL.String())
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
		return nil, errors.New(resp.Status + ":" + req.URL.String())
	}
}
