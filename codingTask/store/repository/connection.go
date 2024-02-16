package repository

import (
	"database/sql"
	"fmt"
)

const driverName = "postgres"

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open(driverName, fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("postgres: ping db: %w", err)
	}
	return db, nil
}
