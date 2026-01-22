package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// App
	Port      string
	JWTSecret string

	// Database
	DatabaseURL string

	// Google OAuth
	GoogleClientID string

	// SMTP Email
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
	SMTPFrom string

	// Web Push (VAPID)
	VapidPublicKey  string
	VapidPrivateKey string
	VapidSubject    string
}

func Load() Config {
	// load .env file kalau ada
	if err := godotenv.Load(); err != nil {
		log.Println("No .env found, using system env")
	}

	return Config{
		// App
		Port:      getEnv("APP_PORT", "3000"),
		JWTSecret: getEnv("JWT_SECRET", "secret"),

		// Database
		DatabaseURL: getEnv("DATABASE_URL", ""),

		// Google
		GoogleClientID: getEnv("GOOGLE_CLIENT_ID", ""),

		// SMTP
		SMTPHost: getEnv("SMTP_HOST", ""),
		SMTPPort: getEnv("SMTP_PORT", "587"),
		SMTPUser: getEnv("SMTP_USER", ""),
		SMTPPass: getEnv("SMTP_PASS", ""),
		SMTPFrom: getEnv("SMTP_FROM", ""),

		// VAPID Push
		VapidPublicKey:  getEnv("VAPID_PUBLIC_KEY", ""),
		VapidPrivateKey: getEnv("VAPID_PRIVATE_KEY", ""),
		VapidSubject:    getEnv("VAPID_SUBJECT", "mailto:example@example.com"),
	}
}

func getEnv(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
