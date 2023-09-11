package prommetrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	Requests              prometheus.CounterVec
	FailedTxs             prometheus.CounterVec
	RequestsLatency       prometheus.HistogramVec
	HistoricQueries       prometheus.GaugeVec
	SendQueue             prometheus.GaugeVec
	HistoricQueryRequests prometheus.CounterVec
	ABCIRequests          prometheus.CounterVec
	LightBlockRequests    prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		Requests: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "icq",
			Name:      "requests",
			Help:      "number of host requests",
		}, []string{"name", "type"}),
		FailedTxs: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "icq",
			Name:      "failed_txs",
			Help:      "number of failed txs",
		}, []string{"name"}),
		RequestsLatency: *prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "icq",
			Name:      "request_duration_seconds",
			Help:      "Latency of requests",
			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
		}, []string{"name", "type"}),
		HistoricQueries: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "icq",
			Name:      "historic_queries",
			Help:      "historic queue size",
		}, []string{"name"}),
		SendQueue: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "icq",
			Name:      "send_queue",
			Help:      "send queue size",
		}, []string{"name"}),
		HistoricQueryRequests: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "icq",
			Name:      "historic_reqs",
			Help:      "number of historic query requests",
		}, []string{"name"}),
		ABCIRequests: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "icq",
			Name:      "abci_reqs",
			Help:      "number of abci requests",
		}, []string{"name", "type"}),
		LightBlockRequests: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "icq",
			Name:      "lightblock_reqs",
			Help:      "number of lightblock requests",
		}, []string{"name"}),
	}
	reg.MustRegister(m.Requests, m.RequestsLatency, m.HistoricQueries, m.SendQueue, m.FailedTxs, m.HistoricQueryRequests, m.LightBlockRequests, m.ABCIRequests)
	return m
}
