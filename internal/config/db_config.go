package config

import (
	"os"
)

type DB struct {
	Username string
	Password string
	DbName   string
	Host     string
	Port     string
}

func GetDBCfg() *DB {
	return &DB{
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DbName:   os.Getenv("POSTGRES_DB"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
	}
}
