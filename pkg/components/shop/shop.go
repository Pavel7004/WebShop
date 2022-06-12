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

func (s *Shop) RegisterUser(ctx context.Context, user *domain.RegisterUserRequest) (string, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	span.SetTag("user_request", *user)

	return s.db.RegisterUser(ctx, user)
}

func (s *Shop) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	span.SetTag("id", id)

	return s.db.GetUserById(ctx, id)
}

func (s *Shop) GetItemsByOwnerId(ctx context.Context, id string) ([]*domain.Item, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	span.SetTag("id", id)

	return s.db.GetItemsByOwnerId(ctx, id)
}

func (s *Shop) GetRecentlyAddedUsers(ctx context.Context, count int64) ([]*domain.User, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	span.SetTag("count", count)

	return s.db.GetRecentlyAddedUsers(ctx, count)
}

func (s *Shop) UpdateItem(ctx context.Context, id string, in *domain.UpdateItemRequest) (int64, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	span.SetTag("id", id)

	return s.db.UpdateItem(ctx, id, in)
}

func (s *Shop) CreateOrder(ctx context.Context, req *domain.CreateOrderRequest) (string, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	return s.db.CreateOrder(ctx, req)
}

func (s *Shop) PayOrder(ctx context.Context, orderID string) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	modCount, err := s.db.UpdateOrder(ctx, orderID, domain.UpdateOrderRequest{
		Status: &domain.PAID,
	})
	if err != nil {
		return err
	}

	if modCount < 1 {
		return domain.ErrOrderNotProcessed
	}

	return nil
}

func (s *Shop) ProcessOrder(ctx context.Context, orderID string) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	return nil
}
