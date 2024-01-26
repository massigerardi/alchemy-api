package ethereum

import (
	"reflect"
	"testing"
)

func TestNewLogRequest(t *testing.T) {
	type args struct {
		address   []string
		fromBlock string
		toBlock   string
		topics    []string
	}
	address := make([]string, 1)
	address[0] = "0xb59f67a8bff5d8cd03f6ac17265c550ed8f33907"
	topics := make([]string, 3)
	topics[0] = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	topics[1] = "0x00000000000000000000000000b46c2526e227482e2ebb8f4c69e4674d262e75"
	topics[2] = "0x00000000000000000000000054a2d42a40f51259dedd1978f6c118a0f0eff078"

	params := args{
		address:   address,
		fromBlock: "0x429d3b",
		toBlock:   Pending,
		topics:    topics,
	}

	want := NewLogRequest(address, "0x429d3b", Pending, topics...)

	tests := []struct {
		name string
		args args
		want LogRequest
	}{
		{"Test New", params, want},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLogRequest(tt.args.address, tt.args.fromBlock, tt.args.toBlock, tt.args.topics...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLogRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
