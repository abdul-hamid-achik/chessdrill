package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Client struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewClient(ctx context.Context, uri, dbName string) (*Client, error) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB")

	return &Client{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

func (c *Client) Database() *mongo.Database {
	return c.database
}

func (c *Client) Collection(name string) *mongo.Collection {
	return c.database.Collection(name)
}

func (c *Client) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

// CreateIndexes creates all required indexes for the application
func (c *Client) CreateIndexes(ctx context.Context) error {
	// Users collection indexes
	usersCollection := c.Collection("users")
	_, err := usersCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create users indexes: %w", err)
	}

	// Sessions collection indexes
	sessionsCollection := c.Collection("sessions")
	_, err = sessionsCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys:    bson.D{{Key: "token", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "expires_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create sessions indexes: %w", err)
	}

	// Drill sessions collection indexes
	drillSessionsCollection := c.Collection("drill_sessions")
	_, err = drillSessionsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create drill_sessions indexes: %w", err)
	}

	// Attempts collection indexes
	attemptsCollection := c.Collection("attempts")
	_, err = attemptsCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "session_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "drill_type", Value: 1},
				{Key: "answered_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "correct_answer", Value: 1},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create attempts indexes: %w", err)
	}

	log.Println("MongoDB indexes created successfully")
	return nil
}
