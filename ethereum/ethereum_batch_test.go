package ethereum

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"alchemy-api/mocks"
	"github.com/ybbus/jsonrpc/v3"
)

func TestEthClient_GetContractCodeBatch(t *testing.T) {
	type fields struct {
		client jsonrpc.RPCClient
	}
	type args struct {
		addresses      []string
		blockNumberOpt []string
	}

	client := mocks.GetMockClient()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ContractCodeResponses
		wantErr bool
	}{
		{
			name:   "Test Success",
			fields: fields{client: client},
			args: args{
				addresses:      []string{"0x549c660ce2b988f588769d6ad87be801695b2be3", "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"},
				blockNumberOpt: []string{Latest}},
			want: ContractCodeResponses{
				&ContractCodeResponse{Address: "0x549c660ce2b988f588769d6ad87be801695b2be3", Code: mocks.EoaCode, Error: nil},
				&ContractCodeResponse{Address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", Code: mocks.UsdcCode, Error: nil},
			},
		},
		{
			name:   "Test Wrong Address",
			fields: fields{client: client},
			args: args{
				addresses:      []string{"0x549c660ce2b988f588769d6ad87be80169", "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"},
				blockNumberOpt: []string{Latest}},
			wantErr: true},
		{
			name:   "Test Wrong Response",
			fields: fields{client: client},
			args: args{
				addresses:      []string{"0x549c660ce2b988f588769d6ad87be801695b2be1", "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"},
				blockNumberOpt: []string{Latest}},
			want: ContractCodeResponses{
				&ContractCodeResponse{Address: "0x549c660ce2b988f588769d6ad87be801695b2be1", Code: "", Error: fmt.Errorf("-123: wrong Response")},
				&ContractCodeResponse{Address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", Code: mocks.UsdcCode, Error: nil},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := EthClient{
				client: &ETHClientRaw{tt.fields.client},
			}
			got, err := c.GetContractCodeBatch(tt.args.addresses, tt.args.blockNumberOpt...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractsCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetContractsCode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEthClient_GetBalanceBatch(t *testing.T) {
	type fields struct {
		client jsonrpc.RPCClient
	}
	type args struct {
		addresses      []string
		blockNumberOpt []string
	}

	client := mocks.GetMockClient()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    BalanceResponses
		wantErr bool
	}{
		{
			name:   "Test Success",
			fields: fields{client: client},
			args: args{
				addresses:      []string{"0x549c660ce2b988f588769d6ad87be801695b2be3", "0x558FA75074cc7cF045C764aEd47D37776Ea697d2"},
				blockNumberOpt: []string{Latest}},
			want: BalanceResponses{
				&BalanceResponse{Address: "0x549c660ce2b988f588769d6ad87be801695b2be3", Amount: *big.NewInt(20066469208092992), Error: nil},
				&BalanceResponse{Address: "0x558FA75074cc7cF045C764aEd47D37776Ea697d2", Amount: *big.NewInt(452046866901000), Error: nil},
			},
		},
		{
			name:   "Test Error",
			fields: fields{client: client},
			args: args{
				addresses:      []string{"0x549c660ce2b988f588769d6ad87be80169", "0x558FA75074cc7cF045C764aEd47D37776Ea697d2"},
				blockNumberOpt: []string{Latest}},
			wantErr: true,
		},
		{
			name:   "Test Remote Error",
			fields: fields{client: client},
			args: args{
				addresses:      []string{"0x549c660ce2b988f588769d6ad87be80169", "0x558FA75074cc7cF045C764aEd47D37776Ea697d1"},
				blockNumberOpt: []string{Latest}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := EthClient{
				client: &ETHClientRaw{tt.fields.client},
			}
			got, err := c.GetBalanceBatch(tt.args.addresses, tt.args.blockNumberOpt...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBalanceBatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBalanceBatch() got = %v, want %v", got, tt.want)
			}
		})
	}
}
