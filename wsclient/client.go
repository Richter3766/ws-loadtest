package wsclient

import (
	"encoding/json"
	"log"
	"sync"
	"time"
	"ws-loadtest/metrics"
	// ... 기타 필요 패키지 ...
)
func RunClient(clientNum int, wsURL string, roomID int, jwt string, stopAll <-chan struct{}, wg *sync.WaitGroup) {
    defer wg.Done()

    conn, wsnet, err := ConnectWSSTOMP(wsURL, jwt)
    if err != nil {
        log.Printf("[클라이언트 %d] 연결 실패: %v", clientNum, err)
        return
    }

    done := make(chan struct{})
    var subWg sync.WaitGroup
    subWg.Add(1)
    err = Subscribe(conn, roomID, clientNum, handleMsg, done, &subWg)
    if err != nil {
        log.Printf("[클라이언트 %d] 구독 실패: %v", clientNum, err)
        return
    }

    go PublishLoopWithID(conn, roomID, 1, clientNum)

    <-stopAll
    close(done)
    
	conn.Disconnect()
	wsnet.Close()
    log.Printf("[클라이언트 %d] 종료", clientNum)
}

// TODO: 각 클라별로 집계 데이터 쌓아 통계적 분석 해보기
func handleMsg(body string, clientNum int) {
    // 1. 응답 메시지 파싱
    var resp MessageResponseDto
    if err := json.Unmarshal([]byte(body), &resp); err != nil {
        log.Printf("메시지 Unmarshal 오류: %v", err)
        return
    }
    metrics.Default.Delivery.AddRecv(clientNum, resp.ClientMessageId)
    // 2. 레이턴시 계산 (현재 수신 시각 - clientSentAt)
    if resp.ClientSentAt > 0 {
        sentTime := time.Unix(0, resp.ClientSentAt) // UnixNano 기준
        receivedTime := time.Now()
        latency := receivedTime.Sub(sentTime)

        // 3. 중앙 지표 저장
        metrics.Default.Latency.Add(latency)
    }
    // log.Println(body)
}
