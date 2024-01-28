package ethereum

import (
	"alchemy-api/utils"
)

func (c EthClient) GetContractCodeBatch(addresses []string, blockNumberOpt ...string) (ContractCodeResponses, error) {
	responses, err := c.client.GetContractCodeBatchRaw(addresses, blockNumberOpt...)
	if err != nil {
		return nil, err
	}
	contractCodeResponses := make(ContractCodeResponses, len(addresses))
	for i, response := range responses {
		address := addresses[i]
		code, codeError := utils.GetString(response)
		contractCodeResponses[i] = &ContractCodeResponse{
			Address: address,
			Code:    code,
			Error:   codeError,
		}
	}
	return contractCodeResponses, nil
}

func (c EthClient) GetBalanceBatch(addresses []string, blockNumberOpt ...string) (BalanceResponses, error) {
	responses, err := c.client.GetBalanceBatch(addresses, blockNumberOpt...)
	if err != nil {
		return nil, err
	}
	balanceResponses := make(BalanceResponses, len(addresses))
	for i, response := range responses {
		address := addresses[i]
		amount, codeError := utils.GetBigInt(response)
		balanceResponses[i] = &BalanceResponse{
			Address: address,
			Amount:  *amount,
			Error:   codeError,
		}
	}
	return balanceResponses, nil
}
