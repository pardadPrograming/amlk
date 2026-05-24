package httptransport

import (
	"net/http"
	"strconv"

	"amlakcrm/backend/internal/domain"
	"amlakcrm/backend/internal/service"
)

type api struct {
	deps RouterDeps
}

func (a *api) requestOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone string `json:"phone"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "درخواست معتبر نیست")
		return
	}
	challenge, err := a.deps.Auth.RequestOTP(r.Context(), req.Phone)
	if err != nil {
		writeError(w, http.StatusBadRequest, "otp_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"phone":           challenge.Phone,
		"expiresAt":       challenge.ExpiresAt,
		"developmentCode": challenge.Development,
	})
}

func (a *api) verifyOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "درخواست معتبر نیست")
		return
	}
	tokens, err := a.deps.Auth.VerifyOTP(r.Context(), req.Phone, req.Code, r.UserAgent(), r.RemoteAddr)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "otp_invalid", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tokens)
}

func (a *api) latestOTP(w http.ResponseWriter, r *http.Request) {
	if a.deps.Config.AppEnv != "development" && a.deps.Config.AppEnv != "test" {
		writeError(w, http.StatusNotFound, "not_found", "مسیر یافت نشد")
		return
	}
	challenge, err := a.deps.Auth.LatestOTP(r.Context())
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "هنوز کدی ارسال نشده است")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"phone":     challenge.Phone,
		"code":      challenge.Development,
		"expiresAt": challenge.ExpiresAt,
		"sentAt":    challenge.LastSentAt,
	})
}

func (a *api) refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "درخواست معتبر نیست")
		return
	}
	tokens, err := a.deps.Auth.Refresh(r.Context(), req.RefreshToken, r.UserAgent(), r.RemoteAddr)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid_session", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tokens)
}

func (a *api) logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	_ = decodeJSON(r, &req)
	_ = a.deps.Auth.Logout(r.Context(), req.RefreshToken)
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *api) me(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	businesses, _ := a.deps.Businesses.ListForUser(r.Context(), user.ID)
	inbox, _ := a.deps.Invites.Inbox(r.Context(), user.Phone)
	writeJSON(w, http.StatusOK, map[string]interface{}{"user": user, "businesses": businesses, "invitations": inbox})
}

func (a *api) profile(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	var req struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		CityID    string `json:"cityId"`
	}
	if err := decodeJSON(r, &req); err != nil || req.FirstName == "" || req.LastName == "" || req.CityID == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "نام، نام خانوادگی و شهر الزامی است")
		return
	}
	updated, err := a.deps.Auth.CompleteProfile(r.Context(), user.ID, req.FirstName, req.LastName, req.CityID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "profile_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (a *api) securityProfile(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	data, err := a.deps.Auth.SecurityProfile(r.Context(), user.ID, currentSessionID(r.Context()))
	if err != nil {
		writeError(w, http.StatusBadRequest, "security_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (a *api) updatePrivacy(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	var req domain.PrivacySettings
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "درخواست معتبر نیست")
		return
	}
	settings, err := a.deps.Auth.UpdatePrivacy(r.Context(), user.ID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "privacy_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

func (a *api) sessions(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	items, err := a.deps.Auth.ActiveSessions(r.Context(), user.ID, currentSessionID(r.Context()))
	if err != nil {
		writeError(w, http.StatusBadRequest, "sessions_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) revokeSession(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	sessionID := r.PathValue("sessionId")
	if sessionID == currentSessionID(r.Context()) {
		writeError(w, http.StatusBadRequest, "current_session", "برای خروج از نشست فعلی از گزینه خروج استفاده کنید")
		return
	}
	if err := a.deps.Auth.RevokeUserSession(r.Context(), user.ID, sessionID); err != nil {
		writeError(w, http.StatusNotFound, "session_not_found", "نشست پیدا نشد")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *api) notifications(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	unreadOnly := r.URL.Query().Get("unread") == "true"
	items, err := a.deps.Store.ListNotifications(r.Context(), user.ID, unreadOnly)
	if err != nil {
		writeError(w, http.StatusBadRequest, "notifications_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) markNotificationRead(w http.ResponseWriter, r *http.Request, notificationID string) {
	user, _ := currentUser(r.Context())
	if err := a.deps.Store.MarkNotificationRead(r.Context(), user.ID, notificationID); err != nil {
		writeError(w, http.StatusNotFound, "notification_not_found", "اعلان پیدا نشد")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *api) createBusiness(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	var req domain.Business
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "اطلاعات کسب‌وکار معتبر نیست")
		return
	}
	business, err := a.deps.Businesses.Create(r.Context(), user, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "business_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, business)
}

func (a *api) listBusinesses(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	items, err := a.deps.Businesses.ListForUser(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "business_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) leaveBusiness(w http.ResponseWriter, r *http.Request, businessID string) {
	user, _ := currentUser(r.Context())
	if err := a.deps.Businesses.Leave(r.Context(), user.ID, businessID); err != nil {
		writeError(w, http.StatusForbidden, "leave_business_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *api) dashboard(w http.ResponseWriter, r *http.Request, businessID string) {
	data, err := a.deps.Businesses.Dashboard(r.Context(), businessID)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "کسب‌وکار یافت نشد")
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (a *api) members(w http.ResponseWriter, r *http.Request, businessID string) {
	items, err := a.deps.Businesses.Members(r.Context(), businessID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "members_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) updateMember(w http.ResponseWriter, r *http.Request, businessID, memberID string) {
	user, _ := currentUser(r.Context())
	var req struct {
		Role              domain.Role         `json:"role"`
		CommissionPercent *float64            `json:"commissionPercent"`
		Status            domain.MemberStatus `json:"status"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "درخواست معتبر نیست")
		return
	}
	commission := -1.0
	if req.CommissionPercent != nil {
		commission = *req.CommissionPercent
	}
	member, err := a.deps.Businesses.UpdateMember(r.Context(), user.ID, businessID, memberID, req.Role, commission, req.Status)
	if err != nil {
		writeError(w, http.StatusForbidden, "forbidden", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, member)
}

