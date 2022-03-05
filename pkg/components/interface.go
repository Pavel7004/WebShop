package components

import (
	"context"

	"github.com/Pavel7004/WebShop/pkg/domain"
)

type Shop interface {
	GetItemById(ctx context.Context, id string) (*domain.Item, error)
	AddItem(ctx context.Context, item *domain.AddItemRequest) (string, error)
	GetItemsByPrice(ctx context.Context, from, to float64) ([]*domain.Item, error)
}
