package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var dbClient *MongoDBClient

func NewClient(config *DatabaseConfig) (*MongoDBClient, error) {
	var err error
	var client *mongo.Client

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opt := options.Client().
		ApplyURI(config.Uri).
		SetMinPoolSize(5).
		SetMaxPoolSize(100).
		SetMaxConnIdleTime(10 * time.Minute)

	client, err = mongo.Connect(ctx, opt)
	if err != nil {
		return nil, err
	}

	db := client.Database(config.DBName)
	dbClient = &MongoDBClient{
		config: config,
		client: client,
		db:     db,
	}
	dbClient.Ping(ctx)
	return dbClient, nil
}

func GetDB() *MongoDBClient {
	if dbClient == nil {
		panic("MongoDBClient is not initialized!")
	}
	return dbClient
}

type MongoDBClient struct {
	config *DatabaseConfig
	client *mongo.Client
	db     *mongo.Database
}

func (c *MongoDBClient) String() string {
	return fmt.Sprintf("MongoDBClient{config: %+v}", c.config)
}

func (c *MongoDBClient) Shutdown(ctx context.Context) error {
	if err := c.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("error disconnecting from MongoDB: %w\n", err)
	}
	return nil
}

func (c *MongoDBClient) Ping(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := c.client.Ping(ctx, nil); err != nil {
		log.Fatalf("failed to ping MongoDB: %v\n", err.Error())
	}
}

func (c *MongoDBClient) Collection(name string) *mongo.Collection {
	return c.db.Collection(name)
}

func (c *MongoDBClient) ListDatabases(ctx context.Context) ([]string, error) {
	databases, err := c.client.ListDatabaseNames(ctx, bson.D{{"empty", false}})
	if err != nil {
		return nil, err
	}
	return databases, nil
}

func (c *MongoDBClient) CreateUniqueIndex(coll *mongo.Collection, field string, ctx context.Context) error {
	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{field: 1},
		Options: options.Index().SetUnique(true),
	})
	return err
}
