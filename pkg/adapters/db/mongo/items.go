package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/Pavel7004/Common/tracing"
	"github.com/Pavel7004/WebShop/pkg/adapters/db/mongo/models"
	"github.com/Pavel7004/WebShop/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *DB) AddItem(ctx context.Context, item *domain.AddItemRequest) (string, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	ownerId, err := primitive.ObjectIDFromHex(item.OwnerID)
	if err != nil {
		return "", domain.ErrInvalidId
	}

	ctx, cancel := context.WithTimeout(ctx, db.cfg.Timeout)
	defer cancel()

	res, err := db.collectionItems.InsertOne(ctx, bson.M{
		"name":        item.Name,
		"owner_id":    ownerId,
		"price":       item.Price,
		"description": item.Description,
		"created_at":  time.Now(),
	})

	if err != nil {
		return "", err
	}

	obj, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", domain.ErrInvalidId
	}

	span.SetTag("result_id", obj.Hex())

	return obj.Hex(), nil
}

func (db *DB) GetItemById(ctx context.Context, id string) (*domain.Item, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	span.SetTag("item_id", id)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidId
	}

	ctx, cancel := context.WithTimeout(ctx, db.cfg.Timeout)
	defer cancel()

	var result models.Item
	if err := db.collectionItems.FindOne(ctx, bson.M{"_id": objectID}).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrItemNotFound
		}

		return nil, err
	}

	return result.ConvertToDomain(), nil
}

func (db *DB) GetItemsByPrice(ctx context.Context, from, to float64) ([]*domain.Item, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	span.SetTag("from", from)
	span.SetTag("to", to)

	return db.findItems(ctx, bson.M{"price": bson.M{"$gte": from, "$lte": to}})
}

func (db *DB) GetRecentlyAddedItems(ctx context.Context, period time.Duration) ([]*domain.Item, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	span.SetTag("period", period.String())

	timeBound := time.Now().Add(-period)

	return db.findItems(ctx, bson.M{"created_at": bson.M{"$gte": timeBound}})
}
