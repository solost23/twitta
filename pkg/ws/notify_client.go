package ws

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Notification 是推送给用户的通知体。
type Notification struct {
	ToUserID  string    `json:"toUserId"`
	FromID    string    `json:"fromId"`
	FromName  string    `json:"fromName"`
	RoomID    string    `json:"roomId"`
	Content   string    `json:"content"`
	Type      string    `json:"type"` // "message" | "read_receipt"
	CreatedAt time.Time `json:"createdAt"`
}

// NotifyClient 代表一个通知 WebSocket 连接。
type NotifyClient struct {
	hub    *NotifyHub
	conn   *websocket.Conn
	send   chan *Notification
	userID string
}

func NewNotifyClient(hub *NotifyHub, conn *websocket.Conn, userID string) *NotifyClient {
	return &NotifyClient{
		hub:    hub,
		conn:   conn,
		send:   make(chan *Notification, 32),
		userID: userID,
	}
}

// SendNotification 向该客户端发送一条通知，非阻塞。
func (c *NotifyClient) SendNotification(n *Notification) {
	select {
	case c.send <- n:
	default:
	}
}

// WritePump 将通知写回 WebSocket，并定期发送 ping 保活。
func (c *NotifyClient) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.hub.Unregister(c)
		_ = c.conn.Close()
	}()

	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		select {
		case n, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			b, err := json.Marshal(n)
			if err != nil {
				continue
			}
			if err = c.conn.WriteMessage(websocket.TextMessage, b); err != nil {
				zap.S().Warnf("notify write error: %v", err)
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
