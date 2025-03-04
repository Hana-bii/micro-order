package query

import (
	"context"

	"github.com/Hana-bii/gorder-v2/common/decorator"
	"github.com/Hana-bii/gorder-v2/common/genproto/orderpb"
	domain "github.com/Hana-bii/gorder-v2/stock/domain/stock"
	"github.com/sirupsen/logrus"
)

type CheckIfItemsInStock struct {
	Items []*orderpb.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*orderpb.Item]

type checkIfItemsInStockHandler struct {
	stockRepo domain.Repository
}

func NewCheckIfItemsInStockHandler(
	stockRepo domain.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) CheckIfItemsInStockHandler {
	if stockRepo == nil {
		panic("nil stockRepo")
	}
	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*orderpb.Item](
		checkIfItemsInStockHandler{stockRepo: stockRepo},
		logger,
		metricsClient,
	)
}

var stub = map[string]string{
	"1": "price_1QwbhsG3A3fstR5HxhAZ6mGi",
	"2": "price_1QvdV8G3A3fstR5H7iJWbUUd",
	"3": "price_1QvdUXG3A3fstR5HqwBkCOwv",
}

// 具体查询方法
func (c checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*orderpb.Item, error) {
	var res []*orderpb.Item
	for _, item := range query.Items {
		// TODO: 改成从数据库 or Stripe 获取
		priceID, ok := stub[item.ID]
		if !ok {
			priceID = stub["1"]
		}
		res = append(res, &orderpb.Item{
			ID:       item.ID,
			Quantity: item.Quantity,
			PriceID:  priceID,
		})
	}
	return res, nil
}
