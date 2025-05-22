// metrics/metrics.go
package metrics

import (
	"encoding/json"
	"os"
)

type Metrics struct {
    Connection *ConnectionMetrics
    Message    *MessageMetrics
    Latency    *LatencyMetrics
    Stability  *StabilityMetrics
    Delivery   *DeliveryStats
    // 추가 지표 있으면 확장
}

// 전역 싱글톤 (가장 간단)
var Default = &Metrics{
    Connection: &ConnectionMetrics{},
    Message:    &MessageMetrics{},
    Latency:    &LatencyMetrics{},
    Stability:  &StabilityMetrics{},
    Delivery:  NewDeliveryStats(),
}

type ExportedMetrics struct {
    Connection *ConnectionMetrics `json:"connection"`
    Message    *MessageMetrics    `json:"message"`
    Latency    LatencyStats       `json:"latency"`
    Stability  *StabilityMetrics  `json:"stability"`
    LostMessages LostMessagesStats   `json:"lost_messages"`
}

func SaveMetricsJSON(filename string) error {
    count, avgMs, maxMs, p95Ms := Default.Latency.Stats()
    lost := CalcLostMessages()
    export := ExportedMetrics{
        Connection:   Default.Connection,
        Message:      Default.Message,
        Latency:      LatencyStats{Count: count, AvgMs: avgMs, MaxMs: maxMs, P95Ms: p95Ms},
        Stability:    Default.Stability,
        LostMessages: lost,
    }
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    enc := json.NewEncoder(file)
    enc.SetIndent("", "  ")
    return enc.Encode(export)
}

func ResetAll() {
    *Default = Metrics{
        Connection: &ConnectionMetrics{},
        Message:    &MessageMetrics{},
        Latency:    &LatencyMetrics{},
        Stability:  &StabilityMetrics{},
        Delivery:   NewDeliveryStats(),
    }
}
