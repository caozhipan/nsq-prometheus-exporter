package controllers

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type ChannelMetric struct {
	val func(channel *Channel) float64
	vec *prometheus.GaugeVec
}

func generateChannelMetrics() []*ChannelMetric {
	labels := []string{"node", "topic", "channel", "paused"}
	namespace := "nsq_channel"
	return []*ChannelMetric{
		&ChannelMetric{
			val: func(channel *Channel) float64 {
				return float64(len(channel.Clients))
			},
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "client_count",
				Help:      "Number of clients",
			}, labels),
		},
		&ChannelMetric{
			val: func(channel *Channel) float64 { return channel.Depth },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "depth",
				Help:      "Queue depth",
			}, labels),
		},
		&ChannelMetric{
			val: func(channel *Channel) float64 { return channel.BackendDepth },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "backend_depth",
				Help:      "Queue backend depth",
			}, labels),
		},
		&ChannelMetric{
			val: func(channel *Channel) float64 { return channel.MessageCount },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "message_count",
				Help:      "Queue message count",
			}, labels),
		},
		&ChannelMetric{
			val: func(channel *Channel) float64 { return channel.InFlightCount },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "in_flight_count",
				Help:      "In flight count",
			}, labels),
		},

		&ChannelMetric{
			val: func(channel *Channel) float64 { return channel.E2eLatency.percentileValue(0) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "e2e_latency_99p",
				Help:      "e2e latency 99th percentile",
			}, labels),
		},

		&ChannelMetric{
			val: func(channel *Channel) float64 { return channel.E2eLatency.percentileValue(1) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "e2e_latency_95p",
				Help:      "e2e latency 95th percentile",
			}, labels),
		},

		&ChannelMetric{
			val: func(channel *Channel) float64 { return channel.DeferredCount },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "deferred_count",
				Help:      "Deferred count",
			}, labels),
		},

		&ChannelMetric{
			val: func(channel *Channel) float64 { return channel.RequeueCount },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "requeue_count",
				Help:      "Requeue Count",
			}, labels),
		},

		&ChannelMetric{
			val: func(channel *Channel) float64 { return channel.TimeoutCount },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "timeout_count",
				Help:      "Timeout count",
			}, labels),
		},
	}

}

func (m *ChannelMetric) collect(stats []*Stats, out chan<- prometheus.Metric) {
	for _, s := range stats {
		for _, topic := range s.Topics {
			for _, channel := range topic.Channels {
				m.vec.With(prometheus.Labels{
					"node":    s.Node.HostName,
					"topic":   topic.Name,
					"channel": channel.Name,
					"paused":  strconv.FormatBool(channel.Paused),
				}).Set(m.val(channel))
			}
		}
	}
	m.vec.Collect(out)
}

func (m *ChannelMetric) describe(ch chan<- *prometheus.Desc) {
	m.vec.Describe(ch)
}
func (m *ChannelMetric) reset() {
	m.vec.Reset()
}
