package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"desc"`
	Price       float64            `bson:"price"`
	CreatedAt   time.Time          `bson:"created_at"`
}
