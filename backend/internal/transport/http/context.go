package httptransport

import (
	"context"

	"amlakcrm/backend/internal/domain"
)

type contextKey string

const userKey contextKey = "user"
const sessionKey contextKey = "session"

func withUser(ctx context.Context, user domain.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func withSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, sessionKey, sessionID)
}

func currentUser(ctx context.Context) (domain.User, bool) {
	user, ok := ctx.Value(userKey).(domain.User)
	return user, ok
}

func currentSessionID(ctx context.Context) string {
	sessionID, _ := ctx.Value(sessionKey).(string)
	return sessionID
}
