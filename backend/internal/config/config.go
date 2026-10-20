package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port          int
	DataDir       string
	DatabaseURL   string
	SessionSecret string
	MaxUploadMB   int
	AdminEmail    string
	AdminPassword string
}

func Load() *Config {
	loadDotEnv(".env")

	port := 8080
	if p := os.Getenv("PORT"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			port = parsed
		}
	}

	maxUpload := 25
	if m := os.Getenv("MAX_UPLOAD_MB"); m != "" {
		if parsed, err := strconv.Atoi(m); err == nil {
			maxUpload = parsed
		}
	}

	return &Config{
		Port:          port,
		DataDir:       getEnv("DATA_DIR", "./data"),
		DatabaseURL:   getEnv("DATABASE_URL", "./data/gallery.db"),
		SessionSecret: getEnv("SESSION_SECRET", ""),
		MaxUploadMB:   maxUpload,
		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@example.com"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
	}
}

// IsInsecureAdminPassword reports whether the configured admin password is
// missing or still the well-known insecure default.
func (c *Config) IsInsecureAdminPassword() bool {
	return c.AdminPassword == "" || c.AdminPassword == "admin123"
}

// loadDotEnv reads a KEY=VALUE .env file without overriding variables that
// are already set in the environment.
func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			continue
		}
		if len(value) >= 2 && (value[0] == '"' && value[len(value)-1] == '"' ||
			value[0] == '\'' && value[len(value)-1] == '\'') {
			value = value[1 : len(value)-1]
		}
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
