package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/squarebrain/gallery/internal/auth"
	"github.com/squarebrain/gallery/internal/store"
	"github.com/squarebrain/gallery/internal/utils"
)

type AuthHandler struct {
	db      *store.DB
	limiter *authRateLimiter
}

func NewAuthHandler(db *store.DB) *AuthHandler {
	return &AuthHandler{db: db, limiter: newAuthRateLimiter()}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	LimitBody(w, r)
	if !h.limiter.Allow(clientIP(r)) {
		WriteError(w, http.StatusTooManyRequests, "Too many attempts, please try again later")
		return
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.Email == "" || req.Password == "" {
		WriteError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	user, err := h.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}

	if user == nil || !auth.VerifyPassword(user.PasswordHash, req.Password) {
		WriteError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Create session
	session, err := auth.CreateSession(r.Context(), h.db, user.ID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	// Set session cookie
	auth.SetSessionCookie(w, session.Token)

	WriteSuccess(w, http.StatusOK, LoginResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	})
}

// Logout handles POST /api/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := auth.GetSessionFromRequest(r)
	if token != "" {
		h.db.DeleteSession(r.Context(), token)
	}

	auth.ClearSessionCookie(w)
	WriteSuccess(w, http.StatusOK, map[string]string{"message": "Logged out"})
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register handles POST /api/auth/register?token=...
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		WriteError(w, http.StatusBadRequest, "Invite token is required")
		return
	}

	var req RegisterRequest
	LimitBody(w, r)
	if !h.limiter.Allow(clientIP(r)) {
		WriteError(w, http.StatusTooManyRequests, "Too many attempts, please try again later")
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.Email == "" || req.Password == "" {
		WriteError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Verify invite
	invite, err := h.db.GetInviteByToken(r.Context(), token)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}

	if invite == nil {
		WriteError(w, http.StatusNotFound, "Invalid invite token")
		return
	}

	if invite.UsedByID != nil {
		WriteError(w, http.StatusBadRequest, "Invite already used")
		return
	}

	if time.Now().After(invite.ExpiresAt) {
		WriteError(w, http.StatusBadRequest, "Invite has expired")
		return
	}

	// Check if user already exists
	existing, err := h.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}

	if existing != nil {
		WriteError(w, http.StatusBadRequest, "Email already registered")
		return
	}

	// Hash password
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	// Create user
	user := &store.User{
		ID:           utils.GenerateShortID(),
		Email:        req.Email,
		PasswordHash: hash,
		Role:         "member",
	}

	if err := h.db.CreateUser(r.Context(), user); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Mark invite as used
	if err := h.db.MarkInviteAsUsed(r.Context(), token, user.ID); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to process invite")
		return
	}

	// Create session
	session, err := auth.CreateSession(r.Context(), h.db, user.ID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	auth.SetSessionCookie(w, session.Token)

	WriteSuccess(w, http.StatusCreated, LoginResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	})
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// GetMe handles GET /api/me
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*store.User)
	WriteSuccess(w, http.StatusOK, UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	})
}

type CreateInviteRequest struct {
	Email string `json:"email"`
}

type InviteResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

// CreateInvite handles POST /api/invites (admin only)
func (h *AuthHandler) CreateInvite(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*store.User)
	if user.Role != "admin" {
		WriteError(w, http.StatusForbidden, "Admin role required")
		return
	}

	var req CreateInviteRequest
	LimitBody(w, r)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.Email == "" {
		WriteError(w, http.StatusBadRequest, "Email is required")
		return
	}

	invite := &store.Invite{
		ID:          utils.GenerateShortID(),
		Email:       req.Email,
		Token:       utils.GenerateToken(),
		CreatedByID: user.ID,
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour), // Valid for 7 days
	}

	if err := h.db.CreateInvite(r.Context(), invite); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to create invite")
		return
	}

	WriteSuccess(w, http.StatusCreated, InviteResponse{
		ID:        invite.ID,
		Email:     invite.Email,
		Token:     invite.Token,
		ExpiresAt: invite.ExpiresAt.Format(time.RFC3339),
	})
}
