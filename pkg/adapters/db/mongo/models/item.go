package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Pavel7004/WebShop/pkg/domain"
)

type Item struct {
	ID          primitive.ObjectID `bson:"_id"`
	OwnerID     primitive.ObjectID `bson:"owner_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"desc"`
	Category    string             `bson:"category"`
	Price       float64            `bson:"price"`
	CreatedAt   time.Time          `bson:"created_at"`
	Quantity    uint64             `bson:"quantity"`
}

func ConvertItemFromDomainRequest(it *domain.AddItemRequest) (*Item, error) {
	if it == nil {
		return nil, domain.ErrItemNotFound // TODO: fix no item error
	}

	ownerID, err := primitive.ObjectIDFromHex(it.OwnerID)
	if err != nil {
		return nil, domain.ErrInvalidId
	}

	return &Item{
		ID:          primitive.NewObjectID(),
		OwnerID:     ownerID,
		Name:        it.Name,
		Description: it.Description,
		Category:    it.Category,
		Price:       it.Price,
		CreatedAt:   time.Now(),
		Quantity:    it.Quantity,
	}, nil
}

func ConvertItemFromDomain(it *domain.Item) (*Item, error) {
	id, err := primitive.ObjectIDFromHex(it.ID)
	if err != nil {
		return nil, domain.ErrInvalidId
	}

	ownerId, err := primitive.ObjectIDFromHex(it.OwnerID)
	if err != nil {
		return nil, domain.ErrInvalidId
	}

	return &Item{
		ID:          id,
		OwnerID:     ownerId,
		Name:        it.Name,
		Description: it.Description,
		Category:    it.Category,
		Price:       it.Price,
		CreatedAt:   it.CreatedAt,
		Quantity:    it.Quantity,
	}, nil
}

func (it *Item) ConvertToDomain() *domain.Item {
	return &domain.Item{
		ID:          it.ID.Hex(),
		OwnerID:     it.OwnerID.Hex(),
		Name:        it.Name,
		Description: it.Description,
		Category:    it.Category,
		Price:       it.Price,
		CreatedAt:   it.CreatedAt,
		Quantity:    it.Quantity,
	}
}

func ConvertItemsToDomain(items []Item) []*domain.Item {
	result := make([]*domain.Item, 0, len(items))

	for _, it := range items {
		result = append(result, it.ConvertToDomain())
	}

	return result
}

func ConvertUpdateReqToBSON(in *domain.UpdateItemRequest) (bson.M, error) {
	if in == nil {
		return nil, domain.ErrNoUpdate
	}
	req := bson.M{}
	if in.Name != nil {
		req["name"] = in.Name
	}
	if in.Description != nil {
		req["desc"] = in.Description
	}
	if in.Category != nil {
		req["category"] = in.Category
	}
	if in.OwnerID != nil {
		req["owner_id"] = in.OwnerID
	}
	if in.Price != nil {
		req["price"] = in.Price
	}
	if in.Quantity != nil {
		req["quantity"] = in.Quantity
	}
	req = bson.M{"$set": req}
	return req, nil
}
