package main

import (
	"fmt"
	"log"
	"social-media-app/database"
	"social-media-app/router"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	PORT = ":8000"
)

func main() {
	DB, configuringError := pgxpool.NewWithConfig(context.Background(), database.GetDBConfig())
	if configuringError != nil {
		log.Fatal("Error while creating connection to the database!!")
	} 
 
	connection, acquiringConnectionError := DB.Acquire(context.Background())
	if acquiringConnectionError != nil {
		log.Fatal("Error while acquiring connection from the database pool!!")
	} 
	defer connection.Release()
 
	pingError := connection.Ping(context.Background())
	if pingError != nil {
		log.Fatal("Could not ping database")
	}
	fmt.Println("Connected to the database!!")
	
	r := router.StartApp(DB)
	
	r.Run(PORT)

	defer DB.Close()
}