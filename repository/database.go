package repository

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func ReadEnv(key string) string {
	return os.Getenv(key)
}

func Open(config PostgresConfig) (*sql.DB, error) {

	db, err := sql.Open("pgx", config.String())
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	return db, nil
}

func ReadPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     ReadEnv("PGHOST"),
		Port:     ReadEnv("PGPORT"),
		User:     ReadEnv("PGUSER"),
		Password: ReadEnv("PGHOST"),
		Database: ReadEnv("PGDATABASE"),
		SSLMode:  "disable",
	}
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (config PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User,
		config.Password, config.Database, config.SSLMode)
}
