package btc_rpc_client

import (
	go_logger "github.com/pefish/go-logger"
	"testing"
	"time"
)

func TestUnisatHttpClient_ListBrc20Balances(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: struct{ address string }{address: "bc1ptsl68uq68xtmzr46snemrf4rextp5kwzkaahjzd36ktx6dvz80yq9k02kc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uhc := NewUnisatHttpClient(go_logger.Logger, 3*time.Second)
			got, err := uhc.ListBrc20Balances(tt.args.address)
			go_logger.Logger.Info(got, err)
		})
	}
}
