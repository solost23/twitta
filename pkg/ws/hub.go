package ws

import (
	"sync"
)

// Hub 管理所有活跃的 WebSocket 连接，按房间（roomID）分组。
// roomID 格式：两个用户 ID 按字典序拼接，如 "uid1:uid2"。
type Hub struct {
	mu      sync.RWMutex
	rooms   map[string]map[*Client]struct{}
	reg     chan *Client
	unreg   chan *Client
	forward chan *Message
}

var Default = NewHub()

func NewHub() *Hub {
	h := &Hub{
		rooms:   make(map[string]map[*Client]struct{}),
		reg:     make(chan *Client, 64),
		unreg:   make(chan *Client, 64),
		forward: make(chan *Message, 256),
	}
	go h.run()
	return h
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.reg:
			h.mu.Lock()
			if h.rooms[c.roomID] == nil {
				h.rooms[c.roomID] = make(map[*Client]struct{})
			}
			h.rooms[c.roomID][c] = struct{}{}
			h.mu.Unlock()

		case c := <-h.unreg:
			h.mu.Lock()
			if room, ok := h.rooms[c.roomID]; ok {
				delete(room, c)
				if len(room) == 0 {
					delete(h.rooms, c.roomID)
				}
			}
			h.mu.Unlock()
			close(c.send)

		case msg := <-h.forward:
			h.mu.RLock()
			room := h.rooms[msg.RoomID]
			h.mu.RUnlock()
			for c := range room {
				select {
				case c.send <- msg:
				default:
					// 发送缓冲满，断开该连接
					h.unreg <- c
				}
			}
		}
	}
}

func (h *Hub) Register(c *Client)   { h.reg <- c }
func (h *Hub) Unregister(c *Client) { h.unreg <- c }
func (h *Hub) Forward(m *Message)   { h.forward <- m }
