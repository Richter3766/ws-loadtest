// metrics/connection.go
package metrics

import "sync/atomic"

type ConnectionMetrics struct {
    WsConnSuccess   int64
    WsConnFail      int64
    StompConnSuccess int64
    StompConnFail    int64
}

func (m *ConnectionMetrics) IncWsConnSuccess() {
    atomic.AddInt64(&m.WsConnSuccess, 1)
}
func (m *ConnectionMetrics) IncWsConnFail() {
    atomic.AddInt64(&m.WsConnFail, 1)
}
func (m *ConnectionMetrics) IncStompConnSuccess() {
    atomic.AddInt64(&m.StompConnSuccess, 1)
}
func (m *ConnectionMetrics) IncStompConnFail() {
    atomic.AddInt64(&m.StompConnFail, 1)
}
