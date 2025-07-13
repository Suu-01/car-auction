// internal/config/config.go
package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DSN        string
	JwtSecret  []byte
	AuctionTTL time.Duration
}

var Cfg *Config

func Load() {
	// .env 파일도 먼저 로드 (선택)
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN is required")
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	ttl, err := strconv.Atoi(os.Getenv("AUCTION_TTL_MINUTES"))
	if err != nil || ttl <= 0 {
		ttl = 60
	}

	Cfg = &Config{
		Port:       port,
		DSN:        dsn,
		JwtSecret:  []byte(secret),
		AuctionTTL: time.Duration(ttl) * time.Minute,
	}
}
