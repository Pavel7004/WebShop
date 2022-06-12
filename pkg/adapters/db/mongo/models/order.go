package models

import (
	"time"

	"github.com/Pavel7004/WebShop/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID         primitive.ObjectID   `bson:"_id"`
	ItemIDs    []primitive.ObjectID `bson:"item_ids"`
	Total      float64              `bson:"total"`
	CreatedAt  time.Time            `bson:"created_at"`
	Status     domain.StatusID      `bson:"status"`
	CustomerID primitive.ObjectID   `bson:"customer_id"`
}

// `bson:"-"`
// `bson:"status,omitempty"`

func ConvertAddOrderRequestFromDomain(ord *domain.CreateOrderRequest) (*Order, error) {
	if ord == nil {
		return nil, domain.ErrNoOrder
	}

	itemIDs := make([]primitive.ObjectID, 0, len(ord.ItemIDs))
	for _, it := range ord.ItemIDs {
		obj, err := primitive.ObjectIDFromHex(it)
		if err != nil {
			return nil, err
		}

		itemIDs = append(itemIDs, obj)
	}

	customer, err := primitive.ObjectIDFromHex(ord.CustomerID)
	if err != nil {
		return nil, domain.ErrInvalidId
	}

	return &Order{
		ID:         primitive.NewObjectID(),
		ItemIDs:    itemIDs,
		CreatedAt:  time.Now(),
		Status:     domain.CREATED,
		CustomerID: customer,
	}, nil
}

func ConvertUpdateOrderReqToBSON(ord *domain.UpdateOrderRequest) (bson.M, error) {
	if ord == nil {
		return nil, domain.ErrNoUpdate
	}

	req := bson.M{}

	if ord.ItemIDs != nil {
		itemIDs := make([]primitive.ObjectID, 0, len(*ord.ItemIDs))
		for _, it := range *ord.ItemIDs {
			obj, err := primitive.ObjectIDFromHex(it)
			if err != nil {
				return nil, err
			}

			itemIDs = append(itemIDs, obj)
		}

		req["item_ids"] = itemIDs
	}

	if ord.Status != nil {
		req["status"] = *ord.Status
	}

	req = bson.M{"$set": req}
	return req, nil
}

func (o *Order) ConvertToDomain() *domain.Order {
	itemIDs := make([]string, 0, len(o.ItemIDs))
	for _, id := range o.ItemIDs {
		itemIDs = append(itemIDs, id.Hex())
	}

	return &domain.Order{
		ID:         o.ID.Hex(),
		ItemIDs:    itemIDs,
		CreatedAt:  o.CreatedAt,
		Status:     o.Status,
		CustomerID: o.CustomerID.Hex(),
	}
}
