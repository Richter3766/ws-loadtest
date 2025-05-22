package wsclient

import (
	"net"
	"time"

	"github.com/gorilla/websocket"
)

// wsNetConn: gorilla/websocket.Conn을 net.Conn처럼 보이게 해주는 래퍼
type wsNetConn struct {
    *websocket.Conn
}

func (w *wsNetConn) Read(b []byte) (int, error) {
    _, r, err := w.NextReader()
    if err != nil {
        return 0, err
    }
    return r.Read(b)
}

func (w *wsNetConn) Write(b []byte) (int, error) {
    err := w.WriteMessage(websocket.TextMessage, b)
    if err != nil {
        return 0, err
    }
    return len(b), nil
}

func (w *wsNetConn) Close() error {
    return w.Conn.Close()
}

func (w *wsNetConn) LocalAddr() net.Addr                { return w.Conn.LocalAddr() }
func (w *wsNetConn) RemoteAddr() net.Addr               { return w.Conn.RemoteAddr() }
func (w *wsNetConn) SetDeadline(t time.Time) error      { return nil }
func (w *wsNetConn) SetReadDeadline(t time.Time) error  { return w.Conn.SetReadDeadline(t) }
func (w *wsNetConn) SetWriteDeadline(t time.Time) error { return w.Conn.SetWriteDeadline(t) }
