package config

import "os"

type Config struct {
	Port      string
	DBDSN     string
	JWTSecret string
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		DBDSN:     getEnv("DB_DSN", "root:web3pass@tcp(127.0.0.1:3306)/web3_learning?parseTime=true&charset=utf8mb4"),
		JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-in-production"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
