package controllers

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

type ClientMetric struct {
	val func(client *Client) float64
	vec *prometheus.GaugeVec
}

func generateClientMetrics() []*ClientMetric {
	labels := []string{"node", "topic", "channel", "client_id", "hostname", "version", "remote_address", "state", "deflate", "snappy", "user_agent", "tls"}
	namespace := "nsq_client"
	return []*ClientMetric{
		&ClientMetric{
			val: func(client *Client) float64 { return client.FinishCount },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "finish_count",
				Help:      "Finish count",
			}, labels),
		},
		&ClientMetric{
			val: func(client *Client) float64 { return client.MessageCount },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "message_count",
				Help:      "Queue message count",
			}, labels),
		},
		&ClientMetric{
			val: func(client *Client) float64 { return client.ReadyCount },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "ready_count",
				Help:      "Ready count",
			}, labels),
		},
		&ClientMetric{
			val: func(client *Client) float64 { return client.InFlightCount },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "in_flight_count",
				Help:      "In flight count",
			}, labels),
		},

		&ClientMetric{
			val: func(client *Client) float64 { return client.RequeueCount },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "requeue_count",
				Help:      "Requeue count",
			}, labels),
		},

		&ClientMetric{
			val: func(client *Client) float64 { return client.SampleRate },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "sample_rate",
				Help:      "Sample Rate",
			}, labels),
		},
	}

}

func (m *ClientMetric) collect(stats []*Stats, out chan<- prometheus.Metric) {
	for _, s := range stats {
		for _, topic := range s.Topics {
			for _, channel := range topic.Channels {
				for _, client := range channel.Clients {
					m.vec.With(prometheus.Labels{
						"node":           s.Node.HostName,
						"topic":          topic.Name,
						"channel":        channel.Name,
						"client_id":      client.ID,
						"hostname":       client.Hostname,
						"version":        client.Version,
						"remote_address": client.RemoteAddress,
						"state":          fmt.Sprintf("%d", client.State),
						"deflate":        fmt.Sprintf("%v", client.Deflate),
						"snappy":         fmt.Sprintf("%v", client.Snappy),
						"user_agent":     client.UserAgent,
						"tls":            fmt.Sprintf("%v", client.TLS),
					}).Set(m.val(client))
				}
			}
		}
	}
	m.vec.Collect(out)
}

func (m *ClientMetric) describe(ch chan<- *prometheus.Desc) {
	m.vec.Describe(ch)
}
func (m *ClientMetric) reset() {
	m.vec.Reset()
}
