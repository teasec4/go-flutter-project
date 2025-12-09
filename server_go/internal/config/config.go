// Package config provides configuration loading from environment variables
package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	Server ServerConfig
	DB     DBConfig
}

// ServerConfig holds server-related settings
type ServerConfig struct {
	Addr         string
	ReadTimeout  int // seconds
	WriteTimeout int // seconds
	IdleTimeout  int // seconds
}

// DBConfig holds database-related settings
type DBConfig struct {
	Path string
}

// Load reads configuration from environment variables with sensible defaults
func Load() *Config {
	cfg := &Config{
		Server: ServerConfig{
			Addr:         getEnv("SERVER_ADDR", ":8080"),
			ReadTimeout:  getEnvInt("SERVER_READ_TIMEOUT", 15),
			WriteTimeout: getEnvInt("SERVER_WRITE_TIMEOUT", 15),
			IdleTimeout:  getEnvInt("SERVER_IDLE_TIMEOUT", 60),
		},
		DB: DBConfig{
			Path: getEnv("DB_PATH", "bank.db"),
		},
	}
	return cfg
}

// getEnv reads environment variable with default fallback
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// getEnvInt reads integer environment variable with default fallback
func getEnvInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}
