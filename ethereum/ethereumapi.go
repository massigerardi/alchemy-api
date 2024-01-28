package ethereum

import (
	"fmt"
	"math/big"

	"github.com/massigerardi/alchemy-api/utils"
)

type EthClient struct {
	client *ETHClientRaw
}

func New(apiKey string) *EthClient {
	return &EthClient{client: NewETHClientRaw(apiKey)}
}

func (c EthClient) GetBlockNumber() (string, error) {
	response, err := c.client.GetBlockNumberRaw()
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

func (c EthClient) GetContractCode(address string, blockNumberOpt ...string) (string, error) {
	response, err := c.client.GetContractCodeRaw(address, blockNumberOpt...)
	if err != nil {
		return "", err
	}
	return utils.GetString(response)
}

func (c EthClient) GetBalance(address string, blockNumberOpt ...string) (*big.Int, error) {
	response, err := c.client.GetBalance(address, blockNumberOpt...)
	if err != nil {
		return nil, err
	}
	return utils.GetBigInt(response)
}

func (c EthClient) GetLogs(request LogRequest) (*LogsResponse, error) {
	response, err := c.client.GetLogs(request)
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
