package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/squarebrain/gallery/internal/auth"
	"github.com/squarebrain/gallery/internal/store"
	"github.com/squarebrain/gallery/internal/utils"
)

func setupTestDB(t *testing.T) *store.DB {
	db, err := store.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	if err := db.Migrate(); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func TestLoginSuccess(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create test user
	hash, _ := auth.HashPassword("password123")
	user := &store.User{
		ID:           utils.GenerateShortID(),
		Email:        "test@example.com",
		PasswordHash: hash,
		Role:         "member",
	}
	if err := db.CreateUser(context.Background(), user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	handler := NewAuthHandler(db)

	body := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp SuccessResponse
	json.NewDecoder(w.Body).Decode(&resp)
	loginResp := resp.Data.(map[string]interface{})

	if loginResp["email"] != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %v", loginResp["email"])
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create test user
	hash, _ := auth.HashPassword("password123")
	user := &store.User{
		ID:           utils.GenerateShortID(),
		Email:        "test@example.com",
		PasswordHash: hash,
		Role:         "member",
	}
	if err := db.CreateUser(context.Background(), user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	handler := NewAuthHandler(db)

	body := LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	handler.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestRegisterWithInviteToken(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create admin user
	hash, _ := auth.HashPassword("admin123")
	admin := &store.User{
		ID:           utils.GenerateShortID(),
		Email:        "admin@example.com",
		PasswordHash: hash,
		Role:         "admin",
	}
	if err := db.CreateUser(context.Background(), admin); err != nil {
		t.Fatalf("Failed to create admin: %v", err)
	}

	// Create invite
	invite := &store.Invite{
		ID:          utils.GenerateShortID(),
		Email:       "newuser@example.com",
		Token:       utils.GenerateToken(),
		CreatedByID: admin.ID,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}
	if err := db.CreateInvite(context.Background(), invite); err != nil {
		t.Fatalf("Failed to create invite: %v", err)
	}

	handler := NewAuthHandler(db)

	body := RegisterRequest{
		Email:    "newuser@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/auth/register?token="+invite.Token, bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()

	handler.Register(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	// Verify user was created
	user, _ := db.GetUserByEmail(context.Background(), "newuser@example.com")
	if user == nil {
		t.Error("User was not created")
	}
}

func TestAccessWithoutSession(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewAuthHandler(db)

	req := httptest.NewRequest("GET", "/api/me", nil)
	w := httptest.NewRecorder()

	AuthMiddleware(db)(
		http.HandlerFunc(handler.GetMe),
	).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}
