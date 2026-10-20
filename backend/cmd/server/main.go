package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/squarebrain/gallery/internal/auth"
	"github.com/squarebrain/gallery/internal/config"
	"github.com/squarebrain/gallery/internal/handlers"
	"github.com/squarebrain/gallery/internal/media"
	"github.com/squarebrain/gallery/internal/store"
	"github.com/squarebrain/gallery/internal/utils"
	"github.com/squarebrain/gallery/internal/web"
)



func main() {
	cfg := config.Load()

	// Create data directory
	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize database
	db, err := store.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Seed admin user if needed
	if err := seedAdmin(db, cfg); err != nil {
		log.Fatalf("Failed to seed admin: %v", err)
	}

	// Periodically purge expired sessions
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if err := db.DeleteExpiredSessions(context.Background()); err != nil {
				log.Printf("Failed to delete expired sessions: %v", err)
			}
		}
	}()

	// Create image processor with 4 workers
	processor := media.NewProcessor(cfg.DataDir, 4, db)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := processor.Shutdown(shutdownCtx); err != nil {
			log.Printf("Error shutting down processor: %v", err)
		}
	}()

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Auth routes (no authentication required)
	authHandler := handlers.NewAuthHandler(db)
	r.Post("/api/auth/login", authHandler.Login)
	r.Post("/api/auth/register", authHandler.Register)

	// Public read-only routes (guest view): browse albums, no mutations
	guestHandler := handlers.NewGuestHandler(db)
	r.Get("/api/guest/books", guestHandler.ListBooks)
	r.Get("/api/guest/books/{id}", guestHandler.GetBook)
	r.Get("/api/guest/books/{id}/pages", guestHandler.GetBookPages)

	// Media serving (public read; photo IDs are unguessable)
	mediaHandler := handlers.NewMediaHandler(db, cfg.DataDir)
	r.Get("/media/{book_id}/{photo_id}/{size}", mediaHandler.ServeMedia)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(handlers.AuthMiddleware(db))

		// Auth
		r.Post("/api/auth/logout", authHandler.Logout)
		r.Get("/api/me", authHandler.GetMe)

		// Admin routes
		r.Group(func(r chi.Router) {
			r.Use(handlers.AdminMiddleware)
			r.Post("/api/invites", authHandler.CreateInvite)
		})

		// Books
		bookHandler := handlers.NewBookHandler(db, cfg.DataDir)
		r.Get("/api/books", bookHandler.ListBooks)
		r.Post("/api/books", bookHandler.CreateBook)
		r.Get("/api/books/{id}", bookHandler.GetBook)
		r.Patch("/api/books/{id}", bookHandler.UpdateBook)
		r.Delete("/api/books/{id}", bookHandler.DeleteBook)

		// Photos
		photoHandler := handlers.NewPhotoHandler(db, cfg.DataDir, processor, cfg.MaxUploadMB)
		r.Post("/api/books/{id}/photos", photoHandler.UploadPhotos)
		r.Patch("/api/photos/{id}", photoHandler.UpdatePhoto)
		r.Delete("/api/photos/{id}", photoHandler.DeletePhoto)
		r.Patch("/api/books/{id}/reorder", photoHandler.ReorderPhotos)
		r.Get("/api/books/{id}/pages", photoHandler.GetBookPages)
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok"}`)
	})

	// Serve the embedded frontend (production/Docker builds). In development
	// the dist directory only holds a placeholder, so the API runs alone and
	// the frontend is served separately by the Vite dev server.
	webFS := web.FS()
	if _, err := fs.Stat(webFS, "index.html"); err == nil {
		r.NotFound(func(w http.ResponseWriter, req *http.Request) {
			path := req.URL.Path

			// API/media 404s stay JSON
			if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/media/") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, `{"error":"not found"}`)
				return
			}

			name := strings.TrimPrefix(path, "/")
			if name == "" {
				name = "index.html"
			}

			if f, err := webFS.Open(name); err == nil {
				f.Close()
				if strings.HasPrefix(path, "/assets/") {
					// Vite hashes asset filenames — safe to cache forever
					w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
				} else if strings.HasSuffix(name, ".html") {
					w.Header().Set("Cache-Control", "no-cache")
				}
				http.ServeFileFS(w, req, webFS, name)
				return
			}

			// Unknown path → SPA fallback (client routes like /books/:id, /admin)
			w.Header().Set("Cache-Control", "no-cache")
			http.ServeFileFS(w, req, webFS, "index.html")
		})
		log.Println("Serving embedded frontend")
	} else {
		log.Println("Frontend embedding disabled (dist/index.html not found).")
		log.Println("Run: cd frontend && npm run dev")
	}

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: r,
		// Guard against slowloris / hanging connections
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      120 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("Starting API server on %s", addr)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func seedAdmin(db *store.DB, cfg *config.Config) error {
	count, err := db.CountUsers(context.Background())
	if err != nil {
		return err
	}

	// Skip if users already exist
	if count > 0 {
		return nil
	}

	// Never seed a fresh install with missing/default credentials
	if cfg.IsInsecureAdminPassword() {
		return fmt.Errorf("ADMIN_PASSWORD is not set or still the insecure default: set ADMIN_PASSWORD (e.g. in backend/.env) before the first start")
	}

	// Hash password
	hash, err := auth.HashPassword(cfg.AdminPassword)
	if err != nil {
		return err
	}

	admin := &store.User{
		ID:           utils.GenerateShortID(),
		Email:        cfg.AdminEmail,
		PasswordHash: hash,
		Role:         "admin",
	}

	if err := db.CreateUser(context.Background(), admin); err != nil {
		return err
	}

	log.Printf("Admin user created: %s", admin.Email)
	return nil
}
