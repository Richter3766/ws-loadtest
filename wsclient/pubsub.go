package wsclient

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-stomp/stomp"
	"github.com/google/uuid"

	"ws-loadtest/metrics"
)

// Subscribe : 메시지 구독 함수 (비동기로 수신 콜백)
func Subscribe(conn *stomp.Conn, roomID int, clientNum int, done <-chan struct{}, subWg *sync.WaitGroup, receivedSet *SafeSet) (error) {
    sub, err := conn.Subscribe(fmt.Sprintf("/sub/room/%d", roomID), stomp.AckAuto)
    if err != nil {
        return err
    }
    go func() {
        defer subWg.Done()
        for {
            select {
            case <-done:
                log.Println("done 신호로 구독 루프 종료")
                return
            case msg, ok := <-sub.C:
                if !ok {
                    log.Println("[클라이언트", clientNum, "] sub.C 채널 예기치 않게 닫힘 → 비정상 종료")
                    metrics.Default.Stability.IncDisconnect()
                    return
                }
                if msg.Err != nil {
                    log.Printf("메시지 수신 오류: %v", msg.Err)
                    continue
                }
                metrics.Default.Message.IncRecv() // 받은 메세지 지표 수집
                messageId, roomId  := handleMsg(string(msg.Body), clientNum, receivedSet)
                if messageId != -1 {
                    sendAck(conn, messageId, roomId)
                }
            }
        }
    }()
    return nil
}

func SubscribeNotify(conn *stomp.Conn, clientNum int, done <-chan struct{}, subWg *sync.WaitGroup, receivedSet *SafeSet) error {
    sub, err := conn.Subscribe("/user/queue/notify", stomp.AckAuto)
    if err != nil {
        return err
    }
    go func() {
        defer subWg.Done()
        for {
            select {
            case <-done:
                log.Printf("[클라이언트 %d] notify 구독 루프 종료", clientNum)
                return
            case msg, ok := <-sub.C:
                if !ok {
                    log.Printf("[클라이언트 %d] notify 채널 예기치 않게 닫힘", clientNum)
                    return
                }
                // 누락 메시지 처리
                metrics.Default.Message.IncRecv()
                messageId, roomId := handleMsg(string(msg.Body), clientNum, receivedSet)
                // log.Printf("클라이언트 %d: ack 전송", clientNum)
                if messageId != -1 {
                    log.Printf("클라이언트 %d: ack 전송", clientNum)
                    sendAck(conn, messageId, roomId)
                }
            }
        }
    }()
    return nil
}


func sendAck(conn *stomp.Conn, messageId int64, roomId int64) {
    // 서버에서 정의한 /pub/ack, 메시지 포맷에 맞춰 전송
    payload := fmt.Sprintf(`{"messageId":"%d", "roomId":"%d"}`, messageId, roomId)
    err := conn.Send("/pub/ack", "application/json", []byte(payload))
    if err != nil {
        log.Printf("ACK 전송 실패: %v", err)
    }
}

// Publish : 메시지 발행 함수
func Publish(conn *stomp.Conn, req MessageRequest) error {
    body, err := json.Marshal(req)
    if err != nil {
        return err
    }
    metrics.Default.Message.IncSent() // 보낸 메세지 지표 수집
    return conn.Send(
        "/pub/message",
        "application/json",
        body,
    )
}

func PublishLoopWithID(conn *stomp.Conn, roomID int, count int, clientNum int) {
    for i := 1; i <= count; i++ {
        messageId := uuid.New().String()
        req := MessageRequest{
            RoomId:  roomID,
            Content: fmt.Sprintf("Hello from client %d, message %d", clientNum, i),
            MessageId: messageId,         
            SentAt: time.Now().UnixNano(),
        }
        metrics.Default.Delivery.AddSent(messageId)
        
        err := Publish(conn, req)
        if err != nil {
            log.Printf("[클라이언트 %d] 메시지 전송 실패(%d): %v", clientNum, i, err)
            return
        }
        log.Printf("[클라이언트 %d] 메시지 전송: %d", clientNum, i)
        time.Sleep(300 * time.Millisecond) // (예시) 메시지 발송 간격
    }
}
