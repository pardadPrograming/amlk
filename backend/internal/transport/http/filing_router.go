package httptransport

import (
	"net/http"
	"time"
)

func NewFilingRouter(deps RouterDeps) http.Handler {
	api := &api{deps: deps}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.Handle("GET /objects/", http.StripPrefix("/objects/", http.FileServer(http.Dir(deps.Config.ObjectStorageDir))))

	protected := http.NewServeMux()
	protected.HandleFunc("GET /api/v1/vaults", api.userVaults)
	protected.HandleFunc("POST /api/v1/vaults", api.createUserVault)
	protected.HandleFunc("GET /api/v1/channels/{channelId}/vault/files", func(w http.ResponseWriter, r *http.Request) {
		api.channelVaultFiles(w, r, r.PathValue("channelId"))
	})
	protected.HandleFunc("POST /api/v1/channels/{channelId}/vault/files", func(w http.ResponseWriter, r *http.Request) {
		api.uploadChannelVaultFile(w, r, r.PathValue("channelId"))
	})
	protected.HandleFunc("POST /api/v1/channels/{channelId}/vault/files/confirm", func(w http.ResponseWriter, r *http.Request) {
		api.confirmChannelVaultFile(w, r, r.PathValue("channelId"))
	})
	protected.HandleFunc("GET /api/v1/channels/{channelId}/vault/files/{fileId}", func(w http.ResponseWriter, r *http.Request) {
		api.channelVaultFile(w, r, r.PathValue("channelId"), r.PathValue("fileId"))
	})
	protected.HandleFunc("GET /api/v1/businesses/{businessId}/vaults", func(w http.ResponseWriter, r *http.Request) {
		api.businessVaults(w, r, r.PathValue("businessId"))
	})
	protected.HandleFunc("POST /api/v1/businesses/{businessId}/vaults", func(w http.ResponseWriter, r *http.Request) {
		api.createBusinessVault(w, r, r.PathValue("businessId"))
	})

	authenticated := authMiddleware(deps.Tokens, deps.Store)(protected)
	root := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) >= 8 && r.URL.Path[:8] == "/api/v1/" {
			authenticated.ServeHTTP(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	})
	return cors(recoverer(requestLogger(deps.Logger)(rateLimit(240, time.Minute)(root))))
}
