package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Pavel7004/WebShop/pkg/domain"
)

type OrderItem struct {
	ID       primitive.ObjectID `bson:"item_id"`
	Quantity uint64             `bson:"quantity"`
}

type Order struct {
	Items      []OrderItem        `bson:"items"`
	CustomerID primitive.ObjectID `bson:"customer_id"`

	ID        primitive.ObjectID `bson:"_id"`
	Total     float64            `bson:"total"`
	CreatedAt time.Time          `bson:"created_at"`
	Status    domain.StatusID    `bson:"status"`
}

// `bson:"-"`
// `bson:"status,omitempty"`

func ConvertAddOrderRequestFromDomain(ord *domain.CreateOrderRequest) (*Order, error) {
	if ord == nil {
		return nil, domain.ErrNoOrder
	}

	itemIDs := make([]OrderItem, 0, len(ord.Items))
	for _, it := range ord.Items {
		obj, err := primitive.ObjectIDFromHex(it.ID)
		if err != nil {
			return nil, err
		}

		itemIDs = append(itemIDs, OrderItem{
			ID:       obj,
			Quantity: it.Quantity,
		})
	}

	customer, err := primitive.ObjectIDFromHex(ord.CustomerID)
	if err != nil {
		return nil, domain.ErrInvalidId
	}

	return &Order{
		Items:      itemIDs,
		CustomerID: customer,
		ID:         primitive.NewObjectID(),
		Total:      0,
		CreatedAt:  time.Now(),
		Status:     domain.CREATED,
	}, nil
}

func ConvertUpdateOrderReqToBSON(ord *domain.UpdateOrderRequest) (bson.M, error) {
	if ord == nil {
		return nil, domain.ErrNoUpdate
	}

	req := bson.M{}

	if ord.Items != nil {
		itemIDs := make([]OrderItem, 0, len(*ord.Items))
		for _, it := range *ord.Items {
			obj, err := primitive.ObjectIDFromHex(it.ID)
			if err != nil {
				return nil, err
			}

			itemIDs = append(itemIDs, OrderItem{
				ID:       obj,
				Quantity: it.Quantity,
			})
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
	items := make([]domain.OrderItem, 0, len(o.Items))
	for _, it := range o.Items {
		items = append(items, domain.OrderItem{
			ID:       it.ID.Hex(),
			Quantity: it.Quantity,
		})
	}

	return &domain.Order{
		ID:         o.ID.Hex(),
		Total:      o.Total,
		Items:      items,
		CreatedAt:  o.CreatedAt,
		Status:     o.Status,
		CustomerID: o.CustomerID.Hex(),
	}
}
