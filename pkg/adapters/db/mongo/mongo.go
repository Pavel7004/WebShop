package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/Pavel7004/Common/tracing"
	"github.com/Pavel7004/WebShop/pkg/adapters/db"
	"github.com/Pavel7004/WebShop/pkg/adapters/db/mongo/models"
	"github.com/Pavel7004/WebShop/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ db.DB = (*DB)(nil)

type DB struct {
	client *mongo.Client

	collectionItems *mongo.Collection
	collectionUsers *mongo.Collection
}

func New() *DB {
	db := new(DB)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	db.client = client
	db.collectionItems = client.Database("shop").Collection("items")
	db.collectionUsers = client.Database("shop").Collection("users")

	return db
}

func (db *DB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return db.client.Disconnect(ctx)
}
