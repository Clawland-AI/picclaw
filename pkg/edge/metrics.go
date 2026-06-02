package edge

import (
	"expvar"
	"sync"
	"sync/atomic"
	"time"
)

// ReporterMetrics holds live counters and gauges for the edge reporter.
// Exported via expvar under "edgereporter" so operators can scrape them
// alongside any Prometheus /debug/vars exporter.
type ReporterMetrics struct {
	mu sync.RWMutex

	// Counters (cumulative)
	RegistrationsTotal  atomic.Int64 `json:"registrations_total"`
	HeartbeatsTotal     atomic.Int64 `json:"heartbeats_total"`
	StatusReportsTotal  atomic.Int64 `json:"status_reports_total"`
	EventsTotal         atomic.Int64 `json:"events_total"`
	QueuedEventsTotal   atomic.Int64 `json:"queued_events_total"`
	DeliveriesTotal     atomic.Int64 `json:"deliveries_total"`      // all HTTP POST attempts
	DeliveriesOK        atomic.Int64 `json:"deliveries_ok"`         // 2xx responses
	DeliveriesFailed    atomic.Int64 `json:"deliveries_failed"`     // non-2xx or network errors
	RetriesTotal        atomic.Int64 `json:"retries_total"`         // retry attempts due to failure
	DroppedEventsTotal  atomic.Int64 `json:"dropped_events_total"`  // events dropped when queue full

	// Gauges (snapshots)
	QueueDepth     atomic.Int64 `json:"queue_depth"`
	LastDeliveryOK atomic.Bool  `json:"last_delivery_ok"`

	// Latency buckets (ms) for the most recent deliveries (ring buffer)
	latencyMu   sync.Mutex
	latencyRing [64]time.Duration
	latencyIdx  int
}

var metrics = &ReporterMetrics{}

func init() {
	expvar.Publish("edgereporter", expvar.Func(func() any {
		return metrics.Snapshot()
	}))
}

// Snapshot returns a consistent read of all metrics.
func (m *ReporterMetrics) Snapshot() map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]any{
		"registrations_total":  m.RegistrationsTotal.Load(),
		"heartbeats_total":     m.HeartbeatsTotal.Load(),
		"status_reports_total": m.StatusReportsTotal.Load(),
		"events_total":         m.EventsTotal.Load(),
		"queued_events_total":  m.QueuedEventsTotal.Load(),
		"deliveries_total":     m.DeliveriesTotal.Load(),
		"deliveries_ok":        m.DeliveriesOK.Load(),
		"deliveries_failed":    m.DeliveriesFailed.Load(),
		"retries_total":        m.RetriesTotal.Load(),
		"dropped_events_total": m.DroppedEventsTotal.Load(),
		"queue_depth":          m.QueueDepth.Load(),
		"last_delivery_ok":     m.LastDeliveryOK.Load(),
	}
}

func (m *ReporterMetrics) recordLatency(d time.Duration) {
	m.latencyMu.Lock()
	m.latencyRing[m.latencyIdx%64] = d
	m.latencyIdx++
	m.latencyMu.Unlock()
}
