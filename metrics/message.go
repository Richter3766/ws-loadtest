// metrics/message.go
package metrics

import "sync/atomic"

type MessageMetrics struct {
    Sent   int64
    Recv   int64
}

func (m *MessageMetrics) IncSent() {
    atomic.AddInt64(&m.Sent, 1)
}
func (m *MessageMetrics) IncRecv() {
    atomic.AddInt64(&m.Recv, 1)
}
