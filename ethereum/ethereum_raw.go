package ethereum

import (
	"context"
	"fmt"

	"alchemy-api/batch"
	"alchemy-api/utils"
	"github.com/ybbus/jsonrpc/v3"
)

const BaseApiUrl = "https://eth-mainnet.g.alchemy.com:443/v2/"

const (
	EthBlockNumber string = "eth_blockNumber"
	EthGetCode            = "eth_getCode"
	EthGetBalance         = "eth_getBalance"
	EthGetLogs            = "eth_getLogs"
	EthGasPrice           = "eth_gasPrice"
)

type ETHClientRaw struct {
	client jsonrpc.RPCClient
}

func getRpcClient(apiKey string) jsonrpc.RPCClient {
	url := fmt.Sprintf("%v%v", BaseApiUrl, apiKey)
	return jsonrpc.NewClient(url)
}

func NewETHClientRaw(apiKey string) *ETHClientRaw {
	return &ETHClientRaw{getRpcClient(apiKey)}
}

func (c ETHClientRaw) GetBlockNumberRaw() (*jsonrpc.RPCResponse, error) {
	return c.client.Call(context.Background(), EthBlockNumber)
}

func (c ETHClientRaw) GetContractCodeRaw(address string, blockNumberOpt ...string) (*jsonrpc.RPCResponse, error) {
	isValid := utils.CheckAddress(address)
	if !isValid {
		return nil, fmt.Errorf("invalid address %v", address)
	}

	blockNumber := Latest
	if len(blockNumberOpt) > 0 {
		blockNumber = blockNumberOpt[0]
	}

	return c.client.Call(context.Background(), EthGetCode, address, blockNumber)
}

func (c ETHClientRaw) GetBalance(address string, blockNumberOpt ...string) (*jsonrpc.RPCResponse, error) {
	isValid := utils.CheckAddress(address)
	if !isValid {
		return nil, fmt.Errorf("invalid address %v", address)
	}

	blockNumber := Latest
	if len(blockNumberOpt) > 0 {
		blockNumber = blockNumberOpt[0]
	}
	return c.client.Call(context.Background(), EthGetBalance, address, blockNumber)
}

func (c ETHClientRaw) GetLogs(request LogRequest) (*jsonrpc.RPCResponse, error) {
	for _, address := range request.Address {
		isValid := utils.CheckAddress(address)
		if !isValid {
			return nil, fmt.Errorf("invalid address %v", address)
		}
	}
	params := make([]interface{}, 1)
	params[0] = request
	return c.client.Call(context.Background(), EthGetLogs, params)
}

func (c ETHClientRaw) GetContractCodeBatchRaw(addresses []string, blockNumberOpt ...string) (jsonrpc.RPCResponses, error) {
	blockNumber := Latest
	if len(blockNumberOpt) > 0 {
		blockNumber = blockNumberOpt[0]
	}
	requests := make(jsonrpc.RPCRequests, len(addresses))
	for i, address := range addresses {
		isValid := utils.CheckAddress(address)
		if !isValid {
			return nil, fmt.Errorf("invalid address %v", address)
		}
		requests[i] = &jsonrpc.RPCRequest{Method: EthGetCode, Params: jsonrpc.Params(addresses[i], blockNumber), ID: i, JSONRPC: "2.0"}
	}
	return batch.DoBatchCall(c.client, requests)
}

func (c ETHClientRaw) GetBalanceBatch(addresses []string, blockNumberOpt ...string) (jsonrpc.RPCResponses, error) {
	blockNumber := Latest
	if len(blockNumberOpt) > 0 {
		blockNumber = blockNumberOpt[0]
	}
	requests := make(jsonrpc.RPCRequests, len(addresses))
	for i, address := range addresses {
		isValid := utils.CheckAddress(address)
		if !isValid {
			return nil, fmt.Errorf("invalid address %v", address)
		}
		requests[i] = &jsonrpc.RPCRequest{Method: EthGetBalance, Params: jsonrpc.Params(addresses[i], blockNumber), ID: i, JSONRPC: "2.0"}
	}
	return batch.DoBatchCall(c.client, requests)
}

func (c ETHClientRaw) GetGasPrice() (*jsonrpc.RPCResponse, error) {
	return nil, fmt.Errorf("not implemented")
}
