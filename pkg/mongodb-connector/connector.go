package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connector mongodb instance
type Connector struct {
	DB     string        // The MongoDB uri
	DBName string        // Database name from MongoDB
	client *mongo.Client // Mongodb client
}

var ctx context.Context

// New to create a new Mongodb connection
func New(c *Connector) (*Connector, error) {
	client, e := mongo.Connect(ctx, options.Client().ApplyURI(c.DB))
	if e != nil {
		return &Connector{}, e
	}

	return &Connector{
		DB:     c.DB,
		DBName: c.DBName,
		client: client,
	}, nil
}

// Ping to check connection status
func (c *Connector) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if e := c.client.Ping(ctx, readpref.Primary()); e != nil {
		return e
	}

	return nil
}
