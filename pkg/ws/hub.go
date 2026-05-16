package ws

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const redisChatPrefix = "chat:"

// Hub 管理本节点所有活跃 WebSocket 连接，通过 Redis Pub/Sub 实现跨节点广播。
// roomID 格式：两个用户 ID 按字典序拼接，如 "uid1:uid2"；群组使用 "group:{groupId}"。
type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[*Client]struct{} // roomID -> 本节点连接集合

	reg     chan *Client
	unreg   chan *Client
	localFw chan *Message // 本节点内部投递（来自 Redis 订阅）

	rdb *redis.Client
	ps  *redis.PubSub
	// 已订阅的 channel 集合，避免重复订阅
	subscribed map[string]struct{}
	subMu      sync.Mutex
}

var Default *Hub

func InitHub(rdb *redis.Client) {
	Default = newHub(rdb)
	InitNotifyHub(rdb)
}

func newHub(rdb *redis.Client) *Hub {
	h := &Hub{
		rooms:      make(map[string]map[*Client]struct{}),
		reg:        make(chan *Client, 64),
		unreg:      make(chan *Client, 64),
		localFw:    make(chan *Message, 512),
		rdb:        rdb,
		ps:         rdb.Subscribe(context.Background()), // 空订阅，后续动态加
		subscribed: make(map[string]struct{}),
	}
	go h.run()
	go h.redisReceive()
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
			h.ensureSubscribed(c.roomID)

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

		case msg := <-h.localFw:
			h.deliverLocal(msg)
		}
	}
}

// redisReceive 持续从 Redis Pub/Sub 读取消息并转发给本地连接。
func (h *Hub) redisReceive() {
	ch := h.ps.Channel()
	for redisMsg := range ch {
		var msg Message
		if err := json.Unmarshal([]byte(redisMsg.Payload), &msg); err != nil {
			zap.S().Warnf("ws: failed to unmarshal redis message: %v", err)
			continue
		}
		h.localFw <- &msg
	}
}

func (h *Hub) deliverLocal(msg *Message) {
	h.mu.RLock()
	room := h.rooms[msg.RoomID]
	h.mu.RUnlock()
	for c := range room {
		select {
		case c.send <- msg:
		default:
			h.unreg <- c
		}
	}
}

// ensureSubscribed 确保本节点已订阅该 roomID 对应的 Redis channel。
func (h *Hub) ensureSubscribed(roomID string) {
	channel := redisChatPrefix + roomID
	h.subMu.Lock()
	defer h.subMu.Unlock()
	if _, ok := h.subscribed[channel]; ok {
		return
	}
	if err := h.ps.Subscribe(context.Background(), channel); err != nil {
		zap.S().Warnf("ws: redis subscribe %s error: %v", channel, err)
		return
	}
	h.subscribed[channel] = struct{}{}
}

// Forward 将消息发布到 Redis，由各节点的订阅者负责投递给本地连接。
func (h *Hub) Forward(msg *Message) {
	b, err := json.Marshal(msg)
	if err != nil {
		return
	}
	channel := redisChatPrefix + msg.RoomID
	if err = h.rdb.Publish(context.Background(), channel, b).Err(); err != nil {
		zap.S().Warnf("ws: redis publish %s error: %v", channel, err)
	}
}

func (h *Hub) Register(c *Client)   { h.reg <- c }
func (h *Hub) Unregister(c *Client) { h.unreg <- c }
