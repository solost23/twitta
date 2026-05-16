package ws

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMsgSize = 4096
)

// Message 是在房间内广播的消息体。
type Message struct {
	RoomID    string    `json:"roomId"`
	FromID    string    `json:"fromId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

// Client 代表一个 WebSocket 连接。
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan *Message
	roomID string
	userID string
}

func NewClient(hub *Hub, conn *websocket.Conn, roomID, userID string) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan *Message, 32),
		roomID: roomID,
		userID: userID,
	}
}

// ReadPump 从 WebSocket 读取消息并转发到 Hub。
func (c *Client) ReadPump(onMessage func(*Message)) {
	defer func() {
		c.hub.Unregister(c)
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMsgSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zap.S().Warnf("ws read error: %v", err)
			}
			break
		}

		var payload struct {
			Content string `json:"content"`
		}
		if err = json.Unmarshal(raw, &payload); err != nil || payload.Content == "" {
			continue
		}

		msg := &Message{
			RoomID:    c.roomID,
			FromID:    c.userID,
			Content:   payload.Content,
			CreatedAt: time.Now(),
		}
		onMessage(msg)
		c.hub.Forward(msg)
	}
}

// WritePump 将 Hub 推送的消息写回 WebSocket。
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			b, _ := json.Marshal(msg)
			if err := c.conn.WriteMessage(websocket.TextMessage, b); err != nil {
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
