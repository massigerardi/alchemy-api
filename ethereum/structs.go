package ethereum

import "math/big"

type BlockIdentifier string

const (
	Pending  string = "pending"
	Latest          = "latest"
	Safe            = "safe"
	Earliest        = "earliest"
)

type LogRequest struct {
	Address   []string `json:"address"`
	FromBlock string   `json:"fromBlock"`
	ToBlock   string   `json:"toBlock"`
	Topics    []string `json:"topics"`
}

func NewLogRequest(address []string, fromBlock string, toBlock string, topics ...string) LogRequest {
	return LogRequest{
		Address:   address,
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Topics:    topics,
	}
}

type ContractCodeResponses []*ContractCodeResponse
type ContractCodeResponse struct {
	Address string `json:"address"`
	Code    string `json:"code"`
	Error   error  `json:"error"`
}

type BalanceResponses []*BalanceResponse
type BalanceResponse struct {
	Address string  `json:"address"`
	Amount  big.Int `json:"amount"`
	Error   error   `json:"error"`
}

type LogsResponses []*LogsResponse

type LogsResponse struct {
	Address          string   `json:"address"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      string   `json:"blockNumber"`
	Data             string   `json:"data"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
	Topics           []string `json:"topics"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}
