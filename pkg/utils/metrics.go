package utils

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	name   string
	port   int
	server *http.Server

	peersTotal       prometheus.Counter
	peersCurrent     prometheus.Gauge
	messagesReceived prometheus.Counter
	messagesSent     prometheus.Counter
	requestsTotal    prometheus.Counter
	requestsDuration *prometheus.HistogramVec
	errorsTotal      prometheus.Counter

	mu sync.RWMutex
}

func NewMetrics(name string) (*Metrics, error) {
	m := &Metrics{
		name: name,
	}

	registry := prometheus.NewRegistry()

	m.peersTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_peers_total", name),
		Help: "Total number of peers discovered",
	})

	m.peersCurrent = promauto.NewGauge(prometheus.GaugeOpts{
		Name: fmt.Sprintf("%s_peers_current", name),
		Help: "Current number of connected peers",
	})

	m.messagesReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_messages_received_total", name),
		Help: "Total number of messages received",
	})

	m.messagesSent = promauto.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_messages_sent_total", name),
		Help: "Total number of messages sent",
	})

	m.requestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_requests_total", name),
		Help: "Total number of requests",
	})

	m.requestsDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    fmt.Sprintf("%s_requests_duration_seconds", name),
		Help:    "Request duration in seconds",
		Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
	}, []string{"method", "status"})

	m.errorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_errors_total", name),
		Help: "Total number of errors",
	})

	registry.MustRegister(
		m.peersTotal,
		m.peersCurrent,
		m.messagesReceived,
		m.messagesSent,
		m.requestsTotal,
		m.requestsDuration,
		m.errorsTotal,
	)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	m.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", 9090),
		Handler: mux,
	}

	return m, nil
}

func (m *Metrics) Start(ctx context.Context, port int) error {
	m.port = port
	m.server.Addr = fmt.Sprintf(":%d", port)

	go func() {
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Metrics server error: %v\n", err)
		}
	}()

	return nil
}

func (m *Metrics) Stop(ctx context.Context) error {
	if m.server == nil {
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return m.server.Shutdown(shutdownCtx)
}

func (m *Metrics) IncPeers() {
	m.peersTotal.Inc()
	m.peersCurrent.Inc()
}

func (m *Metrics) DecPeers() {
	m.peersCurrent.Dec()
}

func (m *Metrics) IncMessagesReceived() {
	m.messagesReceived.Inc()
}

func (m *Metrics) IncMessagesSent() {
	m.messagesSent.Inc()
}

func (m *Metrics) IncRequests() {
	m.requestsTotal.Inc()
}

func (m *Metrics) ObserveRequestDuration(method, status string, duration time.Duration) {
	m.requestsDuration.WithLabelValues(method, status).Observe(duration.Seconds())
}

func (m *Metrics) IncErrors() {
	m.errorsTotal.Inc()
}

func (m *Metrics) SetPeers(count int) {
	m.peersCurrent.Set(float64(count))
}
