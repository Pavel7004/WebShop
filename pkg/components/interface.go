package components

import (
	"context"
	"time"

	"github.com/Pavel7004/WebShop/pkg/domain"
)

type Shop interface {
	AddItem(ctx context.Context, item *domain.AddItemRequest) (string, error)
	UpdateItem(ctx context.Context, id string, in *domain.UpdateItemRequest) (int64, error)
	GetItemById(ctx context.Context, id string) (*domain.Item, error)
	GetItemsByPrice(ctx context.Context, from, to float64) ([]*domain.Item, error)
	GetRecentlyAddedItems(ctx context.Context, period time.Duration) ([]*domain.Item, error)
	GetItemsByOwnerId(ctx context.Context, id string) ([]*domain.Item, error)

	RegisterUser(ctx context.Context, user *domain.RegisterUserRequest) (string, error)
	GetUserById(ctx context.Context, id string) (*domain.User, error)
	GetRecentlyAddedUsers(ctx context.Context, count int64) ([]*domain.User, error)

	CreateOrder(ctx context.Context, req *domain.CreateOrderRequest) (string, error)
	PayOrder(ctx context.Context, orderID string) error
	ProcessOrder(ctx context.Context, orderID string) error
}
