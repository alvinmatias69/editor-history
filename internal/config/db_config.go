package config

import (
	"os"
)

type DB struct {
	Username string
	Password string
	DbName   string
}

func GetDBCfg() *DB {
	return &DB{
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DbName:   os.Getenv("POSTGRES_DB"),
	}
}
