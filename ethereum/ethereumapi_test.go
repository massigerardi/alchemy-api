package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/ybbus/jsonrpc/v3"
)

const usdcCode = "0x608060405260043610"
const eoaCode = "0x"

const apiKey = "S0n286vv7IjOc-rbaBwu9zMsrfjd_CKs"

type mockClient struct {
	jsonrpc.RPCClient
}

func (m mockClient) Call(_ context.Context, method string, params ...interface{}) (*jsonrpc.RPCResponse, error) {
	if method == EthBlockNumber {
		return &jsonrpc.RPCResponse{Result: "0x1234"}, nil
	}
	if method == EthGetCode {
		address := params[0].([]interface{})[0]
		switch address {
		case "0x549c660ce2b988f588769d6ad87be801695b2be3":
			return &jsonrpc.RPCResponse{Result: eoaCode}, nil
		case "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB49":
			return &jsonrpc.RPCResponse{Result: eoaCode}, nil
		case "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48":
			return &jsonrpc.RPCResponse{Result: usdcCode}, nil
		default:
			return &jsonrpc.RPCResponse{Result: ""}, nil
		}
	}
	if method == EthGetBalance {
		address := params[0].([]interface{})[0]
		switch address {
		case "0x549c660ce2b988f588769d6ad87be801695b2be3":
			return &jsonrpc.RPCResponse{Result: "0x474a58f10b7140"}, nil
		case "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48":
			return &jsonrpc.RPCResponse{Result: ""}, nil
		default:
			return &jsonrpc.RPCResponse{Result: "0x0"}, nil
		}
	}
	return nil, fmt.Errorf("method not supported")
}

func (m mockClient) CallBatch(_ context.Context, _ jsonrpc.RPCRequests) (jsonrpc.RPCResponses, error) {
	return jsonrpc.RPCResponses{
		&jsonrpc.RPCResponse{JSONRPC: "2.0", ID: 1, Result: eoaCode, Error: nil},
		&jsonrpc.RPCResponse{JSONRPC: "2.0", ID: 2, Result: usdcCode, Error: nil},
	}, nil
}

func getMockClient() jsonrpc.RPCClient {
	return mockClient{}
}

func TestEthClient_BlockNumber(t *testing.T) {

	tests := []struct {
		name   string
		client jsonrpc.RPCClient
		want   string
	}{
		{"Test Block", getMockClient(), "0x1234"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := EthClient{
				client: tt.client,
			}
			got, err := c.GetBlockNumber()
			if err != nil {
				t.Errorf("Unexpected Error = %v", err)
			}
			if got != tt.want {
				t.Errorf("GetBlockNumber() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		apiKey string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Test_New", args: args{apiKey: apiKey}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.apiKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("New() got = %v", got)
			}
		})
	}
}

func TestEthClient_GetContractCode(t *testing.T) {
	type fields struct {
		client jsonrpc.RPCClient
	}
	type args struct {
		address        string
		blockNumberOpt []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{name: "Not Contract Result", fields: fields{client: getMockClient()}, args: args{address: "0x549c660ce2b988f588769d6ad87be801695b2be3"}, want: eoaCode, wantErr: false},
		{name: "USDC Contract Result", fields: fields{client: getMockClient()}, args: args{address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"}, want: usdcCode, wantErr: false},
		{name: "Not USDC", fields: fields{client: getMockClient()}, args: args{address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB49"}, want: "0x", wantErr: false},
		{name: "Wrong address", fields: fields{client: getMockClient()}, args: args{address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE360"}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := EthClient{
				client: tt.fields.client,
			}
			got, err := c.GetContractCode(tt.args.address, tt.args.blockNumberOpt...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetContractCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEthClient_BatchCall(t *testing.T) {
	type fields struct {
		client jsonrpc.RPCClient
	}
	type args struct {
		requests jsonrpc.RPCRequests
	}

	contractCodeRequests := jsonrpc.RPCRequests{
		&jsonrpc.RPCRequest{Method: EthGetCode, Params: jsonrpc.Params("0x549c660ce2b988f588769d6ad87be801695b2be3", "latest"), ID: 1, JSONRPC: "2.0"},
		&jsonrpc.RPCRequest{Method: EthGetCode, Params: jsonrpc.Params("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", "latest"), ID: 2, JSONRPC: "2.0"},
	}

	contractCodeResults := make([]string, 2)
	contractCodeResults[0] = eoaCode
	contractCodeResults[1] = usdcCode

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{name: "Multiple GetCode", fields: fields{client: getMockClient()}, args: args{requests: contractCodeRequests}, want: contractCodeResults},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := EthClient{
				client: tt.fields.client,
			}
			got, err := c.BatchCall(tt.args.requests)
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchCall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				if !reflect.DeepEqual(got[i], tt.want[i]) {
					t.Errorf("BatchCall() \ngot[%v] = %v\n want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestEthClient_GetBalance(t *testing.T) {
	type fields struct {
		client jsonrpc.RPCClient
	}
	type args struct {
		address        string
		blockNumberOpt []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *big.Int
		wantErr bool
	}{
		{
			name:   "Wrong Address",
			fields: fields{client: getMockClient()},
			args: args{
				address:        "0x549c660ce2b988f588769d6",
				blockNumberOpt: nil},
			wantErr: true,
		},
		{
			name:   "Positive Balance",
			fields: fields{client: getMockClient()},
			args: args{
				address:        "0x549c660ce2b988f588769d6ad87be801695b2be3",
				blockNumberOpt: nil},
			want: big.NewInt(20066469208092992),
		},
		{
			name:   "Error Balance",
			fields: fields{client: getMockClient()},
			args: args{
				address:        "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
				blockNumberOpt: nil},
			wantErr: true,
		},
		{
			name:   "Zero Balance",
			fields: fields{client: getMockClient()},
			args: args{
				address:        "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB49",
				blockNumberOpt: nil},
			want:    big.NewInt(0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := EthClient{
				client: tt.fields.client,
			}
			got, err := c.GetBalance(tt.args.address, tt.args.blockNumberOpt...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBalance() got = %v, want %v", got, tt.want)
			}
		})
	}
}
