package controllers

import "github.com/prometheus/client_golang/prometheus"

type Metric interface {
	collect(stats []*Stats, out chan<- prometheus.Metric)
	describe(ch chan<- *prometheus.Desc)
	reset()
}

type NsqCollector struct {
	ChannelMetrics []*ChannelMetric
	TopicMetrics   []*TopicMetric
	ClientMetrics  []*ClientMetric
	NodeMetrics    []*NodeMetric
}

var (
	Collector *NsqCollector = &NsqCollector{
		ChannelMetrics: generateChannelMetrics(),
		TopicMetrics:   generateTopicMetrics(),
		ClientMetrics:  generateClientMetrics(),
		NodeMetrics:    generateNodeMetrics(),
	}
)

func (c *NsqCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.ChannelMetrics {
		metric.describe(ch)
	}
}

func (c *NsqCollector) Collect(ch chan<- prometheus.Metric) {
	stats, err := getNsqdStats()
	if err != nil {
		return
	}
	for _, metric := range c.ChannelMetrics {
		metric.collect(stats, ch)
	}

	for _, metric := range c.TopicMetrics {
		metric.collect(stats, ch)
	}

	for _, metric := range c.ClientMetrics {
		metric.collect(stats, ch)
	}

	for _, metric := range c.NodeMetrics {
		metric.collect(stats, ch)
	}

}
