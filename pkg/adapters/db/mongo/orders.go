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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DB) CreateOrder(ctx context.Context, reqDom *domain.CreateOrderRequest) (string, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	req, err := models.ConvertAddOrderRequestFromDomain(reqDom)
	if err != nil {
		return "", err
	}

	ids := make([]primitive.ObjectID, 0, len(req.Items))
	qty := make(map[string]uint64, len(req.Items))
	for _, item := range req.Items {
		ids = append(ids, item.ID)
		qty[item.ID.Hex()] = item.Quantity
	}

	opts := new(options.AggregateOptions)
	opts.SetLet(bson.M{
		"items": qty,
	})

	match := bson.D{primitive.E{
		Key: "$match",
		Value: bson.M{
			"_id": bson.M{"$in": ids},
		},
	}}

	itemsQty := bson.M{
		"$getField": bson.M{
			"$concat": bson.A{
				bson.M{
					"$getField": bson.A{
						"$_id",
					},
				},
				".1",
			},
		},
	}

	group := bson.D{primitive.E{
		Key: "$group",
		Value: bson.M{
			"_id": nil,
			"total": bson.M{
				"$sum": bson.M{
					"$multiply": bson.A{
						"$price",
						itemsQty,
					},
				},
			},
		},
	}}

	cur, err := db.collectionItems.Aggregate(ctx, mongo.Pipeline{match, group}, opts)
	if err != nil {
		return "", err
	}

	var total bson.M
	if err := cur.All(ctx, total); err != nil {
		return "", err
	}
	cur.Close(ctx)

	req.Total = total["total"].(float64)

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
