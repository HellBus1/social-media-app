package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// GetDBConfig returns the configuration for the PostgreSQL database connection pool.
func GetDBConfig() *pgxpool.Config {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Set default connection pool configuration
	const (
		defaultMaxConns        = 20
		defaultMinConns        = 5
		defaultMaxConnLifetime = time.Minute * 30
		defaultMaxConnIdleTime = time.Minute * 5
		defaultHealthCheck     = time.Minute
		defaultConnectTimeout  = time.Second * 5
	)

	// Parse database connection URL
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	var databaseURL string
	if os.Getenv("ENV") != "production" {
		databaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	} else {
		databaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=verify-full&sslrootcert=ap-southeast-1-bundle.pem", user, password, host, port, dbname)
	}

	// Parse database configuration
	databaseConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatal("Failed to create database config:", err)
	}

	// Apply default connection pool settings
	databaseConfig.MaxConns = defaultMaxConns
	databaseConfig.MinConns = defaultMinConns
	databaseConfig.MaxConnLifetime = defaultMaxConnLifetime
	databaseConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	databaseConfig.HealthCheckPeriod = defaultHealthCheck
	databaseConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	// Set up connection pool event hooks
	databaseConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("Before acquiring the connection pool to the database!!")
		return true
	}

	databaseConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("After releasing the connection pool to the database!!")
		return true
	}

	databaseConfig.BeforeClose = func(c *pgx.Conn) {
		log.Println("Closed the connection pool to the database!!")
	}

	return databaseConfig
}