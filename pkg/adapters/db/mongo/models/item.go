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

func (it *Item) ConvertToDomainItem() *domain.Item {
	return &domain.Item{
		ID:          it.ID.Hex(),
		Name:        it.Name,
		Description: it.Description,
		Price:       it.Price,
		CreatedAt:   it.CreatedAt,
	}
}
