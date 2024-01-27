package ethereum

import (
	"fmt"

	"alchemy-api/batch"
	"alchemy-api/utils"
	"github.com/ybbus/jsonrpc/v3"
)

func (c EthClient) GetContractCodeBatch(addresses []string, blockNumberOpt ...string) (ContractCodeResponses, error) {
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
	values, err := batch.DoBatchCall(c.client, requests)
	if err != nil {
		return nil, err
	}
	response := make(ContractCodeResponses, len(addresses))
	for i, value := range values {
		address := addresses[i]
		code, codeError := extractValue(*value)
		response[i] = &ContractCodeResponse{
			Address: address,
			Code:    code,
			Error:   codeError,
		}
	}
	return response, nil
}

func extractValue(value jsonrpc.RPCResponse) (string, error) {
	responseError := value.Error
	if responseError == nil {
		code, err := value.GetString()
		if err != nil {
			return "", err
		}
		return code, nil
	}
	return "", fmt.Errorf(responseError.Error())
}
