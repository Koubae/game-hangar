package models

import (
	"context"
	"github.com/koubae/game-hangar/account/internal/infrastructure/database/models"
	"github.com/koubae/game-hangar/account/internal/settings"
	"github.com/koubae/game-hangar/account/pkg/database/mongodb"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/koubae/game-hangar/account/pkg/testings"
)

func TestMain(m *testing.M) {
	// /////////////////////////
	//			SetUp
	// /////////////////////////
	settings.NewConfig()

	config := settings.GetConfig()
	log.Println(config.DatabaseConfig)

	client, err := mongodb.NewClient(config.DatabaseConfig)
	if err != nil {
		panic(err.Error())
	}
	log.Println(client)

	collection := client.Collection("accounts")

	// /////////////////////////
	//			Tests
	// /////////////////////////

	code := m.Run()

	// /////////////////////////
	//			CleanUp
	// /////////////////////////
	// Cleanup after all tests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = collection.Drop(ctx)

	if err := client.Shutdown(ctx); err != nil {
		log.Fatalf("MongoDB error while shutting Down, error %v\n", err)
	}
	log.Println("MongoDB shutdown completed")

	os.Exit(code)
}

const HashedPassword = "$2a$10$GPeYnQMl9mGX1hvIrqTIjeJmPOESnUHFe39Ksm0HifPU8r9YchbbC"

func TestAccount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	account := models.Account{Username: "integration-tests", Password: HashedPassword}

	t.Run("nil values expected on a struct that is not saved in the database", func(t *testing.T) {
		assert.Nil(t, account.ID)
		assert.Nil(t, account.Created)
		assert.Nil(t, account.Updated)
	})

	t.Run("Create", func(t *testing.T) {
		db := mongodb.GetDB()
		collection := db.Collection("accounts")

		account.OnCreate() // Initialize updated/created fields
		result, err := collection.InsertOne(ctx, account)
		assert.NoError(t, err)

		account.OnCreated(result)

		expectedID := result.InsertedID
		assert.Equal(t, expectedID, *account.ID)
	})

}
