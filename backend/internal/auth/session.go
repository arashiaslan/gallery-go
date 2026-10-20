package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/squarebrain/gallery/internal/store"
	"github.com/squarebrain/gallery/internal/utils"
)

const (
	SessionCookieName = "session_token"
	SessionDuration   = 30 * 24 * time.Hour // 30 days
)

// CreateSession creates a new session and returns the token
func CreateSession(ctx context.Context, db *store.DB, userID string) (*store.Session, error) {
	session := &store.Session{
		ID:        utils.GenerateShortID(),
		UserID:    userID,
		Token:     utils.GenerateToken(),
		ExpiresAt: time.Now().Add(SessionDuration),
	}

	if err := db.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// SetSessionCookie sets the session cookie in the response
func SetSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: 2, // http.SameSiteLax
		Expires:  time.Now().Add(SessionDuration),
		MaxAge:   int(SessionDuration.Seconds()),
	})
}

// ClearSessionCookie removes the session cookie
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Path:     "/",
		MaxAge:   -1,
	})
}

// GetSessionFromRequest extracts the session token from the request
func GetSessionFromRequest(r *http.Request) string {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}
