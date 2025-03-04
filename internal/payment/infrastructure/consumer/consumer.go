package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"

	"github.com/Hana-bii/gorder-v2/common/broker"
	"github.com/Hana-bii/gorder-v2/common/genproto/orderpb"
	"github.com/Hana-bii/gorder-v2/payment/app"
	"github.com/Hana-bii/gorder-v2/payment/app/command"
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
	q, err := ch.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		// 在main函数中出问题可以用fatal
		logrus.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("fail to consume: queue=%s, err=%v", q.Name, err)

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
	logrus.Infof("Payment recieve a message from %s, msg=%v", q.Name, string(msg.Body))
	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	tr := otel.Tracer("rabbitmq")
	_, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
	defer span.End()

	o := &orderpb.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		logrus.Infof("fail to unmarshal msg to order, err=#{err}")
		_ = msg.Nack(false, false)
		return
	}
	if _, err := c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{Order: o}); err != nil {
		// TODO: retry
		logrus.Infof("fail to create order, err=#{err}")
		_ = msg.Nack(false, false)
		return

	}
	span.AddEvent("payment.created")
	_ = msg.Ack(false)
	logrus.Info("consume success")
}
