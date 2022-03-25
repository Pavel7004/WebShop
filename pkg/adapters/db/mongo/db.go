package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/Pavel7004/Common/tracing"
	"github.com/Pavel7004/WebShop/pkg/adapters/db"
	"github.com/Pavel7004/WebShop/pkg/adapters/db/mongo/models"
	"github.com/Pavel7004/WebShop/pkg/domain"
	"github.com/Pavel7004/WebShop/pkg/infra/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ db.DB = (*DB)(nil)

type DB struct {
	client *mongo.Client
	cfg    config.MongoCfg

	collectionItems *mongo.Collection
	collectionUsers *mongo.Collection
}

func New(cfg config.Config) *DB {
	db := new(DB)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.Uri))
	if err != nil {
		panic(err)
	}

	db.client = client
	db.cfg = cfg.Mongo
	db.collectionItems = client.Database("shop").Collection("items")
	db.collectionUsers = client.Database("shop").Collection("users")

	return db
}

func (db *DB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return db.client.Disconnect(ctx)
}
