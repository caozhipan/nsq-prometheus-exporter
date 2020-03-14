package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Stats struct {
	Version   string   `json:"version"`
	Health    string   `json:"health"`
	StartTime int64    `json:"start_time"`
	Topics    []*Topic `json:"topics"`
	Node      Node
}

// see https://github.com/nsqio/nsq/blob/master/nsqd/stats.go
type Topic struct {
	Name         string     `json:"topic_name"`
	Paused       bool       `json:"paused"`
	Depth        float64    `json:"depth"`
	BackendDepth float64    `json:"backend_depth"`
	MessageCount float64    `json:"message_count"`
	E2eLatency   E2elatency `json:"e2e_processing_latency"`
	Channels     []*Channel `json:"channels"`
}

type Channel struct {
	Name          string     `json:"channel_name"`
	Paused        bool       `json:"paused"`
	Depth         float64    `json:"depth"`
	BackendDepth  float64    `json:"backend_depth"`
	MessageCount  float64    `json:"message_count"`
	InFlightCount float64    `json:"in_flight_count"`
	DeferredCount float64    `json:"deferred_count"`
	RequeueCount  float64    `json:"requeue_count"`
	TimeoutCount  float64    `json:"timeout_count"`
	E2eLatency    E2elatency `json:"e2e_processing_latency"`
	Clients       []*Client  `json:"clients"`
}

type E2elatency struct {
	Count       float64              `json:"count"`
	Percentiles []map[string]float64 `json:"percentiles"`
}

func (e *E2elatency) percentileValue(idx int) float64 {
	if idx >= len(e.Percentiles) {
		return 0
	}
	return e.Percentiles[idx]["value"]
}

type Client struct {
	ID            string  `json:"client_id"`
	Hostname      string  `json:"hostname"`
	Version       string  `json:"version"`
	RemoteAddress string  `json:"remote_address"`
	State         int32   `json:"state"`
	FinishCount   float64 `json:"finish_count"`
	MessageCount  float64 `json:"message_count"`
	ReadyCount    float64 `json:"ready_count"`
	InFlightCount float64 `json:"in_flight_count"`
	RequeueCount  float64 `json:"requeue_count"`
	ConnectTime   int64   `json:"connect_ts"`
	SampleRate    float64 `json:"sample_rate"`
	Deflate       bool    `json:"deflate"`
	Snappy        bool    `json:"snappy"`
	TLS           bool    `json:"tls"`
	UserAgent     string  `json:"user_agent"`
}

func getPercentile(t *Topic, percentile int) float64 {
	if len(t.E2eLatency.Percentiles) > 0 {
		if percentile == 99 {
			return t.E2eLatency.Percentiles[0]["value"]
		} else if percentile == 95 {
			return t.E2eLatency.Percentiles[1]["value"]
		}
	}
	return 0
}

func getNsqdStatsByNode(node Node) (*Stats, error) {
	nsqdURL := fmt.Sprintf("http://%s:%d/stats?format=json", node.BroadcastAddress, node.HttpPort)
	resp, err := http.Get(nsqdURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var stats Stats
	if err = json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}
	stats.Node = node
	return &stats, nil
}

func getNsqdStats() (stats []*Stats, err error) {
	if NsqNodes == nil || len(NsqNodes.Producers) == 0 {
		panic("cannot find any nodes")
	}
	for _, node := range NsqNodes.Producers {
		st, err := getNsqdStatsByNode(node)
		if err != nil {
			return nil, err
		}
		stats = append(stats, st)
	}
	return
}
