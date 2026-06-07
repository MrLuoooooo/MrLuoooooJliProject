package ws

import (
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// PushMessage WebSocket 推送的消息体
type PushMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// conn 包装一个 WebSocket 连接，带发送队列
type conn struct {
	userID uint
	ws     *websocket.Conn
	send   chan []byte
}

// Manager 管理所有用户的 WebSocket 连接
// 同一个用户可以从多个设备连接，每个设备一条连接
type Manager struct {
	mu      sync.RWMutex
	conns   map[uint][]*conn
	connCnt int64
}

func NewManager() *Manager {
	return &Manager{conns: make(map[uint][]*conn)}
}

// ConnCount 返回当前连接数
func (m *Manager) ConnCount() int64 {
	return atomic.LoadInt64(&m.connCnt)
}

// Add 注册一条连接
func (m *Manager) Add(userID uint, c *websocket.Conn) {
	wc := &conn{userID: userID, ws: c, send: make(chan []byte, 64)}
	m.mu.Lock()
	m.conns[userID] = append(m.conns[userID], wc)
	m.mu.Unlock()
	atomic.AddInt64(&m.connCnt, 1)
	go wc.writePump(m)
	go func() {
		wc.readPump()
		m.remove(userID, wc)
	}()
}

// SendToUser 推送给指定用户的全部连接，非阻塞
func (m *Manager) SendToUser(userID uint, msg PushMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	m.mu.RLock()
	conns := m.conns[userID]
	m.mu.RUnlock()
	for _, c := range conns {
		select {
		case c.send <- data:
		default:
			zap.S().Warn("ws推送队列已满，丢弃消息", "userID", userID)
		}
	}
}

// Broadcast 推送给所有在线用户，非阻塞
func (m *Manager) Broadcast(msg PushMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	m.mu.RLock()
	// 迭代期间不能长时间持锁，先快照 userIDs
	userIDs := make([]uint, 0, len(m.conns))
	for uid := range m.conns {
		userIDs = append(userIDs, uid)
	}
	m.mu.RUnlock()
	for _, uid := range userIDs {
		m.mu.RLock()
		conns := m.conns[uid]
		m.mu.RUnlock()
		for _, c := range conns {
			select {
			case c.send <- data:
			default:
			}
		}
	}
}

func (m *Manager) remove(userID uint, wc *conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	list := m.conns[userID]
	for i, c := range list {
		if c == wc {
			m.conns[userID] = append(list[:i], list[i+1:]...)
			break
		}
	}
	if len(m.conns[userID]) == 0 {
		delete(m.conns, userID)
	}
	atomic.AddInt64(&m.connCnt, -1)
}

func (c *conn) writePump(m *Manager) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			c.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.ws.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *conn) readPump() {
	defer c.ws.Close()
	c.ws.SetReadLimit(4096)
	c.ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	for {
		if _, _, err := c.ws.ReadMessage(); err != nil {
			break
		}
	}
}
