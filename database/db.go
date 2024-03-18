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
	_ "github.com/lib/pq"
)

func GetDBConfig() (*pgxpool.Config) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

    const defaultMaxConns = int32(4)
    const defaultMinConns = int32(0)
    const defaultMaxConnLifetime = time.Hour
    const defaultMaxConnIdleTime = time.Minute * 30
    const defaultHealthCheckPeriod = time.Minute
    const defaultConnectTimeout = time.Second * 5

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	var databaseUrl string
	if os.Getenv("ENV") != "production" {
        databaseUrl = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	} else {
        databaseUrl = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=verify-full&sslrootcert=ap-southeast-1-bundle.pem", user, password, host, port, dbname)
	}

    databaseConfig, err := pgxpool.ParseConfig(databaseUrl)
    if err!=nil {
     log.Fatal("Failed to create a config, error: ", err)
    }
   
    databaseConfig.MaxConns = defaultMaxConns
    databaseConfig.MinConns = defaultMinConns
    databaseConfig.MaxConnLifetime = defaultMaxConnLifetime
    databaseConfig.MaxConnIdleTime = defaultMaxConnIdleTime
    databaseConfig.HealthCheckPeriod = defaultHealthCheckPeriod
    databaseConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

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
