package batch

import (
	"context"

	"github.com/ybbus/jsonrpc/v3"
)

func DoBatchCall(client jsonrpc.RPCClient, requests jsonrpc.RPCRequests) (jsonrpc.RPCResponses, error) {
	responses, err := client.CallBatch(context.Background(), requests)
	if err != nil {
		println(err)
		return nil, err
	}
	return responses, nil
}
