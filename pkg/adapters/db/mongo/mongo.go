package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/Pavel7004/Common/tracing"
	"github.com/Pavel7004/WebShop/pkg/adapters/db"
	"github.com/Pavel7004/WebShop/pkg/adapters/db/mongo/models"
	"github.com/Pavel7004/WebShop/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ db.DB = (*DB)(nil)

type DB struct {
	client *mongo.Client

	collectionItems *mongo.Collection
	collectionUsers *mongo.Collection
}

func New() *DB {
	db := new(DB)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	db.client = client
	db.collectionItems = client.Database("shop").Collection("items")
	db.collectionUsers = client.Database("shop").Collection("users")

	return db
}

func (db *DB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return db.client.Disconnect(ctx)
}

func (db *DB) RegisterUser(ctx context.Context, user *domain.RegisterUserRequest) (string, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
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

func (db *DB) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	span.SetTag("user_id", id)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidId
	}

	span.SetTag("object_id", objectID)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
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
