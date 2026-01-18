package metrics

import (
	"sync"
	"time"
)

// Metrics holds application metrics
type Metrics struct {
	mu sync.RWMutex

	// Request counts
	TotalRequests   int64
	SuccessRequests int64 // 2xx
	ClientErrors    int64 // 4xx
	ServerErrors    int64 // 5xx

	// Latency tracking
	TotalLatency time.Duration
	MinLatency   time.Duration
	MaxLatency   time.Duration

	// Status code breakdown
	StatusCodes map[int]int64

	// Endpoint tracking
	EndpointCounts map[string]int64
}

// MetricsSnapshot is a read-only snapshot of metrics
type MetricsSnapshot struct {
	TotalRequests   int64
	SuccessRequests int64
	ClientErrors    int64
	ServerErrors    int64
	AvgLatency      time.Duration
	MinLatency      time.Duration
	MaxLatency      time.Duration
	StatusCodes     map[int]int64
	EndpointCounts  map[string]int64
}

var once sync.Once
var instance *Metrics

func GetInstance() *Metrics {
	once.Do(func() {
		instance = &Metrics{
			StatusCodes: make(map[int]int64),
			EndpointCounts: make(map[string]int64),
			MinLatency: time.Hour,
		}
	})
	return instance
}

func (m *Metrics) RecordRequest(statusCode int, latency time.Duration, endpoint string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.TotalRequests++

	// Status code categories
	switch {
	case statusCode >= 200 && statusCode < 300:
		m.SuccessRequests++
	case statusCode >= 400 && statusCode < 500:
		m.ClientErrors++
	case statusCode >= 500:
		m.ServerErrors++
	}

	// Status code breakdown
	m.StatusCodes[statusCode]++

	// Endpoint tracking
	m.EndpointCounts[endpoint]++

	m.TotalLatency += latency

	if latency < m.MinLatency {
		m.MinLatency = latency
	}

	if latency > m.MaxLatency {
		m.MaxLatency = latency
	}

}

func (m *Metrics) GetSnapshot() MetricsSnapshot {
	m.mu.RLock()

	defer m.mu.RUnlock()

	//ask this question after you are done with this section
	avgLatency := time.Duration(0)

	if m.TotalRequests > 0 {
		avgLatency = m.TotalLatency / time.Duration(m.TotalRequests)
	}

	//copy maps to avoid race conditions 
	statusCodes := make(map[int]int64)

	for k, v := range m.StatusCodes {
		statusCodes[k] = v
	}

	endpointCounts := make(map[string]int64)

	for k, v := range m.EndpointCounts {
		endpointCounts[k] = v
	}

	return MetricsSnapshot{
		TotalRequests:   m.TotalRequests,
		SuccessRequests: m.SuccessRequests,
		ClientErrors:    m.ClientErrors,
		ServerErrors:    m.ServerErrors,
		AvgLatency:      avgLatency,
		MinLatency:      m.MinLatency,
		MaxLatency:      m.MaxLatency,
		StatusCodes:     statusCodes,
		EndpointCounts:  endpointCounts,
	}

}