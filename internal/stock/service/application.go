package service

import (
	"context"

	"github.com/Hana-bii/gorder-v2/common/metrics"
	"github.com/Hana-bii/gorder-v2/stock/adapters"
	"github.com/Hana-bii/gorder-v2/stock/app"
	"github.com/Hana-bii/gorder-v2/stock/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(_ context.Context) app.Application {
	stockRepo := adapters.NewMemoryStockRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, logger, metricsClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricsClient),
		},
	}
}
