package order

import (
	"context"

	"github.com/Gabbu98/orders-api/model"
)

type OrderRepository interface {
	Insert(ctx context.Context, order model.Order) error
	FindByID(ctx context.Context, id uint64) (model.Order, error)
	DeleteByID(ctx context.Context, id uint64) error
	Update(ctx context.Context, order model.Order) error
	FindAll(ctx context.Context, page FindAllPage) (FindResult, error)
}