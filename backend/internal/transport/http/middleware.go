package httptransport

import (
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"amlakcrm/backend/internal/repository"
	"amlakcrm/backend/internal/service"
)

type rateBucket struct {
	count int
	reset time.Time
}

func requestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("request", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(start).Milliseconds())
		})
	}
}

func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				writeError(w, http.StatusInternalServerError, "internal_error", "خطای داخلی سرور")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func rateLimit(limit int, window time.Duration) func(http.Handler) http.Handler {
	var mu sync.Mutex
	buckets := map[string]rateBucket{}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			key := ip + ":" + r.URL.Path
			mu.Lock()
			b := buckets[key]
			if time.Now().After(b.reset) {
				b = rateBucket{reset: time.Now().Add(window)}
			}
			b.count++
			buckets[key] = b
			mu.Unlock()
			if b.count > limit {
				writeError(w, http.StatusTooManyRequests, "rate_limited", "تعداد درخواست‌ها بیش از حد مجاز است")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func authMiddleware(tokens *service.TokenService, store repository.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "unauthorized", "نیاز به ورود دارید")
				return
			}
			userID, sessionID, err := tokens.ParseAccess(strings.TrimPrefix(header, "Bearer "))
			if err != nil {
				writeError(w, http.StatusUnauthorized, "unauthorized", "نشست معتبر نیست")
				return
			}
			if sessionID != "" {
				if _, err := store.GetSession(r.Context(), sessionID); err != nil {
					writeError(w, http.StatusUnauthorized, "unauthorized", "نشست معتبر نیست")
					return
				}
				_ = store.TouchSession(r.Context(), sessionID, time.Now().UTC())
			}
			user, err := store.GetUser(r.Context(), userID)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "unauthorized", "کاربر یافت نشد")
				return
			}
			ctx := withSessionID(withUser(r.Context(), user), sessionID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
