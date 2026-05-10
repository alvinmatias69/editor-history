package config

import (
	"os"
	"time"
)

type Redis struct {
	Addr      string
	KeyExpiry time.Duration
}

func GetRedisCfg() (*Redis, error) {
	keyExpiry, err := time.ParseDuration(os.Getenv("REDIS_EXPIRY_DURATION"))
	if err != nil {
		return nil, err
	}

	return &Redis{
		Addr:      os.Getenv("REDIS_ADDRESS"),
		KeyExpiry: keyExpiry,
	}, nil
}
