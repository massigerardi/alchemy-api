package mocks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ybbus/jsonrpc/v3"
)

const JS = `[
        {
            "address": "0xb59f67a8bff5d8cd03f6ac17265c550ed8f33907",
            "blockHash": "0x8243343df08b9751f5ca0c5f8c9c0460d8a9b6351066fae0acbd4d3e776de8bb",
            "blockNumber": "0x429d3b",
            "data": "0x000000000000000000000000000000000000000000000000000000012a05f200",
            "logIndex": "0x56",
            "removed": false,
            "topics": [
                "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
                "0x00000000000000000000000000b46c2526e227482e2ebb8f4c69e4674d262e75",
                "0x00000000000000000000000054a2d42a40f51259dedd1978f6c118a0f0eff078"
            ],
            "transactionHash": "0xab059a62e22e230fe0f56d8555340a29b2e9532360368f810595453f6fdd213b",
            "transactionIndex": "0xac"
        }
    ]`

const UsdcCode = "0x608060405260043610"
const EoaCode = "0x"

func getMap(obj interface{}) (map[string]interface{}, error) {
	js, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(js, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type mockClient struct {
	jsonrpc.RPCClient
}

func (m mockClient) Call(_ context.Context, method string, params ...interface{}) (*jsonrpc.RPCResponse, error) {
	if method == "eth_blockNumber" {
		return &jsonrpc.RPCResponse{Result: "0x1234"}, nil
	}
	if method == "eth_getCode" {
		address := params[0].([]interface{})[0]
		switch address {
		case "0x549c660ce2b988f588769d6ad87be801695b2be3":
			return &jsonrpc.RPCResponse{Result: EoaCode}, nil
		case "0x549c660ce2b988f588769d6ad87be801695b2be1":
			return &jsonrpc.RPCResponse{Error: &jsonrpc.RPCError{
				Code:    -123,
				Message: "wrong Response",
				Data:    nil,
			}}, nil
		case "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB49":
			return &jsonrpc.RPCResponse{Result: EoaCode}, nil
		case "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48":
			return &jsonrpc.RPCResponse{Result: UsdcCode}, nil
		default:
			return &jsonrpc.RPCResponse{Result: ""}, nil
		}
	}
	if method == "eth_getBalance" {
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
	if method == "eth_getLogs" {
		received := params[0].([]interface{})
		request, _ := getMap(received[0])
		toBlock := request["toBlock"]
		switch toBlock {
		case "pending":
			return &jsonrpc.RPCResponse{Result: nil, Error: &jsonrpc.RPCError{
				Code:    -1234,
				Message: "Test Error",
				Data:    nil,
			}}, nil
		case "latest":

			result := make([]map[string]interface{}, 1)
			err := json.Unmarshal([]byte(JS), &result)
			if err != nil {
				return nil, err
			}
			return &jsonrpc.RPCResponse{
				Result: result,
				Error:  nil}, nil
		default:
			return &jsonrpc.RPCResponse{Result: "0x474a58f10b7140"}, nil
		}
	}
	return nil, fmt.Errorf("method not supported")
}

func (m mockClient) CallBatch(_ context.Context, requests jsonrpc.RPCRequests) (jsonrpc.RPCResponses, error) {
	responses := make(jsonrpc.RPCResponses, len(requests))
	for i, request := range requests {
		method := request.Method
		if method == "eth_getCode" {
			address := request.Params.([]interface{})[0].(string)
			switch address {
			case "0x549c660ce2b988f588769d6ad87be801695b2be3":
				responses[i] = &jsonrpc.RPCResponse{Result: EoaCode}
			case "0x549c660ce2b988f588769d6ad87be801695b2be1":
				responses[i] = &jsonrpc.RPCResponse{Error: &jsonrpc.RPCError{
					Code:    -123,
					Message: "wrong Response",
					Data:    nil,
				}}
			case "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB49":
				responses[i] = &jsonrpc.RPCResponse{Result: EoaCode}
			case "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48":
				responses[i] = &jsonrpc.RPCResponse{Result: UsdcCode}
			default:
				responses[i] = &jsonrpc.RPCResponse{Result: ""}
			}
		}
	}
	return responses, nil
}

func GetMockClient() jsonrpc.RPCClient {
	return mockClient{}
}
