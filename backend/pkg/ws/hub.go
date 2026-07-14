package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/crm-platform/backend/pkg/cache"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Message represents a WebSocket message.
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Client represents a connected WebSocket client.
type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   string
	TenantID string
}

// Hub manages all WebSocket connections and broadcasts.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	redis      *cache.RedisClient
	mu         sync.RWMutex
}

// NewHub creates a new WebSocket hub with Redis Pub/Sub for scaling.
func NewHub(redis *cache.RedisClient) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		redis:      redis,
	}
}

// Run starts the hub event loop.
func (h *Hub) Run() {
	// Subscribe to Redis Pub/Sub for horizontal scaling
	go h.subscribeRedis()

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			slog.Debug("WebSocket client connected", "user", client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			slog.Debug("WebSocket client disconnected", "user", client.UserID)

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// subscribeRedis listens for messages on the "ws:broadcast" channel.
func (h *Hub) subscribeRedis() {
	if h.redis == nil {
		return
	}
	ctx := context.Background()
	sub := h.redis.Subscribe(ctx, "ws:broadcast")
	defer sub.Close()

	ch := sub.Channel()
	for msg := range ch {
		h.broadcast <- []byte(msg.Payload)
	}
}

// BroadcastToTenant sends a message to all clients of a specific tenant.
func (h *Hub) BroadcastToTenant(tenantID string, msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		slog.Error("Failed to marshal WS message", "error", err)
		return
	}
	// Publish to Redis for cross-instance broadcast
	if h.redis != nil {
		h.redis.Publish(context.Background(), "ws:broadcast", string(data))
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		if client.TenantID == tenantID {
			select {
			case client.Send <- data:
			default:
			}
		}
	}
}

// BroadcastToUser sends a message to a specific user.
func (h *Hub) BroadcastToUser(userID string, msg Message) {
	data, _ := json.Marshal(msg)
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		if client.UserID == userID {
			select {
			case client.Send <- data:
			default:
			}
		}
	}
}

// HandleWebSocket is the HTTP handler for WebSocket upgrade.
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("WebSocket upgrade failed", "error", err)
		return
	}

	client := &Client{
		Hub:      h,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserID:   r.URL.Query().Get("user_id"),
		TenantID: r.URL.Query().Get("tenant_id"),
	}

	h.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *Client) writePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}
