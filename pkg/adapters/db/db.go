package db

import (
	"context"
	"time"

	"github.com/Pavel7004/WebShop/pkg/domain"
)

type DB interface {
	AddItem(ctx context.Context, item *domain.AddItemRequest) (string, error)
	UpdateItem(ctx context.Context, id string, in *domain.UpdateItemRequest) (int64, error)
	GetItemById(ctx context.Context, id string) (*domain.Item, error)
	GetItemsByPrice(ctx context.Context, from, to float64) ([]*domain.Item, error)
	GetRecentlyAddedItems(ctx context.Context, period time.Duration) ([]*domain.Item, error)
	GetItemsByOwnerId(ctx context.Context, id string) ([]*domain.Item, error)

	RegisterUser(ctx context.Context, user *domain.RegisterUserRequest) (string, error)
	GetUserById(ctx context.Context, id string) (*domain.User, error)
	GetRecentlyAddedUsers(ctx context.Context, count int64) ([]*domain.User, error)

	Close() error
}
