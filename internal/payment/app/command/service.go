package command

import (
	"context"
	"github.com/Hana-bii/gorder-v2/common/genproto/orderpb"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, order *orderpb.Order) error
}
