package controllers

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

type NodeMetric struct {
	val func(node Node) float64
	vec *prometheus.GaugeVec
}

func generateNodeMetrics() []*NodeMetric {
	labels := []string{"hostname", "remote_address", "broadcast_address", "http_port", "tcp_port", "version"}
	namespace := "nsq_node"
	return []*NodeMetric{
		&NodeMetric{
			val: func(node Node) float64 { return float64(len(node.Topics)) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "topic_count",
				Help:      "topic count",
			}, labels),
		},
	}

}

func (m *NodeMetric) collect(stats []*Stats, out chan<- prometheus.Metric) {
	for _, node := range NsqNodes.Producers {
		m.vec.With(prometheus.Labels{
			"hostname":          node.HostName,
			"remote_address":    node.RemoteAddress,
			"broadcast_address": node.BroadcastAddress,
			"http_port":         fmt.Sprintf("%d", node.HttpPort),
			"tcp_port":          fmt.Sprintf("%d", node.TcpPort),
			"version":           node.Version,
		}).Set(m.val(node))
	}
	m.vec.Collect(out)
}

func (m *NodeMetric) describe(ch chan<- *prometheus.Desc) {
	m.vec.Describe(ch)
}
func (m *NodeMetric) reset() {
	m.vec.Reset()
}
