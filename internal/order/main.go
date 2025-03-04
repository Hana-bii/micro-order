package main

import (
	"context"
	"github.com/Hana-bii/gorder-v2/common/tracing"

	"github.com/Hana-bii/gorder-v2/common/broker"
	"github.com/Hana-bii/gorder-v2/common/config"
	"github.com/Hana-bii/gorder-v2/common/discovery"
	"github.com/Hana-bii/gorder-v2/common/genproto/orderpb"
	"github.com/Hana-bii/gorder-v2/common/logging"
	"github.com/Hana-bii/gorder-v2/common/server"
	"github.com/Hana-bii/gorder-v2/order/infrastructure/consumer"
	"github.com/Hana-bii/gorder-v2/order/ports"
	"github.com/Hana-bii/gorder-v2/order/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	logging.Init()
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("order.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	// 主函数结束时close grpc conn
	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	// 注册Consul服务发现
	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()

	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.passwd"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	defer func() {
		_ = closeCh()
		_ = ch.Close()
	}()
	go consumer.NewConsumer(application).Listen(ch)

	// 丢入协程，防止阻塞HTTP服务
	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		svc := ports.NewGRPCServer(application)
		orderpb.RegisterOrderServiceServer(server, svc)
	})

	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		router.StaticFile("/success", "../../public/success.html")
		ports.RegisterHandlersWithOptions(router, HTTPServer{app: application}, ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})

}
