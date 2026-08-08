package main

import (
	"database/sql"
	"time"
)

// DatabaseConfig: hold the connection pool config
type DatabaseConfig struct {
	Driver          string // holds info for what kind of database it is
	DSN             string // Data source name, reps the url of the db
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxIdleTime time.Duration // prevents ideal connection from occupying the network for extended periods
	ConnMaxLifetime time.Duration
}

// DefaultDatabaseConfig: returns the recommended config
// based on load testing evidence
func DefaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		MaxOpenConns:    30,
		MaxIdleConns:    30,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 30 * time.Second,
	}
}

func initDB(cfg *DatabaseConfig) (*sql.DB, error) {
	// create conn pool
	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, err
	}

	// configure the conn pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}
