package main

import (
	"context"
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"amlakcrm/backend/internal/config"
	"amlakcrm/backend/internal/service"

	"nhooyr.io/websocket"
)

const shardCount = 64

type client struct {
	userID string
	conn   *websocket.Conn
	send   chan []byte
}

type shard struct {
	mu      sync.RWMutex
	clients map[*client]struct{}
}

type hub struct {
	shards [shardCount]*shard
	count  atomic.Int64
}

func newHub() *hub {
	h := &hub{}
	for i := range h.shards {
		h.shards[i] = &shard{clients: map[*client]struct{}{}}
	}
	return h
}

func (h *hub) add(c *client) {
	s := h.shards[h.index(c.userID)]
	s.mu.Lock()
	s.clients[c] = struct{}{}
	s.mu.Unlock()
	h.count.Add(1)
}

func (h *hub) remove(c *client) {
	s := h.shards[h.index(c.userID)]
	s.mu.Lock()
	delete(s.clients, c)
	s.mu.Unlock()
	h.count.Add(-1)
}

func (h *hub) publish(userID string, payload []byte) int {
	s := h.shards[h.index(userID)]
	s.mu.RLock()
	defer s.mu.RUnlock()
	delivered := 0
	for c := range s.clients {
		if c.userID != userID {
			continue
		}
		select {
		case c.send <- payload:
			delivered++
		default:
		}
	}
	return delivered
}

func (h *hub) index(userID string) uint32 {
	sum := sha1.Sum([]byte(userID))
	return binary.BigEndian.Uint32(sum[:4]) % shardCount
}

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	tokenService := service.NewTokenService(cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	h := newHub()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"status":      "ok",
			"connections": h.count.Load(),
			"shards":      shardCount,
		})
	})
	mux.HandleFunc("GET /ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(w, r, h, tokenService, logger)
	})
	mux.HandleFunc("POST /internal/publish", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID  string          `json:"userId"`
			Type    string          `json:"type"`
			Payload json.RawMessage `json:"payload"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request"})
			return
		}
		body, _ := json.Marshal(req)
		delivered := h.publish(req.UserID, body)
		writeJSON(w, http.StatusOK, map[string]int{"delivered": delivered})
	})

	addr := env("REALTIME_ADDR", ":8090")
	server := &http.Server{
		Addr:              addr,
		Handler:           withCORS(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("realtime service listening", "addr", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}

func serveWS(w http.ResponseWriter, r *http.Request, h *hub, tokens *service.TokenService, logger *slog.Logger) {
	rawToken := r.URL.Query().Get("token")
	if rawToken == "" {
		rawToken = strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	}
	userID, _, err := tokens.ParseAccess(rawToken)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		CompressionMode: websocket.CompressionDisabled,
		OriginPatterns:  []string{"*"},
	})
	if err != nil {
		return
	}
	c := &client{
		userID: userID,
		conn:   conn,
		send:   make(chan []byte, 32),
	}
	h.add(c)
	defer func() {
		h.remove(c)
		close(c.send)
		_ = conn.Close(websocket.StatusNormalClosure, "closed")
	}()

	ctx := r.Context()
	go writeLoop(ctx, c)
	readLoop(ctx, c, logger)
}

func readLoop(ctx context.Context, c *client, logger *slog.Logger) {
	for {
		typ, payload, err := c.conn.Read(ctx)
		if err != nil {
			return
		}
		if typ != websocket.MessageText {
			continue
		}
		if string(payload) == "ping" {
			select {
			case c.send <- []byte("pong"):
			default:
				logger.Warn("dropping pong for slow client", "userId", c.userID)
			}
		}
	}
}

func writeLoop(ctx context.Context, c *client) {
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.conn.Ping(ctx); err != nil {
				return
			}
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err := c.conn.Write(writeCtx, websocket.MessageText, msg)
			cancel()
			if err != nil {
				return
			}
		}
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func env(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