func (a *api) createInvitation(w http.ResponseWriter, r *http.Request, businessID string) {
	user, _ := currentUser(r.Context())
	var req struct {
		Phone             string      `json:"phone"`
		Role              domain.Role `json:"role"`
		CommissionPercent float64     `json:"commissionPercent"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "درخواست معتبر نیست")
		return
	}
	if req.Role == "" {
		req.Role = domain.RoleConsultant
	}
	invite, err := a.deps.Invites.Create(r.Context(), user.ID, businessID, req.Phone, req.Role, req.CommissionPercent)
	if err != nil {
		writeError(w, http.StatusForbidden, "invite_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, invite)
}

func (a *api) listInvitations(w http.ResponseWriter, r *http.Request, businessID string) {
	items, err := a.deps.Invites.ListForBusiness(r.Context(), businessID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invite_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) invitationInbox(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	items, err := a.deps.Invites.Inbox(r.Context(), user.Phone)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invite_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) acceptInvitation(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	invite, err := a.deps.Invites.Accept(r.Context(), user, r.PathValue("invitationId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invite_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, invite)
}

func (a *api) rejectInvitation(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	invite, err := a.deps.Invites.Reject(r.Context(), user, r.PathValue("invitationId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invite_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, invite)
}

func (a *api) uploadBusinessLogo(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	if err := r.ParseMultipartForm(3 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "upload_failed", "فایل معتبر نیست")
		return
	}
	_, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "upload_failed", "فایل ارسال نشده است")
		return
	}
	file, err := a.deps.Files.SaveBusinessLogo(r.Context(), user.ID, header)
	if err != nil {
		writeError(w, http.StatusBadRequest, "upload_failed", err.Error())
		return
	}
	w.Header().Set("X-File-Size", strconv.FormatInt(file.Size, 10))
	writeJSON(w, http.StatusCreated, file)
}

func (a *api) uploadFile(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	if a.deps.Uploads == nil {
		writeError(w, http.StatusServiceUnavailable, "upload_unavailable", "upload service is unavailable")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 220<<20)
	if err := r.ParseMultipartForm(220 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "upload_failed", "فایل معتبر نیست یا حجم آن بیش از حد مجاز است")
		return
	}
	_, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "upload_failed", "فایل ارسال نشده است")
		return
	}
	result, err := a.deps.Uploads.Upload(r.Context(), user, service.UploadInput{
		Purpose:    r.FormValue("purpose"),
		TargetType: r.FormValue("targetType"),
		TargetID:   firstNonEmptyForm(r, "targetId", "channelId", "propertyId", "businessId"),
		BusinessID: r.FormValue("businessId"),
	}, header)
	if err != nil {
		writeError(w, http.StatusForbidden, "upload_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func (a *api) listLocations(w http.ResponseWriter, r *http.Request, businessID string) {
	user, _ := currentUser(r.Context())
	items, err := a.deps.Locations.List(r.Context(), user.ID, businessID)
	if err != nil {
		writeError(w, http.StatusForbidden, "locations_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) createArea(w http.ResponseWriter, r *http.Request, businessID string) {
	user, _ := currentUser(r.Context())
	var req struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "درخواست معتبر نیست")
		return
	}
	area, err := a.deps.Locations.CreateArea(r.Context(), user.ID, businessID, req.Name)
	if err != nil {
		writeError(w, http.StatusForbidden, "area_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, area)
}

func (a *api) deleteArea(w http.ResponseWriter, r *http.Request, businessID, areaID string) {
	user, _ := currentUser(r.Context())
	if err := a.deps.Locations.DeleteArea(r.Context(), user.ID, businessID, areaID); err != nil {
		writeError(w, http.StatusBadRequest, "area_delete_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *api) createStreet(w http.ResponseWriter, r *http.Request, businessID, areaID string) {
	user, _ := currentUser(r.Context())
	var req struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "درخواست معتبر نیست")
		return
	}
	street, err := a.deps.Locations.CreateStreet(r.Context(), user.ID, businessID, areaID, req.Name)
	if err != nil {
		writeError(w, http.StatusForbidden, "street_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, street)
}

func (a *api) deleteStreet(w http.ResponseWriter, r *http.Request, businessID, areaID, streetID string) {
	user, _ := currentUser(r.Context())
	if err := a.deps.Locations.DeleteStreet(r.Context(), user.ID, businessID, areaID, streetID); err != nil {
		writeError(w, http.StatusBadRequest, "street_delete_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *api) createNeighborhood(w http.ResponseWriter, r *http.Request, businessID, areaID, streetID string) {
	user, _ := currentUser(r.Context())
	var req struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "درخواست معتبر نیست")
		return
	}
	neighborhood, err := a.deps.Locations.CreateNeighborhood(r.Context(), user.ID, businessID, areaID, streetID, req.Name)
	if err != nil {
		writeError(w, http.StatusForbidden, "neighborhood_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, neighborhood)
}

func (a *api) deleteNeighborhood(w http.ResponseWriter, r *http.Request, businessID, areaID, streetID, neighborhoodID string) {
	user, _ := currentUser(r.Context())
	if err := a.deps.Locations.DeleteNeighborhood(r.Context(), user.ID, businessID, areaID, streetID, neighborhoodID); err != nil {
		writeError(w, http.StatusBadRequest, "neighborhood_delete_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *api) listProperties(w http.ResponseWriter, r *http.Request, businessID string) {
	user, _ := currentUser(r.Context())
	items, err := a.deps.Properties.List(r.Context(), user.ID, businessID)
	if err != nil {
		writeError(w, http.StatusForbidden, "properties_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) latestProperties(w http.ResponseWriter, r *http.Request, businessID string) {
	user, _ := currentUser(r.Context())
	limit, offset := channelPageParams(r)
	items, total, err := a.deps.Properties.Latest(
		r.Context(),
		user.ID,
		businessID,
		domain.PropertyFileType(r.URL.Query().Get("type")),
		limit,
		offset,
	)
	if err != nil {
		writeError(w, http.StatusForbidden, "latest_properties_failed", err.Error())
		return
	}
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	w.Header().Set("X-Limit", strconv.Itoa(limit))
	w.Header().Set("X-Offset", strconv.Itoa(offset))
	writeJSON(w, http.StatusOK, items)
}

func (a *api) createProperty(w http.ResponseWriter, r *http.Request, businessID string) {
	user, _ := currentUser(r.Context())
	var req domain.PropertyFile
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "درخواست معتبر نیست")
		return
	}
	file, err := a.deps.Properties.Create(r.Context(), user.ID, businessID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "property_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, file)
}

func (a *api) updateProperty(w http.ResponseWriter, r *http.Request, businessID, propertyID string) {
	user, _ := currentUser(r.Context())
	var req domain.PropertyFile
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "درخواست معتبر نیست")
		return
	}
	file, err := a.deps.Properties.Update(r.Context(), user.ID, businessID, propertyID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "property_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, file)
}

func (a *api) requestPropertyShare(w http.ResponseWriter, r *http.Request, businessID, propertyID string) {
	user, _ := currentUser(r.Context())
	var req struct {
		CommissionPercent float64 `json:"commissionPercent"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "درخواست معتبر نیست")
		return
	}
	request, err := a.deps.Properties.RequestShare(r.Context(), user, businessID, propertyID, req.CommissionPercent)
	if err != nil {
		writeError(w, http.StatusBadRequest, "share_request_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, request)
}

func (a *api) propertyShareRequests(w http.ResponseWriter, r *http.Request, businessID string) {
	user, _ := currentUser(r.Context())
	scope := r.URL.Query().Get("scope")
	items, err := a.deps.Properties.ShareRequests(r.Context(), user.ID, businessID, scope)
	if err != nil {
		writeError(w, http.StatusForbidden, "share_requests_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) propertyOffers(w http.ResponseWriter, r *http.Request, businessID string) {
	user, _ := currentUser(r.Context())
	items, err := a.deps.Properties.Offers(r.Context(), user.ID, businessID, r.URL.Query().Get("scope"))
	if err != nil {
		writeError(w, http.StatusForbidden, "property_offers_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) sendPropertyOffer(w http.ResponseWriter, r *http.Request, businessID, offerID string) {
	user, _ := currentUser(r.Context())
	var req struct {
		CommissionPercent float64 `json:"commissionPercent"`
	}
	_ = decodeJSON(r, &req)
	offer, err := a.deps.Properties.SendOffer(r.Context(), user.ID, businessID, offerID, req.CommissionPercent)
	if err != nil {
		writeError(w, http.StatusBadRequest, "property_offer_send_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, offer)
}

func (a *api) respondPropertyOffer(w http.ResponseWriter, r *http.Request, businessID, offerID string, approve bool) {
	user, _ := currentUser(r.Context())
	offer, err := a.deps.Properties.RespondOffer(r.Context(), user.ID, businessID, offerID, approve)
	if err != nil {
		writeError(w, http.StatusBadRequest, "property_offer_response_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, offer)
}

func (a *api) finalizePropertyOffer(w http.ResponseWriter, r *http.Request, businessID, offerID string, approve bool) {
	user, _ := currentUser(r.Context())
	offer, err := a.deps.Properties.FinalizeOffer(r.Context(), user.ID, businessID, offerID, approve)
	if err != nil {
		writeError(w, http.StatusBadRequest, "property_offer_finalize_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, offer)
}

func (a *api) decidePropertyShare(w http.ResponseWriter, r *http.Request, businessID, requestID string, approve bool) {
	user, _ := currentUser(r.Context())
	request, err := a.deps.Properties.DecideShareRequest(r.Context(), user.ID, businessID, requestID, approve)
	if err != nil {
		writeError(w, http.StatusForbidden, "share_request_decision_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, request)
}

func (a *api) receivePropertyShare(w http.ResponseWriter, r *http.Request, businessID, requestID string) {
	user, _ := currentUser(r.Context())
	file, err := a.deps.Properties.ReceiveSharedFile(r.Context(), user.ID, businessID, requestID)
	if err != nil {
		writeError(w, http.StatusForbidden, "share_receive_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, file)
}

func (a *api) uploadPropertyMedia(w http.ResponseWriter, r *http.Request, businessID, propertyID string) {
	user, _ := currentUser(r.Context())
	r.Body = http.MaxBytesReader(w, r.Body, 220<<20)
	if err := r.ParseMultipartForm(220 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "upload_failed", "فایل معتبر نیست یا حجم آن بیش از حد مجاز است")
		return
	}
	_, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "upload_failed", "فایل ارسال نشده است")
		return
	}
	file, err := a.deps.Properties.AddMedia(r.Context(), user.ID, businessID, propertyID, header)
	if err != nil {
		writeError(w, http.StatusBadRequest, "upload_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, file)
}

func (a *api) confirmPropertyMedia(w http.ResponseWriter, r *http.Request, businessID, propertyID string) {
	user, _ := currentUser(r.Context())
	var input struct {
		FileID string `json:"fileId"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request")
		return
	}
	file, err := a.deps.Properties.AddMediaFromUpload(r.Context(), user.ID, businessID, propertyID, input.FileID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "media_confirm_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, file)
}

func (a *api) contactTags(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.deps.Contacts.SystemTags())
}

func (a *api) listContacts(w http.ResponseWriter, r *http.Request, businessID string) {
	user, _ := currentUser(r.Context())
	items, err := a.deps.Contacts.List(r.Context(), user.ID, businessID, service.ContactFilter{
		Query: r.URL.Query().Get("q"),
		Tag:   r.URL.Query().Get("tag"),
		Phone: r.URL.Query().Get("phone"),
	})
	if err != nil {
		writeError(w, http.StatusForbidden, "contacts_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) createContact(w http.ResponseWriter, r *http.Request, businessID string) {
	user, _ := currentUser(r.Context())
	var req domain.Contact
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request")
		return
	}
	contact, err := a.deps.Contacts.Create(r.Context(), user.ID, businessID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "contact_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, contact)
}

func (a *api) addProfileContactCategories(w http.ResponseWriter, r *http.Request, businessID string) {
	user, _ := currentUser(r.Context())
	var req struct {
		Phone       string   `json:"phone"`
		DisplayName string   `json:"displayName"`
		Tags        []string `json:"tags"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request")
		return
	}
	result, err := a.deps.Contacts.AddProfileToCategories(r.Context(), user.ID, businessID, req.Phone, req.DisplayName, req.Tags)
	if err != nil {
		writeError(w, http.StatusBadRequest, "contact_categories_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (a *api) updateContact(w http.ResponseWriter, r *http.Request, businessID, contactID string) {
	user, _ := currentUser(r.Context())
	var req domain.Contact
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request")
		return
	}
	contact, err := a.deps.Contacts.Update(r.Context(), user.ID, businessID, contactID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "contact_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, contact)
}

func (a *api) requestMatches(w http.ResponseWriter, r *http.Request, businessID, contactID, requestID string) {
	user, _ := currentUser(r.Context())
	limit, offset := matchPageParams(r)
	items, total, err := a.deps.Matching.RequestMatches(r.Context(), user.ID, businessID, contactID, requestID, limit, offset)
	if err != nil {
		writeError(w, http.StatusBadRequest, "matches_failed", err.Error())
		return
	}
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	w.Header().Set("X-Limit", strconv.Itoa(limit))
	w.Header().Set("X-Offset", strconv.Itoa(offset))
	writeJSON(w, http.StatusOK, items)
}

func (a *api) listChannels(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	items, err := a.deps.Channels.List(r.Context(), user)
	if err != nil {
		writeError(w, http.StatusBadRequest, "channels_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) createPrivateChannel(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	var input struct {
		Phone string `json:"phone"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request")
		return
	}
	channel, err := a.deps.Channels.PrivateChat(r.Context(), user, input.Phone)
	if err != nil {
		writeError(w, http.StatusBadRequest, "private_channel_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, channel)
}

func (a *api) userVaults(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	items, err := a.deps.Channels.UserVaults(r.Context(), user)
	if err != nil {
		writeError(w, http.StatusBadRequest, "vaults_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) createUserVault(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	var input struct {
		Title string `json:"title"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request")
		return
	}
	vault, err := a.deps.Channels.CreateUserVault(r.Context(), user, input.Title)
	if err != nil {
		writeError(w, http.StatusBadRequest, "vault_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, vault)
}

func (a *api) businessVaults(w http.ResponseWriter, r *http.Request, businessID string) {
	user, _ := currentUser(r.Context())
	items, err := a.deps.Channels.BusinessVaults(r.Context(), user.ID, businessID)
	if err != nil {
		writeError(w, http.StatusForbidden, "vaults_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) createBusinessVault(w http.ResponseWriter, r *http.Request, businessID string) {
	user, _ := currentUser(r.Context())
	var input struct {
		Title string `json:"title"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request")
		return
	}
	vault, err := a.deps.Channels.CreateBusinessVault(r.Context(), user.ID, businessID, input.Title)
	if err != nil {
		writeError(w, http.StatusForbidden, "vault_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, vault)
}

func (a *api) userMainChannel(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	channel, err := a.deps.Channels.UserMain(r.Context(), user)
	if err != nil {
		writeError(w, http.StatusBadRequest, "channel_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, channel)
}

func (a *api) businessMainChannel(w http.ResponseWriter, r *http.Request, businessID string) {
	user, _ := currentUser(r.Context())
	channel, err := a.deps.Channels.BusinessMain(r.Context(), user.ID, businessID)
	if err != nil {
		writeError(w, http.StatusForbidden, "channel_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, channel)
}

func (a *api) channelMessages(w http.ResponseWriter, r *http.Request, channelID string) {
	user, _ := currentUser(r.Context())
	limit, offset := channelPageParams(r)
	fromUnread := r.URL.Query().Get("fromUnread") == "true"
	window := fromUnread || r.URL.Query().Get("window") == "true"
	page, err := a.deps.Channels.Messages(r.Context(), user.ID, channelID, limit, offset, fromUnread)
	if err != nil {
		writeError(w, http.StatusForbidden, "channel_messages_failed", err.Error())
		return
	}
	w.Header().Set("X-Total-Count", strconv.Itoa(page.Total))
	w.Header().Set("X-Limit", strconv.Itoa(page.Limit))
	w.Header().Set("X-Offset", strconv.Itoa(page.Offset))
	if window {
		writeJSON(w, http.StatusOK, page)
		return
	}
	writeJSON(w, http.StatusOK, page.Items)
}

func (a *api) createChannelMessage(w http.ResponseWriter, r *http.Request, channelID string) {
	user, _ := currentUser(r.Context())
	var input domain.ChannelMessage
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request")
		return
	}
	message, err := a.deps.Channels.SendMessage(r.Context(), user, channelID, input)
	if err != nil {
		writeError(w, http.StatusBadRequest, "channel_message_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, message)
}

func (a *api) updateChannelMessage(w http.ResponseWriter, r *http.Request, channelID, messageID string) {
	user, _ := currentUser(r.Context())
	var input domain.ChannelMessage
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request")
		return
	}
	message, err := a.deps.Channels.UpdateMessage(r.Context(), user, channelID, messageID, input)
	if err != nil {
		writeError(w, http.StatusForbidden, "channel_message_update_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, message)
}

func (a *api) deleteChannelMessage(w http.ResponseWriter, r *http.Request, channelID, messageID string) {
	user, _ := currentUser(r.Context())
	if err := a.deps.Channels.DeleteMessage(r.Context(), user, channelID, messageID); err != nil {
		writeError(w, http.StatusForbidden, "channel_message_delete_failed", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *api) uploadChannelMedia(w http.ResponseWriter, r *http.Request, channelID string) {
	user, _ := currentUser(r.Context())
	if err := r.ParseMultipartForm(220 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "upload_failed", "فایل معتبر نیست یا حجم آن بیش از حد مجاز است")
		return
	}
	_, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "upload_failed", "فایل ارسال نشده است")
		return
	}
	media, err := a.deps.Channels.SaveMedia(r.Context(), user, channelID, header)
	if err != nil {
		writeError(w, http.StatusBadRequest, "upload_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, media)
}

func (a *api) channelVaultFiles(w http.ResponseWriter, r *http.Request, channelID string) {
	user, _ := currentUser(r.Context())
	limit, offset := channelPageParams(r)
	items, total, err := a.deps.Channels.VaultFiles(r.Context(), user.ID, channelID, limit, offset)
	if err != nil {
		writeError(w, http.StatusForbidden, "channel_vault_failed", err.Error())
		return
	}
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	w.Header().Set("X-Limit", strconv.Itoa(limit))
	w.Header().Set("X-Offset", strconv.Itoa(offset))
	writeJSON(w, http.StatusOK, items)
}

func (a *api) channelVaultFile(w http.ResponseWriter, r *http.Request, channelID, fileID string) {
	user, _ := currentUser(r.Context())
	file, err := a.deps.Channels.VaultFile(r.Context(), user.ID, channelID, fileID)
	if err != nil {
		writeError(w, http.StatusForbidden, "channel_vault_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, file)
}

func (a *api) uploadChannelVaultFile(w http.ResponseWriter, r *http.Request, channelID string) {
	user, _ := currentUser(r.Context())
	if err := r.ParseMultipartForm(60 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "upload_failed", "فایل معتبر نیست یا حجم آن بیش از حد مجاز است")
		return
	}
	_, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "upload_failed", "فایل ارسال نشده است")
		return
	}
	file, err := a.deps.Channels.SaveVaultFile(
		r.Context(),
		user,
		channelID,
		header,
		r.FormValue("title"),
		r.FormValue("note"),
	)
	if err != nil {
		writeError(w, http.StatusBadRequest, "upload_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, file)
}

func (a *api) confirmChannelVaultFile(w http.ResponseWriter, r *http.Request, channelID string) {
	user, _ := currentUser(r.Context())
	var input struct {
		FileID string `json:"fileId"`
		Title  string `json:"title"`
		Note   string `json:"note"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request")
		return
	}
	file, err := a.deps.Channels.SaveVaultFileFromUpload(r.Context(), user, channelID, input.FileID, input.Title, input.Note)
	if err != nil {
		writeError(w, http.StatusBadRequest, "vault_file_confirm_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, file)
}

func (a *api) channelMembers(w http.ResponseWriter, r *http.Request, channelID string) {
	user, _ := currentUser(r.Context())
	items, err := a.deps.Channels.Members(r.Context(), user.ID, channelID)
	if err != nil {
		writeError(w, http.StatusForbidden, "channel_members_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) updateChannelMemberRole(w http.ResponseWriter, r *http.Request, channelID, memberID string) {
	user, _ := currentUser(r.Context())
	var input struct {
		Admin bool `json:"admin"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request")
		return
	}
	member, err := a.deps.Channels.SetMemberAdmin(r.Context(), user.ID, channelID, memberID, input.Admin)
	if err != nil {
		writeError(w, http.StatusForbidden, "channel_member_role_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, member)
}

func (a *api) inviteChannelMember(w http.ResponseWriter, r *http.Request, channelID string) {
	user, _ := currentUser(r.Context())
	var input struct {
		Phone string `json:"phone"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request")
		return
	}
	member, err := a.deps.Channels.InviteByPhone(r.Context(), user, channelID, input.Phone)
	if err != nil {
		writeError(w, http.StatusBadRequest, "channel_invite_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, member)
}

func queryInt(r *http.Request, key string, fallback int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func firstNonEmptyForm(r *http.Request, keys ...string) string {
	for _, key := range keys {
		if value := r.FormValue(key); value != "" {
			return value
		}
	}
	return ""
}

func channelPageParams(r *http.Request) (int, int) {
	limit := queryInt(r, "limit", 30)
	offset := queryInt(r, "offset", 0)
	if limit <= 0 {
		limit = 30
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func matchPageParams(r *http.Request) (int, int) {
	limit := queryInt(r, "limit", 50)
	offset := queryInt(r, "offset", 0)
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func (a *api) catalogCities(w http.ResponseWriter, r *http.Request) {
	items, err := a.deps.Platform.ListCities(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "catalog_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) catalogCityLocations(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	if user.CityID == "" || user.CityID != r.PathValue("cityId") {
		writeError(w, http.StatusForbidden, "city_forbidden", "دسترسی به موقعیت های این شهر مجاز نیست")
		return
	}
	data, err := a.deps.Platform.CityLocations(r.Context(), r.PathValue("cityId"))
	if err != nil {
		writeError(w, http.StatusNotFound, "city_not_found", "city not found")
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (a *api) catalogSearchLocations(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	if user.CityID == "" || user.CityID != r.PathValue("cityId") {
		writeError(w, http.StatusForbidden, "city_forbidden", "دسترسی به موقعیت های این شهر مجاز نیست")
		return
	}
	data, err := a.deps.Platform.SearchCityLocations(r.Context(), r.PathValue("cityId"), r.URL.Query().Get("q"))
	if err != nil {
		writeError(w, http.StatusNotFound, "city_not_found", "city not found")
		return
	}
	writeJSON(w, http.StatusOK, data)
}

func (a *api) createLocationSuggestion(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	var req domain.LocationSuggestion
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request")
		return
	}
	item, err := a.deps.Platform.CreateLocationSuggestion(r.Context(), user, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "suggestion_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (a *api) adminMe(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	account, err := a.deps.Platform.Me(r.Context(), user)
	if err != nil {
		writeError(w, http.StatusForbidden, "forbidden", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, account)
}

func (a *api) adminUsers(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	items, err := a.deps.Platform.ListUsers(r.Context(), user)
	if err != nil {
		writeError(w, http.StatusForbidden, "forbidden", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) adminBusinesses(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	items, err := a.deps.Platform.ListBusinesses(r.Context(), user)
	if err != nil {
		writeError(w, http.StatusForbidden, "forbidden", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) adminAccounts(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	items, err := a.deps.Platform.ListAdminAccounts(r.Context(), user)
	if err != nil {
		writeError(w, http.StatusForbidden, "forbidden", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) adminSaveAccount(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	var req struct {
		UserID string             `json:"userId"`
		Roles  []domain.AdminRole `json:"roles"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request")
		return
	}
	account, err := a.deps.Platform.SaveAdminAccount(r.Context(), user, req.UserID, req.Roles)
	if err != nil {
		writeError(w, http.StatusForbidden, "forbidden", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, account)
}

func (a *api) adminSettings(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	settings, err := a.deps.Platform.Settings(r.Context(), user)
	if err != nil {
		writeError(w, http.StatusForbidden, "forbidden", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

func (a *api) adminUpdateSettings(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	var req struct {
		OTPAPIKey        string `json:"otpApiKey"`
		ServiceSMSAPIKey string `json:"serviceSmsApiKey"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request")
		return
	}
	settings, err := a.deps.Platform.UpdateSettings(r.Context(), user, req.OTPAPIKey, req.ServiceSMSAPIKey)
	if err != nil {
		writeError(w, http.StatusForbidden, "forbidden", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

func (a *api) adminCreateCity(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	var req struct {
		Name string `json:"name"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request")
		return
	}
	city, err := a.deps.Platform.CreateCity(r.Context(), user, req.Name)
	if err != nil {
		writeError(w, http.StatusForbidden, "forbidden", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, city)
}

func (a *api) adminLocationSuggestions(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	status := domain.LocationSuggestionStatus(r.URL.Query().Get("status"))
	items, err := a.deps.Platform.ListLocationSuggestions(r.Context(), user, status)
	if err != nil {
		writeError(w, http.StatusForbidden, "forbidden", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *api) adminApproveLocationSuggestion(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	item, err := a.deps.Platform.ApproveLocationSuggestion(r.Context(), user, r.PathValue("suggestionId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "approve_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *api) adminRejectLocationSuggestion(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	var req struct {
		Note string `json:"note"`
	}
	_ = decodeJSON(r, &req)
	item, err := a.deps.Platform.RejectLocationSuggestion(r.Context(), user, r.PathValue("suggestionId"), req.Note)
	if err != nil {
		writeError(w, http.StatusBadRequest, "reject_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, item)
}
