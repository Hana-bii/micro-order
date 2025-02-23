package service

import (
	"context"
	grpcClient "github.com/Hana-bii/gorder-v2/common/client"
	"github.com/Hana-bii/gorder-v2/common/metrics"
	"github.com/Hana-bii/gorder-v2/order/adapters"
	"github.com/Hana-bii/gorder-v2/order/adapters/grpc"
	"github.com/Hana-bii/gorder-v2/order/app"
	"github.com/Hana-bii/gorder-v2/order/app/command"
	"github.com/Hana-bii/gorder-v2/order/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	// 实际的存储基础设施在这里更换，不影响业务逻辑

	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	stockGRPC := grpc.NewStockGRPC(stockClient)
	return newApplication(ctx, stockGRPC), func() {
		_ = closeStockClient()
	}

}

func newApplication(_ context.Context, stockGRPC query.StockService) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}

	// 注入服务
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, logger, metricsClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricsClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricsClient),
		},
	}
}
