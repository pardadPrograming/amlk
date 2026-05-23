package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type AccessClaims struct {
	UserID    string `json:"userId"`
	SessionID string `json:"sessionId"`
	jwt.RegisteredClaims
}

func NewTokenService(secret string, accessTTL, refreshTTL time.Duration) *TokenService {
	return &TokenService{secret: []byte(secret), accessTTL: accessTTL, refreshTTL: refreshTTL}
}

func (s *TokenService) IssueAccess(userID string, sessionID string) (string, time.Time, error) {
	expires := time.Now().UTC().Add(s.accessTTL)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessClaims{
		UserID:    userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expires),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	})
	value, err := token.SignedString(s.secret)
	return value, expires, err
}

func (s *TokenService) ParseAccess(value string) (string, string, error) {
	claims := &AccessClaims{}
	token, err := jwt.ParseWithClaims(value, claims, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil || !token.Valid || claims.UserID == "" {
		return "", "", errors.New("invalid token")
	}
	return claims.UserID, claims.SessionID, nil
}

func (s *TokenService) NewRefresh() (string, time.Time) {
	var b [32]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:]), time.Now().UTC().Add(s.refreshTTL)
}
