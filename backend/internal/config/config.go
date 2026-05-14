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
		DBDSN:     getEnv("DB_DSN", "gotoweb3:Goto@Web3@tcp(127.0.0.1:3306)/goto_web3?parseTime=true&charset=utf8mb4"),
		JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-in-production"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
