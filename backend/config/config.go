package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	DBDriver   string // "postgres" or "sqlite"
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBTimeZone string
	JWTSecret  string
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		return val
	}
	return defaultVal
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		AppPort:    getEnv("PORT", "8080"),
		DBDriver:   getEnv("DB_DRIVER", "auto"), // "auto", "postgres", or "sqlite"
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "R4Tu022804-"),
		DBName:     getEnv("DB_NAME", "db_excelso_work_sample"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		DBTimeZone: getEnv("DB_TIMEZONE", "Asia/Jakarta"),
		JWTSecret:  getEnv("JWT_SECRET", "super-secret-jwt-key-senior-work-sample-2026"),
	}
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.DBSSLMode, c.DBTimeZone,
	)
}
