package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Server
	ServerPort  string
	Environment string

	// JWT
	JWTSecret          string
	JWTExpirationHours int
	JWTRefreshDays     int

	OpenIaKey string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	AppConfig = &Config{
		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "saas_medico"),

		// Server
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),

		// JWT
		JWTSecret:          getEnv("JWT_SECRET", "your-super-secret-key-change-in-production"),
		JWTExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		JWTRefreshDays:     getEnvAsInt("JWT_REFRESH_DAYS", 7),

		OpenIaKey: getEnv("OPEN_AI_API_KEY", "ANOTHER-API-KEY"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
