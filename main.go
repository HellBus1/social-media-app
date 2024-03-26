package main

import (
	"context"
	"fmt"
	"log"
	"social-media-app/database"
	"social-media-app/router"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	PORT = ":8080"
)

func main() {
	// Create a new database connection pool with the configured settings
	DB, err := pgxpool.NewWithConfig(context.Background(), database.GetDBConfig())
	if err != nil {
		log.Fatal("Error creating database connection pool:", err)
	}
	defer DB.Close()

	// Ping the database to ensure connectivity
	if err := DB.Ping(context.Background()); err != nil {
		log.Fatal("Could not ping database:", err)
	}
	fmt.Println("Connected to the database!!")

	// Start the application router with the database connection pool
	r := router.StartApp(DB)
	r.Run(PORT)
}
