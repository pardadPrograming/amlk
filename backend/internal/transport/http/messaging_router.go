package httptransport

import (
	"net/http"
	"time"
)

func NewMessagingRouter(deps RouterDeps) http.Handler {
	api := &api{deps: deps}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	protected := http.NewServeMux()
	protected.HandleFunc("GET /api/v1/channels", api.listChannels)
	protected.HandleFunc("POST /api/v1/channels/private", api.createPrivateChannel)
	protected.HandleFunc("GET /api/v1/channels/me/main", api.userMainChannel)
	protected.HandleFunc("GET /api/v1/channels/{channelId}/messages", func(w http.ResponseWriter, r *http.Request) {
		api.channelMessages(w, r, r.PathValue("channelId"))
	})
	protected.HandleFunc("POST /api/v1/channels/{channelId}/messages", func(w http.ResponseWriter, r *http.Request) {
		api.createChannelMessage(w, r, r.PathValue("channelId"))
	})
	protected.HandleFunc("PATCH /api/v1/channels/{channelId}/messages/{messageId}", func(w http.ResponseWriter, r *http.Request) {
		api.updateChannelMessage(w, r, r.PathValue("channelId"), r.PathValue("messageId"))
	})
	protected.HandleFunc("DELETE /api/v1/channels/{channelId}/messages/{messageId}", func(w http.ResponseWriter, r *http.Request) {
		api.deleteChannelMessage(w, r, r.PathValue("channelId"), r.PathValue("messageId"))
	})
	protected.HandleFunc("POST /api/v1/channels/{channelId}/media", func(w http.ResponseWriter, r *http.Request) {
		api.uploadChannelMedia(w, r, r.PathValue("channelId"))
	})
	protected.HandleFunc("GET /api/v1/channels/{channelId}/members", func(w http.ResponseWriter, r *http.Request) {
		api.channelMembers(w, r, r.PathValue("channelId"))
	})
	protected.HandleFunc("PATCH /api/v1/channels/{channelId}/members/{memberId}", func(w http.ResponseWriter, r *http.Request) {
		api.updateChannelMemberRole(w, r, r.PathValue("channelId"), r.PathValue("memberId"))
	})
	protected.HandleFunc("POST /api/v1/channels/{channelId}/invites", func(w http.ResponseWriter, r *http.Request) {
		api.inviteChannelMember(w, r, r.PathValue("channelId"))
	})
	protected.HandleFunc("GET /api/v1/businesses/{businessId}/channel", func(w http.ResponseWriter, r *http.Request) {
		api.businessMainChannel(w, r, r.PathValue("businessId"))
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
