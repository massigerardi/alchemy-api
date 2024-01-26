package ethereum

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ybbus/jsonrpc/v3"
)

const (
	EthBlockNumber string = "eth_blockNumber"
	EthGetCode            = "eth_getCode"
	EthGetBalance         = "eth_getBalance"
	EthGetLogs            = "eth_getLogs"
)

const BaseApiUrl = "https://eth-mainnet.g.alchemy.com:443/v2/"

func getRpcClient(apiKey string) jsonrpc.RPCClient {
	url := strings.Builder{}
	url.WriteString(BaseApiUrl)
	url.WriteString(apiKey)
	return jsonrpc.NewClient(url.String())
}

func checkEther(s string) (bool, error) {
	s = strings.TrimPrefix(s, "0x")
	_, err := hex.DecodeString(s)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getStringResponse(client jsonrpc.RPCClient, method string, params ...interface{}) (string, error) {
	response, err := client.Call(context.Background(), method, params)
	if err != nil {
		return "", err
	}
	if response.Error != nil {
		return "", fmt.Errorf("remote Error: %v", response.Error.Error())
	}
	result, err := response.GetString()
	if err != nil {
		return "", err
	}
	return result, nil
}

type EthClient struct {
	client jsonrpc.RPCClient
}

func New(apiKey string) (*EthClient, error) {
	return &EthClient{getRpcClient(apiKey)}, nil
}

func (c EthClient) GetBlockNumber() (string, error) {
	return getStringResponse(c.client, EthBlockNumber)
}

func (c EthClient) GetContractCode(address string, blockNumberOpt ...string) (string, error) {
	_, err := checkEther(address)
	if err != nil {
		return "", err
	}

	blockNumber := "latest"
	if len(blockNumberOpt) > 0 {
		blockNumber = blockNumberOpt[0]
	}

	return getStringResponse(c.client, EthGetCode, address, blockNumber)
}

func (c EthClient) GetBalance(address string, blockNumberOpt ...string) (*big.Int, error) {
	_, err := checkEther(address)
	if err != nil {
		return nil, err
	}

	blockNumber := Latest
	if len(blockNumberOpt) > 0 {
		blockNumber = blockNumberOpt[0]
	}
	result, err := getStringResponse(c.client, EthGetBalance, address, blockNumber)
	if err != nil {
		return nil, err
	}
	n := new(big.Int)
	bigint, success := n.SetString(result, 0)
	if !success {
		return nil, fmt.Errorf("failed conversion for %v", result)
	}
	return bigint, err
}

func (c EthClient) GetLogs(request LogRequest) (*LogsResponse, error) {
	for _, address := range request.Address {
		_, err := checkEther(address)
		if err != nil {
			return nil, err
		}
	}
	params := make([]interface{}, 1)
	params[0] = request
	response, err := c.client.Call(context.Background(), EthGetLogs, params)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, fmt.Errorf("remote Error: %v", response.Error.Error())
	}
	var results []LogsResponse
	err = response.GetObject(&results)
	if err != nil {
		return nil, err
	}
	return &results[0], nil
}
