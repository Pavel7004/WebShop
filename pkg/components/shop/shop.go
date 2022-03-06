package shop

import (
	"context"
	"time"

	"github.com/Pavel7004/Common/tracing"
	dbi "github.com/Pavel7004/WebShop/pkg/adapters/db"
	"github.com/Pavel7004/WebShop/pkg/components"
	"github.com/Pavel7004/WebShop/pkg/domain"
)

type Shop struct {
	db dbi.DB
}

var _ components.Shop = (*Shop)(nil)

func New(db dbi.DB) *Shop {
	return &Shop{
		db: db,
	}
}

func (s *Shop) GetItemById(ctx context.Context, id string) (*domain.Item, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	span.SetTag("id", id)

	return s.db.GetItemById(ctx, id)
}

func (s *Shop) AddItem(ctx context.Context, item *domain.AddItemRequest) (string, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	span.SetTag("item_request", item)

	return s.db.AddItem(ctx, item)
}

func (s *Shop) GetItemsByPrice(ctx context.Context, from, to float64) ([]*domain.Item, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	span.SetTag("from", from)
	span.SetTag("to", to)

	return s.db.GetItemsByPrice(ctx, from, to)
}

func (s *Shop) GetRecentlyAddedItems(ctx context.Context, period time.Duration) ([]*domain.Item, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	span.SetTag("period", period.String())

	return s.db.GetRecentlyAddedItems(ctx, period)
}
