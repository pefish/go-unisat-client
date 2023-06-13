package btc_rpc_client

import (
	"fmt"
	go_http "github.com/pefish/go-http"
	go_logger "github.com/pefish/go-logger"
	"time"
)

type UnisatHttpClient struct {
	timeout time.Duration
	logger  go_logger.InterfaceLogger
	baseUrl string
}

func NewUnisatHttpClient(
	logger go_logger.InterfaceLogger,
	httpTimeout time.Duration,
) *UnisatHttpClient {
	return &UnisatHttpClient{
		timeout: httpTimeout,
		logger:  logger,
		baseUrl: "https://unisat.io/api",
	}
}

func (uhc *UnisatHttpClient) GetBtcBalance(address string) (string, error) {
	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  struct {
			Amount string `json:"amount"`
		} `json:"result"`
	}
	_, err := go_http.NewHttpRequester(go_http.WithLogger(go_logger.Logger), go_http.WithTimeout(uhc.timeout)).GetForStruct(go_http.RequestParam{
		Url: fmt.Sprintf("%s/v2/address/balance", uhc.baseUrl),
		Params: map[string]interface{}{
			"address": address,
		},
		Headers: map[string]interface{}{
			"X-Client":  "UniSat Wallet",
			"X-Version": true,
		},
	}, &result)
	if err != nil {
		return "", err
	}
	if result.Status == "0" {
		return "", fmt.Errorf("Get balance error - %s", result.Message)
	}
	return result.Result.Amount, nil
}

type Brc20BalanceResult struct {
	Ticker              string `json:"ticker"`
	TransferableBalance string `json:"transferableBalance"`
	AvailableBalance    string `json:"availableBalance"`
}

func (uhc *UnisatHttpClient) ListBrc20Balances(address string) (map[string]Brc20BalanceResult, error) {
	result := make(map[string]Brc20BalanceResult, 0)

	cursor := 0
	for {
		size := 100
		var httpResult struct {
			Status  string `json:"status"`
			Message string `json:"message"`
			Result  struct {
				List  []Brc20BalanceResult `json:"list"`
				Total uint64               `json:"total"`
			} `json:"result"`
		}
		_, err := go_http.NewHttpRequester(go_http.WithLogger(go_logger.Logger), go_http.WithTimeout(uhc.timeout)).GetForStruct(go_http.RequestParam{
			Url: fmt.Sprintf("%s/v3/brc20/tokens", uhc.baseUrl),
			Params: map[string]interface{}{
				"address": address,
				"cursor":  cursor,
				"size":    size,
			},
			Headers: map[string]interface{}{
				"X-Client":  "UniSat Wallet",
				"X-Version": true,
			},
		}, &httpResult)
		if err != nil {
			return nil, err
		}
		if httpResult.Status == "0" {
			return nil, fmt.Errorf("Get balance error - %s", httpResult.Message)
		}
		for _, b := range httpResult.Result.List {
			result[b.Ticker] = b
		}
		cursor += size
		if httpResult.Result.Total < uint64(size) {
			break
		}
	}
	return result, nil
}
