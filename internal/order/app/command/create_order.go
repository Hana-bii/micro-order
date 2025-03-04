package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Hana-bii/gorder-v2/common/broker"
	"github.com/Hana-bii/gorder-v2/common/decorator"
	"github.com/Hana-bii/gorder-v2/order/app/query"
	"github.com/Hana-bii/gorder-v2/order/convertor"
	domain "github.com/Hana-bii/gorder-v2/order/domain/order"
	"github.com/Hana-bii/gorder-v2/order/entity"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type CreateOrder struct {
	CustomerID string
	Items      []*entity.ItemWithQuantity
}

type CreateOrderResult struct {
	OrderID string
}

type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]

// 面向接口抽象
type createOrderHandler struct {
	orderRepo domain.Repository
	stockGRPC query.StockService
	channel   *amqp.Channel
}

func NewCreateOrderHandler(
	orderRepo domain.Repository,
	stockGRPC query.StockService,
	channel *amqp.Channel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) CreateOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}
	if stockGRPC == nil {
		panic("nil stockGRPC")
	}
	if channel == nil {
		panic("nil channel")
	}
	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
		createOrderHandler{
			orderRepo: orderRepo,
			stockGRPC: stockGRPC,
			channel:   channel,
		},
		logger,
		metricsClient,
	)
}

func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {

	// 注册队列
	// mq 为异步链路，无法带上span，不能使用上下文进行链路追踪
	q, err := c.channel.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	t := otel.Tracer("rabbitmq")
	ctx, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.publish", q.Name))
	defer span.End()
	validItems, err := c.validate(ctx, cmd.Items)
	if err != nil {
		return nil, err
	}
	o, err := c.orderRepo.Create(ctx, &domain.Order{
		CustomerID: cmd.CustomerID,
		Items:      validItems,
	})
	if err != nil {
		return nil, err
	}

	marshalledOrder, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	header := broker.InjectRabbitMQHeaders(ctx)

	// 在消息队列中发布事件
	err = c.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         marshalledOrder,
		Headers:      header,
	})
	if err != nil {
		return nil, err
	}

	return &CreateOrderResult{OrderID: o.ID}, nil
}

// 校验请求，合并同key-value
func (c createOrderHandler) validate(ctx context.Context, items []*entity.ItemWithQuantity) ([]*entity.Item, error) {
	if len(items) == 0 {
		return nil, errors.New("must have ar least one item")
	}
	// 合并数量
	items = packItems(items)
	// 检查库存, 这里调用STOCK的grpc，所以将item转一层，由内部实体结构转为orderpb标准结构
	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, convertor.NewItemWithQuantityConvertor().EntitiesToProtos(items))
	if err != nil {
		return nil, err
	}
	return convertor.NewItemConvertor().ProtosToEntities(resp.Items), nil
	//var ids []string
	//for _, item := range items {
	//	ids = append(ids, item.ID)
	//}
	//return c.stockGRPC.GetItems(ctx, ids)
}

func packItems(items []*entity.ItemWithQuantity) []*entity.ItemWithQuantity {
	merged := make(map[string]int32)
	for _, item := range items {
		merged[item.ID] += item.Quantity
	}
	var resp []*entity.ItemWithQuantity
	for id, quantity := range merged {
		resp = append(resp, &entity.ItemWithQuantity{
			ID:       id,
			Quantity: quantity,
		})
	}
	return resp
}
