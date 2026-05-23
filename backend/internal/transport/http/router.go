package httptransport

import (
	"log/slog"
	"net/http"
	"time"

	"amlakcrm/backend/internal/config"
	"amlakcrm/backend/internal/repository"
	"amlakcrm/backend/internal/service"
)

type RouterDeps struct {
	Config     config.Config
	Logger     *slog.Logger
	Store      repository.Store
	Tokens     *service.TokenService
	Auth       *service.AuthService
	Businesses *service.BusinessService
	Invites    *service.InvitationService
	Files      *service.FileService
	Uploads    *service.UploadService
	Locations  *service.LocationService
	Properties *service.PropertyService
	Contacts   *service.ContactService
	Channels   *service.ChannelService
	Matching   service.RequestMatcher
	Platform   *service.PlatformService
}

func NewRouter(deps RouterDeps) http.Handler {
	api := &api{deps: deps}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.Handle("GET /objects/", http.StripPrefix("/objects/", http.FileServer(http.Dir(deps.Config.ObjectStorageDir))))

	mux.HandleFunc("POST /api/v1/auth/request-otp", api.requestOTP)
	mux.HandleFunc("POST /api/v1/auth/verify-otp", api.verifyOTP)
	mux.HandleFunc("POST /api/v1/auth/refresh", api.refresh)
	mux.HandleFunc("POST /api/v1/auth/logout", api.logout)
	mux.HandleFunc("GET /api/v1/test/latest-otp", api.latestOTP)
	mux.HandleFunc("GET /api/v1/catalog/cities", api.catalogCities)

	protected := http.NewServeMux()
	protected.HandleFunc("GET /api/v1/catalog/cities/{cityId}/locations", api.catalogCityLocations)
	protected.HandleFunc("GET /api/v1/catalog/cities/{cityId}/locations/search", api.catalogSearchLocations)
	protected.HandleFunc("GET /api/v1/auth/me", api.me)
	protected.HandleFunc("PATCH /api/v1/auth/profile", api.profile)
	protected.HandleFunc("GET /api/v1/auth/security", api.securityProfile)
	protected.HandleFunc("PATCH /api/v1/auth/privacy", api.updatePrivacy)
	protected.HandleFunc("GET /api/v1/auth/sessions", api.sessions)
	protected.HandleFunc("DELETE /api/v1/auth/sessions/{sessionId}", api.revokeSession)
	protected.HandleFunc("GET /api/v1/notifications", api.notifications)
	protected.HandleFunc("POST /api/v1/notifications/{notificationId}/read", func(w http.ResponseWriter, r *http.Request) {
		api.markNotificationRead(w, r, r.PathValue("notificationId"))
	})
	protected.HandleFunc("POST /api/v1/businesses", api.createBusiness)
	protected.HandleFunc("GET /api/v1/businesses", api.listBusinesses)
	protected.HandleFunc("GET /api/v1/invitations/inbox", api.invitationInbox)
	protected.HandleFunc("POST /api/v1/invitations/{invitationId}/accept", api.acceptInvitation)
	protected.HandleFunc("POST /api/v1/invitations/{invitationId}/reject", api.rejectInvitation)
	protected.HandleFunc("POST /api/v1/files/business-logo", api.uploadBusinessLogo)
	protected.HandleFunc("POST /api/v1/uploads", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToFileService(w, r) {
			return
		}
		api.uploadFile(w, r)
	})
	protected.HandleFunc("POST /api/v1/location-suggestions", api.createLocationSuggestion)
	protected.HandleFunc("GET /api/v1/contact-tags", api.contactTags)
	protected.HandleFunc("GET /api/v1/channels", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToMessaging(w, r) {
			return
		}
		api.listChannels(w, r)
	})
	protected.HandleFunc("POST /api/v1/channels/private", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToMessaging(w, r) {
			return
		}
		api.createPrivateChannel(w, r)
	})
	protected.HandleFunc("GET /api/v1/vaults", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToFiling(w, r) {
			return
		}
		api.userVaults(w, r)
	})
	protected.HandleFunc("POST /api/v1/vaults", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToFiling(w, r) {
			return
		}
		api.createUserVault(w, r)
	})
	protected.HandleFunc("GET /api/v1/channels/me/main", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToMessaging(w, r) {
			return
		}
		api.userMainChannel(w, r)
	})
	protected.HandleFunc("GET /api/v1/channels/{channelId}/messages", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToMessaging(w, r) {
			return
		}
		api.channelMessages(w, r, r.PathValue("channelId"))
	})
	protected.HandleFunc("POST /api/v1/channels/{channelId}/messages", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToMessaging(w, r) {
			return
		}
		api.createChannelMessage(w, r, r.PathValue("channelId"))
	})
	protected.HandleFunc("PATCH /api/v1/channels/{channelId}/messages/{messageId}", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToMessaging(w, r) {
			return
		}
		api.updateChannelMessage(w, r, r.PathValue("channelId"), r.PathValue("messageId"))
	})
	protected.HandleFunc("DELETE /api/v1/channels/{channelId}/messages/{messageId}", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToMessaging(w, r) {
			return
		}
		api.deleteChannelMessage(w, r, r.PathValue("channelId"), r.PathValue("messageId"))
	})
	protected.HandleFunc("POST /api/v1/channels/{channelId}/media", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToMessaging(w, r) {
			return
		}
		api.uploadChannelMedia(w, r, r.PathValue("channelId"))
	})
	protected.HandleFunc("GET /api/v1/channels/{channelId}/vault/files", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToFiling(w, r) {
			return
		}
		api.channelVaultFiles(w, r, r.PathValue("channelId"))
	})
	protected.HandleFunc("POST /api/v1/channels/{channelId}/vault/files", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToFiling(w, r) {
			return
		}
		api.uploadChannelVaultFile(w, r, r.PathValue("channelId"))
	})
	protected.HandleFunc("POST /api/v1/channels/{channelId}/vault/files/confirm", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToFiling(w, r) {
			return
		}
		api.confirmChannelVaultFile(w, r, r.PathValue("channelId"))
	})
	protected.HandleFunc("GET /api/v1/channels/{channelId}/vault/files/{fileId}", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToFiling(w, r) {
			return
		}
		api.channelVaultFile(w, r, r.PathValue("channelId"), r.PathValue("fileId"))
	})
	protected.HandleFunc("GET /api/v1/channels/{channelId}/members", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToMessaging(w, r) {
			return
		}
		api.channelMembers(w, r, r.PathValue("channelId"))
	})
	protected.HandleFunc("PATCH /api/v1/channels/{channelId}/members/{memberId}", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToMessaging(w, r) {
			return
		}
		api.updateChannelMemberRole(w, r, r.PathValue("channelId"), r.PathValue("memberId"))
	})
	protected.HandleFunc("POST /api/v1/channels/{channelId}/invites", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToMessaging(w, r) {
			return
		}
		api.inviteChannelMember(w, r, r.PathValue("channelId"))
	})
	protected.HandleFunc("GET /api/v1/admin/me", api.adminMe)
	protected.HandleFunc("GET /api/v1/admin/users", api.adminUsers)
	protected.HandleFunc("GET /api/v1/admin/businesses", api.adminBusinesses)
	protected.HandleFunc("GET /api/v1/admin/accounts", api.adminAccounts)
	protected.HandleFunc("POST /api/v1/admin/accounts", api.adminSaveAccount)
	protected.HandleFunc("GET /api/v1/admin/settings", api.adminSettings)
	protected.HandleFunc("PATCH /api/v1/admin/settings", api.adminUpdateSettings)
	protected.HandleFunc("POST /api/v1/admin/cities", api.adminCreateCity)
	protected.HandleFunc("GET /api/v1/admin/location-suggestions", api.adminLocationSuggestions)
	protected.HandleFunc("POST /api/v1/admin/location-suggestions/{suggestionId}/approve", api.adminApproveLocationSuggestion)
	protected.HandleFunc("POST /api/v1/admin/location-suggestions/{suggestionId}/reject", api.adminRejectLocationSuggestion)
	protected.HandleFunc("GET /api/v1/businesses/{businessId}/dashboard", func(w http.ResponseWriter, r *http.Request) {
		api.dashboard(w, r, r.PathValue("businessId"))
	})
	protected.HandleFunc("DELETE /api/v1/businesses/{businessId}/membership", func(w http.ResponseWriter, r *http.Request) {
		api.leaveBusiness(w, r, r.PathValue("businessId"))
	})
	protected.HandleFunc("GET /api/v1/businesses/{businessId}/members", func(w http.ResponseWriter, r *http.Request) {
		api.members(w, r, r.PathValue("businessId"))
	})
	protected.HandleFunc("GET /api/v1/businesses/{businessId}/channel", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToMessaging(w, r) {
			return
		}
		api.businessMainChannel(w, r, r.PathValue("businessId"))
	})
	protected.HandleFunc("GET /api/v1/businesses/{businessId}/vaults", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToFiling(w, r) {
			return
		}
		api.businessVaults(w, r, r.PathValue("businessId"))
	})
	protected.HandleFunc("POST /api/v1/businesses/{businessId}/vaults", func(w http.ResponseWriter, r *http.Request) {
		if api.proxyToFiling(w, r) {
			return
		}
		api.createBusinessVault(w, r, r.PathValue("businessId"))
	})
	protected.HandleFunc("PATCH /api/v1/businesses/{businessId}/members/{memberId}", func(w http.ResponseWriter, r *http.Request) {
		api.updateMember(w, r, r.PathValue("businessId"), r.PathValue("memberId"))
	})
	protected.HandleFunc("POST /api/v1/businesses/{businessId}/invitations", func(w http.ResponseWriter, r *http.Request) {
		api.createInvitation(w, r, r.PathValue("businessId"))
	})
	protected.HandleFunc("GET /api/v1/businesses/{businessId}/invitations", func(w http.ResponseWriter, r *http.Request) {
		api.listInvitations(w, r, r.PathValue("businessId"))
	})
	protected.HandleFunc("GET /api/v1/businesses/{businessId}/locations", func(w http.ResponseWriter, r *http.Request) {
		api.listLocations(w, r, r.PathValue("businessId"))
	})
	protected.HandleFunc("POST /api/v1/businesses/{businessId}/areas", func(w http.ResponseWriter, r *http.Request) {
		api.createArea(w, r, r.PathValue("businessId"))
	})
	protected.HandleFunc("DELETE /api/v1/businesses/{businessId}/areas/{areaId}", func(w http.ResponseWriter, r *http.Request) {
		api.deleteArea(w, r, r.PathValue("businessId"), r.PathValue("areaId"))
	})
	protected.HandleFunc("POST /api/v1/businesses/{businessId}/areas/{areaId}/streets", func(w http.ResponseWriter, r *http.Request) {
		api.createStreet(w, r, r.PathValue("businessId"), r.PathValue("areaId"))
	})
	protected.HandleFunc("DELETE /api/v1/businesses/{businessId}/areas/{areaId}/streets/{streetId}", func(w http.ResponseWriter, r *http.Request) {
		api.deleteStreet(w, r, r.PathValue("businessId"), r.PathValue("areaId"), r.PathValue("streetId"))
	})
	protected.HandleFunc("POST /api/v1/businesses/{businessId}/areas/{areaId}/streets/{streetId}/neighborhoods", func(w http.ResponseWriter, r *http.Request) {
		api.createNeighborhood(w, r, r.PathValue("businessId"), r.PathValue("areaId"), r.PathValue("streetId"))
	})
	protected.HandleFunc("DELETE /api/v1/businesses/{businessId}/areas/{areaId}/streets/{streetId}/neighborhoods/{neighborhoodId}", func(w http.ResponseWriter, r *http.Request) {
		api.deleteNeighborhood(w, r, r.PathValue("businessId"), r.PathValue("areaId"), r.PathValue("streetId"), r.PathValue("neighborhoodId"))
	})
	protected.HandleFunc("GET /api/v1/businesses/{businessId}/properties", func(w http.ResponseWriter, r *http.Request) {
		api.listProperties(w, r, r.PathValue("businessId"))
	})
	protected.HandleFunc("GET /api/v1/businesses/{businessId}/properties/latest", func(w http.ResponseWriter, r *http.Request) {
		api.latestProperties(w, r, r.PathValue("businessId"))
	})
	protected.HandleFunc("POST /api/v1/businesses/{businessId}/properties", func(w http.ResponseWriter, r *http.Request) {
		api.createProperty(w, r, r.PathValue("businessId"))
	})
	protected.HandleFunc("PATCH /api/v1/businesses/{businessId}/properties/{propertyId}", func(w http.ResponseWriter, r *http.Request) {
		api.updateProperty(w, r, r.PathValue("businessId"), r.PathValue("propertyId"))
	})
	protected.HandleFunc("POST /api/v1/businesses/{businessId}/properties/{propertyId}/share-requests", func(w http.ResponseWriter, r *http.Request) {
		api.requestPropertyShare(w, r, r.PathValue("businessId"), r.PathValue("propertyId"))
	})
	protected.HandleFunc("GET /api/v1/businesses/{businessId}/property-share-requests", func(w http.ResponseWriter, r *http.Request) {
		api.propertyShareRequests(w, r, r.PathValue("businessId"))
	})
	protected.HandleFunc("POST /api/v1/businesses/{businessId}/property-share-requests/{requestId}/approve", func(w http.ResponseWriter, r *http.Request) {
		api.decidePropertyShare(w, r, r.PathValue("businessId"), r.PathValue("requestId"), true)
	})
	protected.HandleFunc("POST /api/v1/businesses/{businessId}/property-share-requests/{requestId}/reject", func(w http.ResponseWriter, r *http.Request) {
		api.decidePropertyShare(w, r, r.PathValue("businessId"), r.PathValue("requestId"), false)
	})
	protected.HandleFunc("POST /api/v1/businesses/{businessId}/property-share-requests/{requestId}/receive", func(w http.ResponseWriter, r *http.Request) {
		api.receivePropertyShare(w, r, r.PathValue("businessId"), r.PathValue("requestId"))
	})
	protected.HandleFunc("POST /api/v1/businesses/{businessId}/properties/{propertyId}/media", func(w http.ResponseWriter, r *http.Request) {
		api.uploadPropertyMedia(w, r, r.PathValue("businessId"), r.PathValue("propertyId"))
	})
	protected.HandleFunc("POST /api/v1/businesses/{businessId}/properties/{propertyId}/media/confirm", func(w http.ResponseWriter, r *http.Request) {
		api.confirmPropertyMedia(w, r, r.PathValue("businessId"), r.PathValue("propertyId"))
	})
	protected.HandleFunc("GET /api/v1/businesses/{businessId}/contacts", func(w http.ResponseWriter, r *http.Request) {
		api.listContacts(w, r, r.PathValue("businessId"))
	})
	protected.HandleFunc("POST /api/v1/businesses/{businessId}/contacts", func(w http.ResponseWriter, r *http.Request) {
		api.createContact(w, r, r.PathValue("businessId"))
	})
	protected.HandleFunc("POST /api/v1/businesses/{businessId}/contacts/profile-categories", func(w http.ResponseWriter, r *http.Request) {
		api.addProfileContactCategories(w, r, r.PathValue("businessId"))
	})
	protected.HandleFunc("PATCH /api/v1/businesses/{businessId}/contacts/{contactId}", func(w http.ResponseWriter, r *http.Request) {
		api.updateContact(w, r, r.PathValue("businessId"), r.PathValue("contactId"))
	})
	protected.HandleFunc("GET /api/v1/businesses/{businessId}/contacts/{contactId}/requests/{requestId}/matches", func(w http.ResponseWriter, r *http.Request) {
		api.requestMatches(w, r, r.PathValue("businessId"), r.PathValue("contactId"), r.PathValue("requestId"))
	})

	authenticated := authMiddleware(deps.Tokens, deps.Store)(protected)
	root := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/auth/request-otp" && r.URL.Path != "/api/v1/auth/verify-otp" && r.URL.Path != "/api/v1/auth/refresh" && r.URL.Path != "/api/v1/auth/logout" && r.URL.Path != "/api/v1/test/latest-otp" && r.URL.Path != "/api/v1/catalog/cities" && len(r.URL.Path) >= 8 && r.URL.Path[:8] == "/api/v1/" {
			authenticated.ServeHTTP(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	})
	return cors(recoverer(requestLogger(deps.Logger)(rateLimit(120, time.Minute)(root))))
}
