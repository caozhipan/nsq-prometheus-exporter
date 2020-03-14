package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Node struct {
	RemoteAddress    string   `json:"remote_address"`
	HostName         string   `json:"hostname"`
	BroadcastAddress string   `json:"broadcast_address"`
	TcpPort          int      `json:"tcp_port"`
	HttpPort         int      `json:"http_port"`
	Version          string   `json:"version"`
	Tombstones       []bool   `json:"tombstones"`
	Topics           []string `json:"topics"`
}

type Nodes struct {
	Producers []Node `json:"producers"`
}

var (
	NsqNodes *Nodes
)

func init() {
	http.DefaultClient.Timeout = 5 * time.Second

}

func SyncNodeList(lookupdAddrs string) {
	addrList := strings.Split(lookupdAddrs, ",")
	if len(addrList) == 0 {
		panic("lookupdAddrs cannot be null")
	}

	for _, addr := range addrList {
		nodes, err := getNodeList(addr)
		if err == nil && nodes != nil {
			NsqNodes = nodes
			return
		}
	}

	panic("cannot sync nodes list from lookupd")
}

func getNodeList(lookupdAddr string) (nodes *Nodes, err error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/nodes", lookupdAddr))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&nodes)
	return
}
