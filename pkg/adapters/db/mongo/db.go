package mongo

import (
	"context"

	"github.com/Pavel7004/WebShop/pkg/adapters/db"
	"github.com/Pavel7004/WebShop/pkg/infra/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var _ db.DB = (*DB)(nil)

type DB struct {
	client *mongo.Client
	cfg    *config.MongoCfg

	collectionItems *mongo.Collection
	collectionUsers *mongo.Collection
}

func New(cfg *config.Config) *DB {
	db := new(DB)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Mongo.Timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.Uri))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		panic(err)
	}

	db.client = client
	db.cfg = &cfg.Mongo
	db.collectionItems = client.Database("shop").Collection("items")
	db.collectionUsers = client.Database("shop").Collection("users")

	return db
}

func (db *DB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), db.cfg.Timeout)
	defer cancel()

	return db.client.Disconnect(ctx)
}
