package batch

import (
	"context"

	"github.com/ybbus/jsonrpc/v3"
)

func DoBatchCall(client jsonrpc.RPCClient, requests jsonrpc.RPCRequests) ([]interface{}, error) {
	responses, err := client.CallBatch(context.Background(), requests)
	if err != nil {
		println(err)
		return nil, err
	}
	results := make([]interface{}, len(responses))
	for i, response := range responses {
		if response.Error != nil {
			results[i] = response.Error
		} else {
			results[i] = response.Result
		}
	}
	return results, nil
}
