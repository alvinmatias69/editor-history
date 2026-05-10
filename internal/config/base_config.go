package config

import (
	"log"
	"os"
	"strconv"
)

type Base struct {
	Port                  string
	DefaultTimeoutSeconds uint64
}

func GetBaseCfg() *Base {
	defaultTimeoutSeconds, err := strconv.ParseUint(os.Getenv("DEFAULT_TIMEOUT_SECONDS"), 10, 64)
	if err != nil {
		log.Fatalf("error while parsing config: %v", err)
	}

	return &Base{
		Port:                  os.Getenv("PORT"),
		DefaultTimeoutSeconds: defaultTimeoutSeconds,
	}
}
