package metrics

import "sync"

type DeliveryStats struct {
	mu      sync.Mutex
	SentIDs map[string]bool         // 전체 발송한 MessageID
	RecvIDs map[int]map[string]bool // 각 클라이언트별 수신 MessageID 집합
}

func NewDeliveryStats() *DeliveryStats {
	return &DeliveryStats{
		SentIDs: make(map[string]bool),
		RecvIDs: make(map[int]map[string]bool),
	}
}

type LostMessagesStats struct {
    TotalLostCount int                 `json:"total_lost_count"`
    LostPerClient  map[int]int         `json:"lost_per_client"`
    LostIds        map[int][]string    `json:"lost_message_ids,omitempty"` // 클라별 누락 메시지ID 리스트(옵션)
}

// 메시지 발송시 호출
func (d *DeliveryStats) AddSent(id string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.SentIDs[id] = true
}

// 메시지 수신시 호출
func (d *DeliveryStats) AddRecv(clientNum int, id string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.RecvIDs[clientNum]; !ok {
		d.RecvIDs[clientNum] = make(map[string]bool)
	}
	d.RecvIDs[clientNum][id] = true
}

// 클라이언트별 손실 메시지 집계
func (d *DeliveryStats) GetLostMessagesPerClient() map[int][]string {
    d.mu.Lock()
    defer d.mu.Unlock()
    lost := make(map[int][]string)
    for clientNum, recvSet := range d.RecvIDs {
        for id := range d.SentIDs {
            if !recvSet[id] {
                lost[clientNum] = append(lost[clientNum], id)
            }
        }
    }
    return lost
}

func CalcLostMessages() LostMessagesStats {
    lostPerClient := map[int]int{}
    lostIds := map[int][]string{}
    lostTotal := 0

    // 클라별 누락 메시지 집계
    lostMap := Default.Delivery.GetLostMessagesPerClient()
    for clientNum, lostList := range lostMap {
        lostPerClient[clientNum] = len(lostList)
        lostIds[clientNum] = lostList
        lostTotal += len(lostList)
    }

    return LostMessagesStats{
        TotalLostCount: lostTotal,
        LostPerClient:  lostPerClient,
        LostIds:        lostIds, // 대용량이면 빼도 됨
    }
}



// 전체 손실 메시지(누군가 한 명이라도 못 받은 메시지)
// func (d *DeliveryStats) GetGlobalLostMessages() []string {
//     d.mu.Lock()
//     defer d.mu.Unlock()
//     lost := []string{}
//     for id := range d.SentIDs {
//         deliveredToAll := true
//         for _, recvSet := range d.RecvIDs {
//             if !recvSet[id] {
//                 deliveredToAll = false
//                 break
//             }
//         }
//         if !deliveredToAll {
//             lost = append(lost, id)
//         }
//     }
//     return lost
// }
