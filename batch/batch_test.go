package batch

import (
	"reflect"
	"testing"

	"alchemy-api/ethereum"
	"alchemy-api/mocks"
	"github.com/ybbus/jsonrpc/v3"
)

func TestBatchCall(t *testing.T) {
	type args struct {
		client   jsonrpc.RPCClient
		requests jsonrpc.RPCRequests
	}

	contractCodeRequests := jsonrpc.RPCRequests{
		&jsonrpc.RPCRequest{Method: ethereum.EthGetCode, Params: jsonrpc.Params("0x549c660ce2b988f588769d6ad87be801695b2be3", "latest"), ID: 1, JSONRPC: "2.0"},
		&jsonrpc.RPCRequest{Method: ethereum.EthGetCode, Params: jsonrpc.Params("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", "latest"), ID: 2, JSONRPC: "2.0"},
	}

	contractCodeResults := make([]string, 2)
	contractCodeResults[0] = mocks.EoaCode
	contractCodeResults[1] = mocks.UsdcCode

	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "Multiple GetCode", args: args{client: mocks.GetMockClient(), requests: contractCodeRequests}, want: contractCodeResults},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoBatchCall(tt.args.client, tt.args.requests)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoBatchCall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				if !reflect.DeepEqual(got[i], tt.want[i]) {
					t.Errorf("DoBatchCall() \ngot[%v] = %v\n want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
