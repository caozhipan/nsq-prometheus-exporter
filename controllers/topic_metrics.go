package controllers

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type TopicMetric struct {
	val func(topic *Topic) float64
	vec *prometheus.GaugeVec
}

func generateTopicMetrics() []*TopicMetric {
	labels := []string{"node", "topic", "paused"}
	namespace := "nsq_topic"
	return []*TopicMetric{
		&TopicMetric{
			val: func(Topic *Topic) float64 {
				return float64(len(Topic.Channels))
			},
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "channel_count",
				Help:      "Number of channel",
			}, labels),
		},
		&TopicMetric{
			val: func(topic *Topic) float64 { return topic.Depth },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "depth",
				Help:      "Queue depth",
			}, labels),
		},
		&TopicMetric{
			val: func(topic *Topic) float64 { return topic.BackendDepth },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "backend_depth",
				Help:      "Queue backend depth",
			}, labels),
		},
		&TopicMetric{
			val: func(topic *Topic) float64 { return topic.MessageCount },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "message_count",
				Help:      "Queue message count",
			}, labels),
		},

		&TopicMetric{
			val: func(topic *Topic) float64 { return topic.E2eLatency.percentileValue(0) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "e2e_latency_99p",
				Help:      "e2e latency 99th percentile",
			}, labels),
		},

		&TopicMetric{
			val: func(topic *Topic) float64 { return topic.E2eLatency.percentileValue(1) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "e2e_latency_95p",
				Help:      "e2e latency 95th percentile",
			}, labels),
		},
	}

}

func (m *TopicMetric) collect(stats []*Stats, out chan<- prometheus.Metric) {
	for _, s := range stats {
		for _, topic := range s.Topics {
			m.vec.With(prometheus.Labels{
				"node":   s.Node.HostName,
				"topic":  topic.Name,
				"paused": strconv.FormatBool(topic.Paused),
			}).Set(m.val(topic))
		}
	}
	m.vec.Collect(out)
}

func (m *TopicMetric) describe(ch chan<- *prometheus.Desc) {
	m.vec.Describe(ch)
}
func (m *TopicMetric) reset() {
	m.vec.Reset()
}
