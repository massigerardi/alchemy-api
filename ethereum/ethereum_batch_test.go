package ethereum

import (
	"fmt"
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

	addresses := make([]string, 2)
	addresses[0] = "0x549c660ce2b988f588769d6ad87be801695b2be3"
	addresses[1] = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"

	wrongAddresses := make([]string, 2)
	wrongAddresses[0] = "0x549c660ce2b988f588769d6ad87be80169"
	wrongAddresses[1] = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"

	blockNumber := make([]string, 1)
	blockNumber[0] = "latest"

	contractCodeResponses := make(ContractCodeResponses, 2)
	contractCodeResponses[0] = &ContractCodeResponse{
		Address: "0x549c660ce2b988f588769d6ad87be801695b2be3",
		Code:    mocks.EoaCode,
		Error:   nil,
	}
	contractCodeResponses[1] = &ContractCodeResponse{
		Address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		Code:    mocks.UsdcCode,
		Error:   nil,
	}

	errorResponseAddresses := make([]string, 2)
	errorResponseAddresses[0] = "0x549c660ce2b988f588769d6ad87be801695b2be1"
	errorResponseAddresses[1] = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"

	contractCodeErrorResponses := make(ContractCodeResponses, 2)
	contractCodeErrorResponses[0] = &ContractCodeResponse{
		Address: "0x549c660ce2b988f588769d6ad87be801695b2be1",
		Code:    "",
		Error:   fmt.Errorf("-123: wrong Response"),
	}
	contractCodeErrorResponses[1] = &ContractCodeResponse{
		Address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		Code:    mocks.UsdcCode,
		Error:   nil,
	}

	client := mocks.GetMockClient()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ContractCodeResponses
		wantErr bool
	}{
		{name: "Test Success", fields: fields{client: client}, args: args{addresses: addresses, blockNumberOpt: blockNumber}, want: contractCodeResponses},
		{name: "Test Wrong Address", fields: fields{client: client}, args: args{addresses: wrongAddresses, blockNumberOpt: blockNumber}, wantErr: true},
		{name: "Test Wrong Response", fields: fields{client: client}, args: args{addresses: errorResponseAddresses, blockNumberOpt: blockNumber}, want: contractCodeErrorResponses},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := EthClient{
				client: tt.fields.client,
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
