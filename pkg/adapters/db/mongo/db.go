package mongo

import (
	"context"

	"github.com/Pavel7004/Common/tracing"
	"github.com/Pavel7004/WebShop/pkg/adapters/db"
	"github.com/Pavel7004/WebShop/pkg/adapters/db/mongo/models"
	"github.com/Pavel7004/WebShop/pkg/domain"
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

func (db *DB) findItems(ctx context.Context, filter interface{}) ([]*domain.Item, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	ctx, cancel := context.WithTimeout(ctx, db.cfg.Timeout)
	defer cancel()

	cur, err := db.collectionItems.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []models.Item
	if err := cur.All(ctx, &results); err != nil {
		return nil, err
	}

	return models.ConvertItemsToDomain(results), nil
}
