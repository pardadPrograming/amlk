package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"amlakcrm/backend/internal/config"
	"amlakcrm/backend/internal/domain"
	"amlakcrm/backend/internal/repository"
	"amlakcrm/backend/internal/support"
)

var phonePattern = regexp.MustCompile(`^(\+98|0)?9\d{9}$`)

type AuthService struct {
	store  repository.Store
	tokens *TokenService
	cfg    config.Config
}

type AuthTokens struct {
	AccessToken      string      `json:"accessToken"`
	RefreshToken     string      `json:"refreshToken"`
	AccessExpiresAt  time.Time   `json:"accessExpiresAt"`
	RefreshExpiresAt time.Time   `json:"refreshExpiresAt"`
	User             domain.User `json:"user"`
}

type SessionView struct {
	domain.Session
	Current bool `json:"current"`
}

func NewAuthService(store repository.Store, tokens *TokenService, cfg config.Config) *AuthService {
	return &AuthService{store: store, tokens: tokens, cfg: cfg}
}

func NormalizePhone(phone string) (string, error) {
	if !phonePattern.MatchString(phone) {
		return "", errors.New("شماره موبایل معتبر نیست")
	}
	if len(phone) == 10 && phone[0] == '9' {
		return "+98" + phone, nil
	}
	if len(phone) == 11 && phone[0] == '0' {
		return "+98" + phone[1:], nil
	}
	return phone, nil
}

func (s *AuthService) RequestOTP(ctx context.Context, phone string) (domain.OTPChallenge, error) {
	normalized, err := NormalizePhone(phone)
	if err != nil {
		return domain.OTPChallenge{}, err
	}
	existing, err := s.store.GetOTP(ctx, normalized)
	if err == nil && time.Since(existing.LastSentAt) < s.cfg.OTPResendAfter {
		return domain.OTPChallenge{}, errors.New("برای ارسال مجدد کد کمی صبر کنید")
	}
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	challenge := domain.OTPChallenge{
		Phone:       normalized,
		Code:        code,
		ExpiresAt:   time.Now().UTC().Add(s.cfg.OTPTTL),
		LastSentAt:  time.Now().UTC(),
		SendCount:   existing.SendCount + 1,
		Development: code,
	}
	return challenge, s.store.SaveOTP(ctx, challenge)
}

func (s *AuthService) LatestOTP(ctx context.Context) (domain.OTPChallenge, error) {
	return s.store.GetLatestOTP(ctx)
}

func (s *AuthService) VerifyOTP(ctx context.Context, phone, code, userAgent, ip string) (AuthTokens, error) {
	normalized, err := NormalizePhone(phone)
	if err != nil {
		return AuthTokens{}, err
	}
	challenge, err := s.store.GetOTP(ctx, normalized)
	if err != nil || time.Now().UTC().After(challenge.ExpiresAt) {
		return AuthTokens{}, errors.New("کد تایید منقضی شده است")
	}
	if challenge.Attempts >= s.cfg.OTPMaxAttempts {
		return AuthTokens{}, errors.New("تعداد تلاش‌های ناموفق بیش از حد مجاز است")
	}
	if challenge.Code != code {
		challenge.Attempts++
		_ = s.store.SaveOTP(ctx, challenge)
		return AuthTokens{}, errors.New("کد تایید صحیح نیست")
	}
	_ = s.store.DeleteOTP(ctx, normalized)

	user, err := s.store.UpsertUserByPhone(ctx, normalized)
	if err != nil {
		return AuthTokens{}, err
	}
	return s.issueTokens(ctx, user, userAgent, ip)
}

func (s *AuthService) Refresh(ctx context.Context, refresh, userAgent, ip string) (AuthTokens, error) {
	session, err := s.store.GetSessionByRefresh(ctx, refresh)
	if err != nil || time.Now().UTC().After(session.ExpiresAt) {
		return AuthTokens{}, errors.New("نشست معتبر نیست")
	}
	user, err := s.store.GetUser(ctx, session.UserID)
	if err != nil {
		return AuthTokens{}, err
	}
	_ = s.store.RevokeSession(ctx, refresh)
	return s.issueTokens(ctx, user, userAgent, ip)
}

func (s *AuthService) Logout(ctx context.Context, refresh string) error {
	return s.store.RevokeSession(ctx, refresh)
}

func (s *AuthService) SecurityProfile(ctx context.Context, userID, currentSessionID string) (map[string]interface{}, error) {
	user, err := s.store.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	sessions, err := s.ActiveSessions(ctx, userID, currentSessionID)
	if err != nil {
		return nil, err
	}
	var lastActivity time.Time
	for _, session := range sessions {
		if session.LastSeenAt.After(lastActivity) {
			lastActivity = session.LastSeenAt
		}
	}
	return map[string]interface{}{
		"user":            user,
		"privacySettings": user.PrivacySettings,
		"sessions":        sessions,
		"lastActivityAt":  lastActivity,
	}, nil
}

