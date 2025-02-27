package service

import (
	"context"
	grpcClient "github.com/Hana-bii/gorder-v2/common/client"
	"github.com/Hana-bii/gorder-v2/common/metrics"
	"github.com/Hana-bii/gorder-v2/payment/adapters"
	"github.com/Hana-bii/gorder-v2/payment/app"
	"github.com/Hana-bii/gorder-v2/payment/app/command"
	"github.com/Hana-bii/gorder-v2/payment/domain"
	"github.com/Hana-bii/gorder-v2/payment/infrastructure/processor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	// 实际的存储基础设施在这里更换，不影响业务逻辑
	orderClient, closeOrderClient, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	orderGRPC := adapters.NewOrderGRPC(orderClient)
	// TODO: processor更换为stripe
	// memoryProcessor := processor.NewInmemProcessor()
	stripeProcessor := processor.NewStripeProcessor(viper.GetString("stripe-key"))
	return newApplication(ctx, orderGRPC, stripeProcessor), func() {
		_ = closeOrderClient()
	}

}

func newApplication(ctx context.Context, orderGRPC command.OrderService, processor domain.Processor) app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, logger, metricsClient),
		},
	}
}
