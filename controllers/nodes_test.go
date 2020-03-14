package controllers

import (
	"strings"
	"testing"
)

const (
	TestNsqLookupdAddressLists = "127.0.0.1:4161"
)

func TestSyncNodeList(t *testing.T) {
	type args struct {
		lookupdAddrs string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		struct {
			name string
			args args
		}{name: "test:sync", args: args{lookupdAddrs: TestNsqLookupdAddressLists}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SyncNodeList(tt.args.lookupdAddrs)
		})
	}
}

func Test_getNodeList(t *testing.T) {
	type args struct {
		lookupdAddr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		struct {
			name    string
			args    args
			wantErr bool
		}{name: "test:getNodeList", args: args{lookupdAddr: strings.Split(TestNsqLookupdAddressLists, ",")[0]}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getNodeList(tt.args.lookupdAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNodeList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
