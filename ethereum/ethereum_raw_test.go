package ethereum

import (
	"reflect"
	"testing"

	"github.com/massigerardi/alchemy-api/mocks"

	"github.com/ybbus/jsonrpc/v3"
)

func TestEthClient_GetContractCodeBatchRaw(t *testing.T) {
	type fields struct {
		client jsonrpc.RPCClient
	}
	type args struct {
		addresses      []string
		blockNumberOpt []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    jsonrpc.RPCResponses
		wantErr bool
	}{
		{name: "Test Success", fields: fields{client: mocks.GetMockClient()}, args: args{
			addresses:      []string{"0x549c660ce2b988f588769d6ad87be801695b2be3", "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"},
			blockNumberOpt: nil,
		}, want: jsonrpc.RPCResponses{
			&jsonrpc.RPCResponse{Result: mocks.EoaCode, Error: nil, ID: 0},
			&jsonrpc.RPCResponse{Result: mocks.UsdcCode, Error: nil, ID: 1},
		}},
		{name: "Test Error", fields: fields{client: mocks.GetMockClient()}, args: args{
			addresses:      []string{"0x549c660ce2b988f588769d6ad87be801695b2be1", "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"},
			blockNumberOpt: nil,
		}, want: jsonrpc.RPCResponses{
			&jsonrpc.RPCResponse{Result: nil, Error: &jsonrpc.RPCError{Code: -123, Message: "wrong Response", Data: nil}, ID: 0},
			&jsonrpc.RPCResponse{Result: mocks.UsdcCode, Error: nil, ID: 1},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ETHClientRaw{
				client: tt.fields.client,
			}
			got, err := c.GetContractCodeBatchRaw(tt.args.addresses, tt.args.blockNumberOpt...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractCodeBatchRaw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetContractCodeBatchRaw() got = %v, want %v", got, tt.want)
			}
		})
	}
}
