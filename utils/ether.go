package utils

import (
	"fmt"
	"math/big"
	"regexp"

	"github.com/ybbus/jsonrpc/v3"
)

func CheckAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}

func GetString(response *jsonrpc.RPCResponse) (string, error) {
	responseError := response.Error
	if responseError == nil {
		code, err := response.GetString()
		if err != nil {
			return "", err
		}
		return code, nil
	}
	return "", fmt.Errorf(responseError.Error())
}

func GetBigInt(response *jsonrpc.RPCResponse) (*big.Int, error) {
	n := new(big.Int)
	if response.Error != nil {
		return nil, fmt.Errorf("remote Error: %v", response.Error.Error())
	}
	result, err := response.GetString()
	if err != nil {
		return nil, err
	}
	bigint, success := n.SetString(result, 0)
	if !success {
		return nil, fmt.Errorf("failed conversion for %v", result)
	}
	return bigint, err

}
