// metrics/latency.go
package metrics

import (
	"sort"
	"sync"
	"time"
)

type LatencyMetrics struct {
    mu        sync.Mutex
    Latencies []time.Duration
}

type LatencyStats struct {
    Count      int     `json:"count"`
    AvgMs      float64 `json:"avg_ms"`
    MaxMs      float64 `json:"max_ms"`
    P95Ms      float64 `json:"p95_ms"`
}

func (m *LatencyMetrics) Add(latency time.Duration) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.Latencies = append(m.Latencies, latency)
}


// 통계(평균, 최댓값, 95th percentile 등)
func (m *LatencyMetrics) Stats() (count int, avgMs, maxMs, p95Ms float64) {
    m.mu.Lock()
    defer m.mu.Unlock()
    n := len(m.Latencies)
    if n == 0 {
        return 0, 0, 0, 0
    }
    total := time.Duration(0)
    max := m.Latencies[0]
    sorted := append([]time.Duration{}, m.Latencies...)
    sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
    for _, v := range m.Latencies {
        total += v
        if v > max {
            max = v
        }
    }
    avgMs = float64(total.Nanoseconds()) / float64(n) / 1e6
    maxMs = float64(max.Nanoseconds()) / 1e6
    p95Idx := int(float64(n-1) * 0.95)
    p95Ms = float64(sorted[p95Idx].Nanoseconds()) / 1e6
    return n, avgMs, maxMs, p95Ms
}
