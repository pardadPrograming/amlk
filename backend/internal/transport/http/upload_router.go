package httptransport

import (
	"net/http"
	"time"
)

func NewUploadRouter(deps RouterDeps) http.Handler {
	api := &api{deps: deps}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.Handle("GET /objects/", http.StripPrefix("/objects/", http.FileServer(http.Dir(deps.Config.ObjectStorageDir))))

	protected := http.NewServeMux()
	protected.HandleFunc("POST /api/v1/uploads", api.uploadFile)

	authenticated := authMiddleware(deps.Tokens, deps.Store)(protected)
	root := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) >= 8 && r.URL.Path[:8] == "/api/v1/" {
			authenticated.ServeHTTP(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	})
	return cors(recoverer(requestLogger(deps.Logger)(rateLimit(120, time.Minute)(root))))
}
