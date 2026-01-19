package config

import (
	"log"
	"os"
	"strconv"
	"time"

)

// Config holds application configuration
type Config struct {
	Server   ServerConfig
	RateLimit RateLimitConfig
}

type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	RequestTimeout time.Duration
}

type RateLimitConfig struct {
	RequestsPerSecond int
	Burst             int
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 120*time.Second),
			RequestTimeout: getDurationEnv("SERVER_REQUEST_TIMEOUT", 5*time.Second),
		},
		RateLimit: RateLimitConfig{
			RequestsPerSecond: getIntEnv("RATE_LIMIT_RPS", 10),
			Burst:             getIntEnv("RATE_LIMIT_BURST", 20),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
		log.Printf("Invalid integer for %s, using default: %d", key, defaultValue)
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		log.Printf("Invalid duration for %s, using default: %v", key, defaultValue)
	}
	return defaultValue
}

// Addr returns the server address
func (c *ServerConfig) Addr() string {
	return ":" + c.Port
}