package ethereum

import (
	"reflect"
	"testing"

	"alchemy-api/mocks"
	"github.com/ybbus/jsonrpc/v3"
)

func TestEthClient_GetContractCodeRawBatch(t *testing.T) {
	type fields struct {
		client jsonrpc.RPCClient
	}
	type args struct {
		requests jsonrpc.RPCRequests
	}

	contractCodeRequests := jsonrpc.RPCRequests{
		&jsonrpc.RPCRequest{Method: "eth_getCode", Params: jsonrpc.Params("0x549c660ce2b988f588769d6ad87be801695b2be3", "latest"), ID: 0, JSONRPC: "2.0"},
		&jsonrpc.RPCRequest{Method: "eth_getCode", Params: jsonrpc.Params("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", "latest"), ID: 1, JSONRPC: "2.0"},
	}

	contractCodeResponses := jsonrpc.RPCResponses{
		&jsonrpc.RPCResponse{Result: mocks.EoaCode, Error: nil, ID: 0},
		&jsonrpc.RPCResponse{Result: mocks.UsdcCode, Error: nil, ID: 1},
	}

	contractCodeErrorRequests := jsonrpc.RPCRequests{
		&jsonrpc.RPCRequest{Method: "eth_getCode", Params: jsonrpc.Params("0x549c660ce2b988f588769d6ad87be801695b2be1", "latest"), ID: 0, JSONRPC: "2.0"},
		&jsonrpc.RPCRequest{Method: "eth_getCode", Params: jsonrpc.Params("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", "latest"), ID: 1, JSONRPC: "2.0"},
	}

	contractCodeErrorResponses := jsonrpc.RPCResponses{
		&jsonrpc.RPCResponse{Result: nil, Error: &jsonrpc.RPCError{Code: -123, Message: "wrong Response", Data: nil}, ID: 0},
		&jsonrpc.RPCResponse{Result: mocks.UsdcCode, Error: nil, ID: 1},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    jsonrpc.RPCResponses
		wantErr bool
	}{
		{name: "Test Success", fields: fields{client: mocks.GetMockClient()}, args: args{requests: contractCodeRequests}, want: contractCodeResponses},
		{name: "Test Error", fields: fields{client: mocks.GetMockClient()}, args: args{requests: contractCodeErrorRequests}, want: contractCodeErrorResponses},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := EthClient{
				client: tt.fields.client,
			}
			got, err := c.GetContractCodeRawBatch(tt.args.requests)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractCodeRawBatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetContractCodeRawBatch() got = %v, want %v", got, tt.want)
			}
		})
	}
}
