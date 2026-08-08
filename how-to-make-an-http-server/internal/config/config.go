package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string

	MaxIdleConns int
	MaxOpenConns int

	MaxIdleTime time.Duration
	MaxLifeTime time.Duration
}

type ServerConfig struct {
	Host              string
	Port              string
	MaxHeaderByte     int
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	ShutdownTimeout   time.Duration
}

func Load() Config {
	return Config{
		Server: ServerConfig{
			Host:            getEnv("HOST", ""),
			Port:            getEnv("PORT", ":8080"),
			ReadTimeout:     getDuration("READ_TIMEOUT", 5*time.Second),
			WriteTimeout:    getDuration("WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:     getDuration("IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: getDuration("SHUTDOWN_TIMEOUT", 15*time.Second),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			User:         getEnv("DB_USER", "postgres"),
			Password:     getEnv("DB_PASSWORD", "postgres"),
			Name:         getEnv("DB_NAME", "user_api"),
			SSLMode:      getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns: getInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getInt("DB_MAX_IDLE_CONNS", 25),
			MaxIdleTime: getDuration(
				"DB_MAX_IDLE_TIME",
				15*time.Minute,
			),
			MaxLifeTime: getDuration(
				"DB_MAX_LIFETIME",
				time.Hour,
			),
		},
	}
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		seconds, err := strconv.Atoi(value)
		if err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	return fallback
}

func getInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		convertedV, err := strconv.Atoi(value)
		if err == nil {
			return convertedV
		}
	}
	return fallback
}
