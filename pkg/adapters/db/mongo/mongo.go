package mongo

import (
	"context"
	"errors"
	"time"

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
}

var (
	ErrInvalidObjectType = errors.New("Returned object is not ObjectID")
)

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

	return db
}

func (db *DB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return db.client.Disconnect(ctx)
}

func (db *DB) GetItemById(ctx context.Context, id string) (*domain.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidId
	}

	var result models.Item
	if err := db.collectionItems.FindOne(ctx, bson.M{"_id": objectID}).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrItemNotFound
		}

		return nil, err
	}

	return &domain.Item{
		ID:          result.ID.Hex(),
		Name:        result.Name,
		Description: result.Description,
		Price:       result.Price,
		CreatedAt:   result.CreatedAt,
	}, nil
}

func (db *DB) AddItem(ctx context.Context, item *domain.AddItemRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	res, err := db.collectionItems.InsertOne(ctx, bson.M{
		"name":        item.Name,
		"price":       item.Price,
		"description": item.Description,
		"created_at":  time.Now(),
	})

	if err != nil {
		return "", err
	}

	resStr, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", ErrInvalidObjectType
	}

	return resStr.Hex(), nil
}

func (db *DB) GetItemsByPrice(ctx context.Context, from, to float64) ([]*domain.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cur, err := db.collectionItems.Find(ctx, bson.M{"price": bson.M{"$gte": from, "$lte": to}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []models.Item
	if err := cur.All(ctx, &results); err != nil {
		return nil, err
	}

	items := make([]*domain.Item, 0, len(results))
	for _, it := range results {
		items = append(items, it.ConvertToDomainItem())
	}

	return items, nil
}

func (db *DB) GetRecentlyAddedItems(ctx context.Context, period time.Duration) ([]*domain.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	timeBound := time.Now().Add(-period)

	cur, err := db.collectionItems.Find(ctx, bson.M{"created_at": bson.M{"$gte": timeBound}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []models.Item
	if err := cur.All(ctx, &results); err != nil {
		return nil, err
	}

	items := make([]*domain.Item, 0, len(results))
	for _, it := range results {
		items = append(items, it.ConvertToDomainItem())
	}

	return items, nil
}
