package wsclient

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-stomp/stomp"
	"github.com/google/uuid"

)

// Subscribe : 메시지 구독 함수 (비동기로 수신 콜백)
func Subscribe(conn *stomp.Conn, roomID int, clientNum int, handleMsg func(string, int), done <-chan struct{}, subWg *sync.WaitGroup) (error) {
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
                    return
                }
                if msg.Err != nil {
                    log.Printf("메시지 수신 오류: %v", msg.Err)
                    continue
                }
                handleMsg(string(msg.Body), clientNum)
            }
        }
    }()
    return nil
}


// Publish : 메시지 발행 함수
func Publish(conn *stomp.Conn, req MessageRequest) error {
    body, err := json.Marshal(req)
    if err != nil {
        return err
    }
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
        
        err := Publish(conn, req)
        if err != nil {
            log.Printf("[클라이언트 %d] 메시지 전송 실패(%d): %v", clientNum, i, err)
            return
        }
        log.Printf("[클라이언트 %d] 메시지 전송: %d", clientNum, i)
        time.Sleep(300 * time.Millisecond) // (예시) 메시지 발송 간격
    }
}
