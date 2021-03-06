package mongo

import (
	"context"
	"errors"

	"github.com/Pavel7004/Common/tracing"
	"github.com/Pavel7004/WebShop/pkg/adapters/db/mongo/models"
	"github.com/Pavel7004/WebShop/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *DB) CreateOrder(ctx context.Context, reqDom *domain.CreateOrderRequest) (string, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	req, err := models.ConvertAddOrderRequestFromDomain(reqDom)
	if err != nil {
		return "", err
	}

	res, err := db.collectionOrders.InsertOne(ctx, req)
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

func (db *DB) GetOrderInfo(ctx context.Context, id string) (*domain.Order, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidId
	}

	ctx, cancel := context.WithTimeout(ctx, db.cfg.Timeout)
	defer cancel()

	var result models.Order
	if err := db.collectionItems.FindOne(ctx, bson.M{"_id": obj}).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrOrderNotFound
		}

		return nil, err
	}

	return result.ConvertToDomain(), nil
}

func (db *DB) UpdateOrder(ctx context.Context, id string, ord domain.UpdateOrderRequest) (int64, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	req, err := models.ConvertUpdateOrderReqToBSON(&ord)
	if err != nil {
		return 0, err
	}

	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, domain.ErrInvalidId
	}

	res, err := db.collectionOrders.UpdateByID(ctx, obj, req)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, domain.ErrItemNotFound
		}

		return 0, err
	}

	return res.ModifiedCount, nil
}
