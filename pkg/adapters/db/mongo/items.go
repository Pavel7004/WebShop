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

	ctx, cancel := context.WithTimeout(ctx, db.cfg.Timeout)
	defer cancel()

	it, err := models.ConvertItemFromDomainRequest(item)
	if err != nil {
		return "", err
	}

	res, err := db.collectionItems.InsertOne(ctx, it)

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

func (db *DB) UpdateItem(ctx context.Context, id string, in *domain.UpdateItemRequest) (int64, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, domain.ErrInvalidId
	}

	req, err := models.ConvertUpdateReqToBSON(in)
	if err != nil {
		return 0, err
	}

	res, err := db.collectionItems.UpdateOne(ctx, bson.M{"_id": userID}, req)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, domain.ErrItemNotFound
		}

		return 0, err
	}

	return res.ModifiedCount, nil
}

func (db *DB) GetItemsTotalCost(ctx context.Context, itemIDs []string) (float64, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	itemObjIDs := make([]primitive.ObjectID, 0, len(itemIDs))
	for _, id := range itemIDs {
		obj, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, domain.ErrInvalidId
		}

		itemObjIDs = append(itemObjIDs, obj)
	}

	items, err := db.findItems(ctx, bson.M{"items": bson.M{
		"$in": iteitemObjIDs,
	}})
	if err != nil {
		return 0, err
	}

	var result float64
	for _, item := range items {
		result += item.Price
	}

	return result, nil
}
