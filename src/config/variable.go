package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
)

var (
	DatabasePort     string
	DatabaseHost     string
	DatabaseUsername string
	DatabasePassword string
	DatabaseName     string
	JwtSecret        string
	FrontendBaseURL  string
)

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	DatabasePort = getEnv("DATABASE_PORT", "5432")
	DatabaseHost = getEnv("DATABASE_HOST", "localhost")
	DatabaseUsername = getEnv("DATABASE_USERNAME", "pg")
	DatabasePassword = getEnv("DATABASE_PASSWORD", "pass")
	DatabaseName = getEnv("DATABASE_NAME", "crud")
	JwtSecret = getEnv("JWT_SECRET", "secret")
	FrontendBaseURL = getEnv("FRONTEND_BASE_URL", "localhost:5173")

}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	slog.Warn(fmt.Sprintf("ENV %v not found using default value: %v", key, defaultValue))
	return defaultValue
}
