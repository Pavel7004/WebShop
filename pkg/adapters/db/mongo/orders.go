package mongo

import (
	"context"
	"errors"

	"github.com/Pavel7004/Common/tracing"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Pavel7004/WebShop/pkg/adapters/db/mongo/models"
	"github.com/Pavel7004/WebShop/pkg/domain"
)

func (db *DB) CreateOrder(ctx context.Context, reqDom *domain.CreateOrderRequest) (string, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	req, err := models.ConvertAddOrderRequestFromDomain(reqDom)
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("event", "error", "message", err.Error())
		log.Error().Err(err).Msg("Error in convertion")
		return "", err
	}

	itemIDs := make([]primitive.ObjectID, len(req.Items))
	itemQuantities := make([]uint64, len(req.Items))

	for i, item := range req.Items {
		itemIDs[i] = item.ID
		itemQuantities[i] = item.Quantity
	}
	log.Debug().Msgf("%v ;; %v", itemIDs, itemQuantities)

	pipeline := mongo.Pipeline{
		// Get database items that are presented in order request.
		{primitive.E{Key: "$match", Value: bson.M{"_id": bson.M{"$in": itemIDs}}}},
		// Populate documents with bought goods amount.
		// $indexOfArray based on record's ID returning an index from itemIDs,
		// which is the same as the index from itemQuantities.
		{
			primitive.E{
				Key: "$addFields",
				Value: bson.M{
					"quantity_bought": bson.M{
						"$arrayElemAt": bson.A{
							itemQuantities,
							bson.M{"$indexOfArray": bson.A{itemIDs, "$_id"}},
						},
					},
				},
			},
		},
		// Calculate final price.
		{
			primitive.E{
				Key: "$group",
				Value: bson.M{
					"_id": nil,
					"total": bson.M{
						"$sum": bson.M{"$multiply": bson.A{"$price", "$quantity_bought"}},
					},
				},
			},
		},
	}
	cursor, err := db.collectionItems.Aggregate(ctx, pipeline)
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("event", "error", "message", err.Error())
		log.Error().Err(err).Msg("Error in aggregation")
		return "", err
	}
	defer cursor.Close(ctx)

	log.Debug().Msgf("Cursor = %v", cursor)
	var result struct {
		Total float64 `bson:"total"`
	}
	if err := cursor.Decode(&result); err != nil {
		span.SetTag("error", true)
		span.LogKV("event", "error", "message", err.Error())
		log.Error().Err(err).Msg("Error in Decode")
		return "", err
	}

	req.Total = result.Total
	res, err := db.collectionOrders.InsertOne(ctx, req)
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("event", "error", "message", err.Error())
		log.Error().Err(err).Msg("Error in InsertOne")
		return "", err
	}

	obj, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		span.SetTag("error", true)
		span.LogKV("event", "error", "message", domain.ErrInvalidId.Error())
		log.Error().Err(domain.ErrInvalidId).Msg("Result isn't object id.")
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

func (db *DB) UpdateOrder(ctx context.Context, id string, ord domain.UpdateOrderRequest) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	req, err := models.ConvertUpdateOrderReqToBSON(&ord)
	if err != nil {
		return err
	}

	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidId
	}

	res, err := db.collectionOrders.UpdateByID(ctx, obj, req)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.ErrItemNotFound
		}

		return err
	}

	if res.ModifiedCount < 1 {
		return domain.ErrItemNotFound
	}

	return nil
}
