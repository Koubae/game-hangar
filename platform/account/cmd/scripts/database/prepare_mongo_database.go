package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/koubae/game-hangar/account/internal/settings"
	"github.com/koubae/game-hangar/account/pkg/database/mongodb"
	"log"
	"time"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	settings.NewConfig()

	log.SetFlags(log.Ldate | log.Ltime)

}

func main() {
	config := settings.GetConfig()
	log.Println(config.DatabaseConfig)

	client, err := mongodb.NewClient(config.DatabaseConfig)
	if err != nil {
		panic(err.Error())
	}
	log.Println(client)

	databases, err := client.ListDatabases(context.Background())
	if err != nil {
		log.Printf("MongoDB error while listing databases, error %v\n", err)
	}
	log.Printf("MongoDB databases: %v\n", databases)

	log.Println("MongoDB Creating index for Account collection")

	collectionAccount := client.Collection("accounts")
	err = client.CreateUniqueIndex(collectionAccount, "username", context.Background())
	if err != nil {
		log.Printf("MongoDB error while creating index for Account collection, error %v\n", err)
	}
	shutdown(client)
}

func shutdown(client *mongodb.MongoDBClient) {
	shutDownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Shutdown(shutDownCtx); err != nil {
		log.Fatalf("MongoDB error while shutting Down, error %v\n", err)
	}
	log.Println("MongoDB shutdown completed")
}