func (s *AuthService) UpdatePrivacy(ctx context.Context, userID string, settings domain.PrivacySettings) (domain.PrivacySettings, error) {
	user, err := s.store.GetUser(ctx, userID)
	if err != nil {
		return domain.PrivacySettings{}, err
	}
	user.PrivacySettings = defaultPrivacySettings()
	updated, err := s.store.UpdateUser(ctx, user)
	if err != nil {
		return domain.PrivacySettings{}, err
	}
	return updated.PrivacySettings, nil
}

func defaultPrivacySettings() domain.PrivacySettings {
	return domain.PrivacySettings{
		ShowPhoneToTeam:    true,
		ShowActivityStatus: true,
		AllowInviteByPhone: true,
	}
}

func (s *AuthService) ActiveSessions(ctx context.Context, userID, currentSessionID string) ([]SessionView, error) {
	sessions, err := s.store.ListSessions(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]SessionView, 0, len(sessions))
	for _, session := range sessions {
		session.RefreshToken = ""
		result = append(result, SessionView{
			Session: session,
			Current: session.ID == currentSessionID,
		})
	}
	return result, nil
}

func (s *AuthService) RevokeUserSession(ctx context.Context, userID, sessionID string) error {
	return s.store.RevokeSessionByID(ctx, userID, sessionID)
}

func (s *AuthService) TouchSession(ctx context.Context, sessionID string) {
	if sessionID == "" {
		return
	}
	_ = s.store.TouchSession(ctx, sessionID, time.Now().UTC())
}

func (s *AuthService) CompleteProfile(ctx context.Context, userID, firstName, lastName, cityID string) (domain.User, error) {
	user, err := s.store.GetUser(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	cityID = strings.TrimSpace(cityID)
	if firstName == "" || lastName == "" || cityID == "" {
		return domain.User{}, errors.New("نام، نام خانوادگی و شهر الزامی است")
	}
	if _, err := s.store.GetCity(ctx, cityID); err != nil {
		return domain.User{}, repository.ErrNotFound
	}
	user.FirstName = firstName
	user.LastName = lastName
	user.DisplayName = strings.TrimSpace(firstName + " " + lastName)
	user.CityID = cityID
	return s.store.UpdateUser(ctx, user)
}

func (s *AuthService) issueTokens(ctx context.Context, user domain.User, userAgent, ip string) (AuthTokens, error) {
	sessionID := support.NewID()
	refresh, refreshExp := s.tokens.NewRefresh()
	access, accessExp, err := s.tokens.IssueAccess(user.ID, sessionID)
	if err != nil {
		return AuthTokens{}, err
	}
	now := time.Now().UTC()
	deviceName, deviceType, browser, osName := describeDevice(userAgent)
	session := domain.Session{
		ID:           sessionID,
		UserID:       user.ID,
		RefreshToken: refresh,
		UserAgent:    userAgent,
		DeviceName:   deviceName,
		DeviceType:   deviceType,
		Browser:      browser,
		OS:           osName,
		IP:           ip,
		LastSeenAt:   now,
		ExpiresAt:    refreshExp,
		CreatedAt:    now,
	}
	if err := s.store.SaveSession(ctx, session); err != nil {
		return AuthTokens{}, err
	}
	return AuthTokens{AccessToken: access, RefreshToken: refresh, AccessExpiresAt: accessExp, RefreshExpiresAt: refreshExp, User: user}, nil
}

func describeDevice(userAgent string) (string, string, string, string) {
	ua := strings.ToLower(userAgent)
	deviceType := "desktop"
	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") || strings.Contains(ua, "iphone") {
		deviceType = "mobile"
	}
	osName := "Unknown"
	switch {
	case strings.Contains(ua, "windows"):
		osName = "Windows"
	case strings.Contains(ua, "android"):
		osName = "Android"
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") || strings.Contains(ua, "ios"):
		osName = "iOS"
	case strings.Contains(ua, "mac os"):
		osName = "macOS"
	case strings.Contains(ua, "linux"):
		osName = "Linux"
	}
	browser := "Browser"
	switch {
	case strings.Contains(ua, "edg/"):
		browser = "Edge"
	case strings.Contains(ua, "chrome/"):
		browser = "Chrome"
	case strings.Contains(ua, "firefox/"):
		browser = "Firefox"
	case strings.Contains(ua, "safari/"):
		browser = "Safari"
	}
	return strings.TrimSpace(browser + " on " + osName), deviceType, browser, osName
}
