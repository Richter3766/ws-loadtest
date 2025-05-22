package wsclient

import (
	"log"

	"github.com/go-stomp/stomp"
	"github.com/gorilla/websocket"

)

// ConnectWSSTOMP: gorilla/websocket으로 WS 연결 후 STOMP 클라이언트 생성
func ConnectWSSTOMP(wsURL string, jwt string) (*stomp.Conn, *wsNetConn, error) {
    wsDialer := websocket.Dialer{}
    wsConn, _, err := wsDialer.Dial(wsURL, nil)
    if err != nil {
        log.Printf("웹소켓 연결 실패: %v", err)
        return nil, nil, err
    }
    netConn := &wsNetConn{wsConn}

    stompConn, err := stomp.Connect(
        netConn,
        
        stomp.ConnOpt.Header("Authorization", "Bearer "+jwt),
    )
    if err != nil {
        log.Printf("STOMP 연결 실패: %v", err)
        return nil, nil, err
    }

    return stompConn, netConn, nil
}

