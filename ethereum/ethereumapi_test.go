package ethereum

import (
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"alchemy-api/mocks"
	"github.com/ybbus/jsonrpc/v3"
)

const usdcCode = "0x608060405260043610"
const eoaCode = "0x"

func TestEthClient_BlockNumber(t *testing.T) {
	type fields struct {
		client jsonrpc.RPCClient
	}

	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{name: "Test Block", fields: fields{client: mocks.GetMockClient()}, want: "0x1234"},
		{name: "Test Error", fields: fields{client: mocks.GetMockClient(true)}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := EthClient{
				client: &ETHClientRaw{tt.fields.client},
			}
			got, err := c.GetBlockNumber()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
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
		name string
		args args
	}{
		{name: "Test_New", args: args{apiKey: "test"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.apiKey)
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
		{name: "Not Contract Result", fields: fields{client: mocks.GetMockClient()}, args: args{address: "0x549c660ce2b988f588769d6ad87be801695b2be3"}, want: eoaCode, wantErr: false},
		{name: "USDC Contract Result", fields: fields{client: mocks.GetMockClient()}, args: args{address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"}, want: usdcCode, wantErr: false},
		{name: "Not USDC", fields: fields{client: mocks.GetMockClient()}, args: args{address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB49"}, want: "0x", wantErr: false},
		{name: "Wrong address", fields: fields{client: mocks.GetMockClient()}, args: args{address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE360"}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := EthClient{
				client: &ETHClientRaw{tt.fields.client},
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
			fields: fields{client: mocks.GetMockClient()},
			args: args{
				address:        "0x549c660ce2b988f588769d6",
				blockNumberOpt: nil},
			wantErr: true,
		},
		{
			name:   "Positive Balance",
			fields: fields{client: mocks.GetMockClient()},
			args: args{
				address:        "0x549c660ce2b988f588769d6ad87be801695b2be3",
				blockNumberOpt: nil},
			want: big.NewInt(20066469208092992),
		},
		{
			name:   "Error Balance",
			fields: fields{client: mocks.GetMockClient()},
			args: args{
				address:        "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
				blockNumberOpt: nil},
			wantErr: true,
		},
		{
			name:   "Remote Error",
			fields: fields{client: mocks.GetMockClient()},
			args: args{
				address:        "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB47",
				blockNumberOpt: nil},
			wantErr: true,
		},
		{
			name:   "Zero Balance",
			fields: fields{client: mocks.GetMockClient()},
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
				client: &ETHClientRaw{tt.fields.client},
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

func TestEthClient_GetLogs(t *testing.T) {
	type fields struct {
		client jsonrpc.RPCClient
	}
	type args struct {
		request LogRequest
	}

	want := make([]LogsResponse, 1)
	err := json.Unmarshal([]byte(mocks.JS), &want)
	if err != nil {

	}

	address := make([]string, 1)
	address[0] = "0xb59f67a8bff5d8cd03f6ac17265c550ed8f33907"
	topics := make([]string, 3)
	topics[0] = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	topics[1] = "0x00000000000000000000000000b46c2526e227482e2ebb8f4c69e4674d262e75"
	topics[2] = "0x00000000000000000000000054a2d42a40f51259dedd1978f6c118a0f0eff078"

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *LogsResponse
		wantErr bool
	}{
		{name: "Get Logs", fields: fields{mocks.GetMockClient()}, args: args{request: NewLogRequest(address, "0x429d3b", Latest, topics...)}, want: &want[0]},
		{name: "Get Remote Error", fields: fields{mocks.GetMockClient()}, args: args{request: NewLogRequest(address, "0x429d3b", Pending, topics...)}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := EthClient{
				client: &ETHClientRaw{tt.fields.client},
			}
			got, err := c.GetLogs(tt.args.request)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetLogs() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else {
				if got == nil {
					t.Errorf("GetLogs() got nil")
				}
				if tt.want != nil && !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetLogs() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
