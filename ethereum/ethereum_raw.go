package ethereum

import (
	"alchemy-api/batch"
	"github.com/ybbus/jsonrpc/v3"
)

func (c EthClient) GetContractCodeRawBatch(requests jsonrpc.RPCRequests) (jsonrpc.RPCResponses, error) {
	return batch.DoBatchCall(c.client, requests)
}
