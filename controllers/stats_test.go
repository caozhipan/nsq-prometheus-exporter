package controllers

import (
	"fmt"
	"testing"
)

func Test_getNsqdStatsByNode(t *testing.T) {

	SyncNodeList(TestNsqLookupdAddressLists)

	type args struct {
		node *Node
	}
	type TestParam struct {
		name    string
		args    args
		wantErr bool
	}

	var params []TestParam

	for _, node := range NsqNodes.Producers {
		params = append(params, TestParam{
			name:    fmt.Sprintf("test:node_%s", node.HostName),
			args:    args{node: &node},
			wantErr: false,
		})
	}
	for _, tt := range params {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getNsqdStatsByNode(tt.args.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNsqdStatsByNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_getNsqdStats(t *testing.T) {
	SyncNodeList(TestNsqLookupdAddressLists)
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
		struct {
			name    string
			wantErr bool
		}{name: "Test_getNsqdStats", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getNsqdStats()
			if (err != nil) != tt.wantErr {
				t.Errorf("getNsqdStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
