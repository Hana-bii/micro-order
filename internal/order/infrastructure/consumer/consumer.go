package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"

	"github.com/Hana-bii/gorder-v2/common/broker"
	"github.com/Hana-bii/gorder-v2/order/app"
	"github.com/Hana-bii/gorder-v2/order/app/command"
	domain "github.com/Hana-bii/gorder-v2/order/domain/order"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	app app.Application
}

func NewConsumer(application app.Application) *Consumer {
	return &Consumer{
		app: application,
	}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.EventOrderPaid, true, false, true, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	// bind到exchange
	err = ch.QueueBind(q.Name, "", broker.EventOrderPaid, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	// 用通道阻塞该函数
	var forever chan struct{}
	go func() {
		for msg := range msgs {
			c.handleMessage(msg, q)
		}
	}()
	<-forever
}

func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue) {
	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	t := otel.Tracer("rabbitmq")
	_, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
	defer span.End()

	// 接受到消息后更新状态
	o := &domain.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		logrus.Infof("error unmarshal mas.body into domain.order, err = #{err}")
		_ = msg.Nack(false, false)
		return

	}

	_, err := c.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
		Order: o,
		UpdateFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			if err := order.IsPaid(); err != nil {
				return nil, err
			}
			return order, nil
		},
	})
	if err != nil {
		logrus.Infof("error updating order, orderID = #{o.ID}, err err = #{err}")
		// TODO: retry
		return
	}
	span.AddEvent("order.updated")
	_ = msg.Ack(false)
	logrus.Infof("order consume paid event success!")
}
