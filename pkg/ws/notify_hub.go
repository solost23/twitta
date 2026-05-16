package ws

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const redisNotifyPrefix = "notify:"

// NotifyHub 管理每个用户的全局通知 WebSocket 连接。
// 用户登录后建立一条长连接，收到新消息时通过 Redis Pub/Sub 推送通知。
type NotifyHub struct {
	mu    sync.RWMutex
	conns map[string]map[*NotifyClient]struct{} // userID -> 连接集合

	reg     chan *NotifyClient
	unreg   chan *NotifyClient
	localFw chan *Notification

	rdb        *redis.Client
	ps         *redis.PubSub
	subscribed map[string]struct{}
	subMu      sync.Mutex
}

var DefaultNotify *NotifyHub

func InitNotifyHub(rdb *redis.Client) {
	DefaultNotify = newNotifyHub(rdb)
}

func newNotifyHub(rdb *redis.Client) *NotifyHub {
	h := &NotifyHub{
		conns:      make(map[string]map[*NotifyClient]struct{}),
		reg:        make(chan *NotifyClient, 64),
		unreg:      make(chan *NotifyClient, 64),
		localFw:    make(chan *Notification, 512),
		rdb:        rdb,
		ps:         rdb.Subscribe(context.Background()),
		subscribed: make(map[string]struct{}),
	}
	go h.run()
	go h.redisReceive()
	return h
}

func (h *NotifyHub) run() {
	for {
		select {
		case c := <-h.reg:
			h.mu.Lock()
			if h.conns[c.userID] == nil {
				h.conns[c.userID] = make(map[*NotifyClient]struct{})
			}
			h.conns[c.userID][c] = struct{}{}
			h.mu.Unlock()
			h.ensureSubscribed(c.userID)

		case c := <-h.unreg:
			h.mu.Lock()
			if conns, ok := h.conns[c.userID]; ok {
				delete(conns, c)
				if len(conns) == 0 {
					delete(h.conns, c.userID)
				}
			}
			h.mu.Unlock()
			close(c.send)

		case n := <-h.localFw:
			h.deliverLocal(n)
		}
	}
}

func (h *NotifyHub) redisReceive() {
	ch := h.ps.Channel()
	for redisMsg := range ch {
		var n Notification
		if err := json.Unmarshal([]byte(redisMsg.Payload), &n); err != nil {
			zap.S().Warnf("notify: failed to unmarshal: %v", err)
			continue
		}
		h.localFw <- &n
	}
}

func (h *NotifyHub) deliverLocal(n *Notification) {
	h.mu.RLock()
	conns := h.conns[n.ToUserID]
	h.mu.RUnlock()
	for c := range conns {
		select {
		case c.send <- n:
		default:
			h.unreg <- c
		}
	}
}

func (h *NotifyHub) ensureSubscribed(userID string) {
	channel := redisNotifyPrefix + userID
	h.subMu.Lock()
	defer h.subMu.Unlock()
	if _, ok := h.subscribed[channel]; ok {
		return
	}
	if err := h.ps.Subscribe(context.Background(), channel); err != nil {
		zap.S().Warnf("notify: redis subscribe %s error: %v", channel, err)
		return
	}
	h.subscribed[channel] = struct{}{}
}

// Notify 向指定用户推送通知（通过 Redis，跨节点有效）。
func (h *NotifyHub) Notify(n *Notification) {
	b, err := json.Marshal(n)
	if err != nil {
		return
	}
	channel := redisNotifyPrefix + n.ToUserID
	if err = h.rdb.Publish(context.Background(), channel, b).Err(); err != nil {
		zap.S().Warnf("notify: redis publish %s error: %v", channel, err)
	}
}

func (h *NotifyHub) Register(c *NotifyClient)   { h.reg <- c }
func (h *NotifyHub) Unregister(c *NotifyClient) { h.unreg <- c }
