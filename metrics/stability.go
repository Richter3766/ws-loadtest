// metrics/stability.go
package metrics

import "sync/atomic"

type StabilityMetrics struct {
    UnexpectedDisconnects int64
    Timeouts              int64
    // 추가로 필요한 이벤트는 여기에 확장
}

func (m *StabilityMetrics) IncDisconnect() {
    atomic.AddInt64(&m.UnexpectedDisconnects, 1)
}
func (m *StabilityMetrics) IncTimeout() {
    atomic.AddInt64(&m.Timeouts, 1)
}
