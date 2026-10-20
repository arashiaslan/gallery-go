package handlers

import (
	"context"
	"net/http"

	"github.com/squarebrain/gallery/internal/auth"
	"github.com/squarebrain/gallery/internal/store"
)

// AuthMiddleware ensures user is authenticated
func AuthMiddleware(db *store.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := auth.GetSessionFromRequest(r)
			if token == "" {
				WriteError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			session, err := db.GetSessionByToken(r.Context(), token)
			if err != nil || session == nil {
				WriteError(w, http.StatusUnauthorized, "Invalid or expired session")
				return
			}

			user, err := db.GetUserByID(r.Context(), session.UserID)
			if err != nil || user == nil {
				WriteError(w, http.StatusUnauthorized, "User not found")
				return
			}

			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AdminMiddleware ensures user is admin
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*store.User)
		if user.Role != "admin" {
			WriteError(w, http.StatusForbidden, "Admin role required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
