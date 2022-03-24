package models

import (
	"time"

	"github.com/Pavel7004/WebShop/pkg/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"desc"`
	Price       float64            `bson:"price"`
	CreatedAt   time.Time          `bson:"created_at"`
}

func (it *Item) ConvertToDomain() *domain.Item {
	return &domain.Item{
		ID:          it.ID.Hex(),
		Name:        it.Name,
		Description: it.Description,
		Price:       it.Price,
		CreatedAt:   it.CreatedAt,
	}
}

func ConvertItemsToDomain(items []Item) []*domain.Item {
	result := make([]*domain.Item, 0, len(items))

	for _, it := range items {
		result = append(result, it.ConvertToDomain())
	}

	return result
}
