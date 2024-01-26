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

	blockNumber := "latest"
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

func (c EthClient) BatchCall(requests jsonrpc.RPCRequests) ([]interface{}, error) {
	responses, err := c.client.CallBatch(context.Background(), requests)
	if err != nil {
		println(err)
		return nil, err
	}
	results := make([]interface{}, len(responses))
	for i, response := range responses {
		results[i] = response.Result
	}
	return results, nil
}
