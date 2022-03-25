package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/Pavel7004/Common/tracing"
	"github.com/Pavel7004/WebShop/pkg/adapters/db/mongo/models"
	"github.com/Pavel7004/WebShop/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *DB) RegisterUser(ctx context.Context, user *domain.RegisterUserRequest) (string, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	ctx, cancel := context.WithTimeout(ctx, db.cfg.Timeout)
	defer cancel()

	res, err := db.collectionUsers.InsertOne(ctx, bson.M{
		"name":       user.Name,
		"email":      user.Email,
		"phone":      user.Phone,
		"created_at": time.Now(),
	})

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

func (db *DB) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	span.SetTag("user_id", id)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidId
	}

	ctx, cancel := context.WithTimeout(ctx, db.cfg.Timeout)
	defer cancel()

	var result models.User
	if err := db.collectionUsers.FindOne(ctx, bson.M{"_id": objectID}).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return result.ConvertToDomain(), nil
}

func (db *DB) GetItemsByOwnerId(ctx context.Context, id string) ([]*domain.Item, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	span.SetTag("owner_id", id)

	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidId
	}

	return db.findItems(ctx, bson.M{"owner_id": obj})
}
