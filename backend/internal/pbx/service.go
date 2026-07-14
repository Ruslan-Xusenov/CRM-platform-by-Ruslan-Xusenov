package pbx

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/crm-platform/backend/internal/shared/middleware"
	"github.com/crm-platform/backend/pkg/broker"
	"github.com/crm-platform/backend/pkg/cache"
	"github.com/crm-platform/backend/pkg/ws"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ActiveCall represents an in-progress call tracked in memory.
type ActiveCall struct {
	ID        string    `json:"id"`
	ChannelID string    `json:"channel_id"`
	Caller    string    `json:"caller"`
	Callee    string    `json:"callee"`
	Direction string    `json:"direction"`
	Status    string    `json:"status"`
	StartedAt time.Time `json:"started_at"`
	TenantID  string    `json:"tenant_id"`
}

type Service struct {
	db          *pgxpool.Pool
	redis       *cache.RedisClient
	mq          *broker.RabbitMQ
	hub         *ws.Hub
	ariURL      string
	ariUser     string
	ariPassword string
	activeCalls map[string]*ActiveCall
	mu          sync.RWMutex
}

func NewService(db *pgxpool.Pool, redis *cache.RedisClient, mq *broker.RabbitMQ, hub *ws.Hub, ariURL, ariUser, ariPass string) *Service {
	return &Service{
		db: db, redis: redis, mq: mq, hub: hub,
		ariURL: ariURL, ariUser: ariUser, ariPassword: ariPass,
		activeCalls: make(map[string]*ActiveCall),
	}
}

// ConnectARI connects to the Asterisk ARI WebSocket and listens for events.
func (s *Service) ConnectARI() {
	wsURL := fmt.Sprintf("%s/events?api_key=%s:%s&app=crm-pbx",
		s.ariURL, s.ariUser, s.ariPassword)
	wsURL = "ws" + wsURL[4:] // Convert http:// to ws://

	for {
		slog.Info("Connecting to Asterisk ARI WebSocket", "url", wsURL)
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			slog.Error("ARI WebSocket connection failed", "error", err)
			time.Sleep(5 * time.Second)
			continue
		}

		slog.Info("✅ Connected to Asterisk ARI")
		s.listenARI(conn)
		conn.Close()
		slog.Warn("ARI WebSocket disconnected, reconnecting...")
		time.Sleep(2 * time.Second)
	}
}

func (s *Service) listenARI(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			slog.Error("ARI read error", "error", err)
			return
		}

		var event map[string]interface{}
		if err := json.Unmarshal(message, &event); err != nil {
			continue
		}

		eventType, _ := event["type"].(string)
		slog.Debug("ARI event", "type", eventType)

		switch eventType {
		case "StasisStart":
			s.handleStasisStart(event)
		case "StasisEnd":
			s.handleStasisEnd(event)
		case "ChannelStateChange":
			s.handleChannelState(event)
		}
	}
}

func (s *Service) handleStasisStart(event map[string]interface{}) {
	channel, _ := event["channel"].(map[string]interface{})
	channelID, _ := channel["id"].(string)
	caller, _ := channel["caller"].(map[string]interface{})
	callerNum, _ := caller["number"].(string)

	args, _ := event["args"].([]interface{})
	callee := ""
	if len(args) > 0 {
		callee, _ = args[0].(string)
	}

	call := &ActiveCall{
		ID:        uuid.NewString(),
		ChannelID: channelID,
		Caller:    callerNum,
		Callee:    callee,
		Direction: "inbound",
		Status:    "ringing",
		StartedAt: time.Now(),
	}

	s.mu.Lock()
	s.activeCalls[channelID] = call
	s.mu.Unlock()

	slog.Info("Call started", "caller", callerNum, "callee", callee, "channel", channelID)

	// Broadcast incoming call to all connected agents
	s.hub.BroadcastToTenant(call.TenantID, ws.Message{
		Type:    "call.incoming",
		Payload: call,
	})
}

func (s *Service) handleStasisEnd(event map[string]interface{}) {
	channel, _ := event["channel"].(map[string]interface{})
	channelID, _ := channel["id"].(string)

	s.mu.Lock()
	call, exists := s.activeCalls[channelID]
	if exists {
		delete(s.activeCalls, channelID)
	}
	s.mu.Unlock()

	if !exists {
		return
	}

	// Save CDR to database
	duration := int(time.Since(call.StartedAt).Seconds())
	ctx := context.Background()

	_, err := s.db.Exec(ctx,
		`INSERT INTO call_logs (tenant_id, caller, callee, direction, status, started_at, ended_at, duration_seconds, asterisk_unique_id)
		 VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7, $8)`,
		call.TenantID, call.Caller, call.Callee, call.Direction, "completed",
		call.StartedAt, duration, channelID)
	if err != nil {
		slog.Error("Failed to save CDR", "error", err)
	}

	slog.Info("Call ended", "caller", call.Caller, "duration", duration)

	s.hub.BroadcastToTenant(call.TenantID, ws.Message{
		Type:    "call.ended",
		Payload: map[string]interface{}{"channel_id": channelID, "duration": duration},
	})
}

