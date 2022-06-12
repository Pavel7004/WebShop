package models

import (
	"time"

	"github.com/Pavel7004/WebShop/pkg/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Phone     string             `bson:"phone"`
	CreatedAt time.Time          `bson:"created_at"`
	Balance   uint64             `bson:"balance"`
}

func (user *User) ConvertToDomain() *domain.User {
	return &domain.User{
		ID:        user.ID.Hex(),
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		Balance:   user.Balance,
	}
}

func ConvertUsersToDomain(users []User) []*domain.User {
	result := make([]*domain.User, 0, len(users))

	for _, user := range users {
		result = append(result, user.ConvertToDomain())
	}

	return result
}