func (s *Service) handleChannelState(event map[string]interface{}) {
	// Track channel state changes (ringing -> answered, etc.)
}

// ─── API Methods ─────────────────────────────────────────────

func (s *Service) OriginateCall(ctx context.Context, from, to, tenantID string) error {
	// Use ARI REST API to originate a call
	url := fmt.Sprintf("%s/channels?endpoint=PJSIP/%s&extension=%s&context=crm-stasis&priority=1&api_key=%s:%s",
		s.ariURL, to, to, s.ariUser, s.ariPassword)

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("originate call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("ARI originate failed: status %d", resp.StatusCode)
	}

	slog.Info("Call originated", "from", from, "to", to)
	return nil
}

func (s *Service) GetActiveCalls() []*ActiveCall {
	s.mu.RLock()
	defer s.mu.RUnlock()
	calls := make([]*ActiveCall, 0, len(s.activeCalls))
	for _, c := range s.activeCalls {
		calls = append(calls, c)
	}
	return calls
}

func (s *Service) ListCallHistory(ctx context.Context, page, pageSize int) (interface{}, error) {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	if pageSize == 0 { pageSize = 50 }
	offset := (page - 1) * pageSize

	rows, err := s.db.Query(ctx,
		`SELECT id, tenant_id, caller, callee, direction, status, started_at, ended_at, duration_seconds, recording_url
		 FROM call_logs WHERE tenant_id=$1 ORDER BY started_at DESC LIMIT $2 OFFSET $3`,
		tid, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type CDR struct {
		ID        uuid.UUID  `json:"id"`
		TenantID  string     `json:"tenant_id"`
		Caller    string     `json:"caller"`
		Callee    string     `json:"callee"`
		Direction string     `json:"direction"`
		Status    string     `json:"status"`
		StartedAt time.Time  `json:"started_at"`
		EndedAt   *time.Time `json:"ended_at"`
		Duration  int        `json:"duration_seconds"`
		Recording *string    `json:"recording_url"`
	}

	var records []CDR
	for rows.Next() {
		var r CDR
		rows.Scan(&r.ID, &r.TenantID, &r.Caller, &r.Callee, &r.Direction, &r.Status,
			&r.StartedAt, &r.EndedAt, &r.Duration, &r.Recording)
		records = append(records, r)
	}
	return records, nil
}

func (s *Service) ListExtensions(ctx context.Context) (interface{}, error) {
	tid, _ := ctx.Value(middleware.TenantIDKey).(string)
	rows, err := s.db.Query(ctx,
		`SELECT id, extension_number, display_name, user_id, enabled FROM extensions WHERE tenant_id=$1`, tid)
	if err != nil { return nil, err }
	defer rows.Close()

	type Ext struct {
		ID       uuid.UUID  `json:"id"`
		Number   string     `json:"extension_number"`
		Name     *string    `json:"display_name"`
		UserID   *uuid.UUID `json:"user_id"`
		Enabled  bool       `json:"enabled"`
	}
	var exts []Ext
	for rows.Next() {
		var e Ext
		rows.Scan(&e.ID, &e.Number, &e.Name, &e.UserID, &e.Enabled)
		exts = append(exts, e)
	}
	return exts, nil
}

// StartRecordingConsumer listens for recording upload events from RabbitMQ.
func (s *Service) StartRecordingConsumer() {
	msgs, err := s.mq.Consume("recording-processor", "recording.events", "recording.uploaded")
	if err != nil {
		slog.Error("Failed to start recording consumer", "error", err)
		return
	}

	for msg := range msgs {
		var event struct {
			UniqueID string `json:"unique_id"`
			S3Bucket string `json:"s3_bucket"`
			S3Key    string `json:"s3_key"`
		}
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			msg.Nack(false, false)
			continue
		}

		recordingURL := fmt.Sprintf("s3://%s/%s", event.S3Bucket, event.S3Key)
		_, err := s.db.Exec(context.Background(),
			`UPDATE call_logs SET recording_url = $1 WHERE asterisk_unique_id = $2`,
			recordingURL, event.UniqueID)
		if err != nil {
			slog.Error("Failed to link recording", "error", err)
			msg.Nack(false, true)
			continue
		}

		slog.Info("Recording linked to CDR", "unique_id", event.UniqueID, "url", recordingURL)
		msg.Ack(false)
	}
}
