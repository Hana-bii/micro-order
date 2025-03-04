package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"
	"io"
	"net/http"

	"github.com/Hana-bii/gorder-v2/common/broker"
	"github.com/Hana-bii/gorder-v2/common/genproto/orderpb"
	"github.com/Hana-bii/gorder-v2/payment/domain"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
)

type PaymentHandler struct {
	channel *amqp.Channel
}

func NewPaymentHandler(ch *amqp.Channel) *PaymentHandler {
	return &PaymentHandler{channel: ch}

}

func (p *PaymentHandler) RegisterRoutes(c *gin.Engine) {
	c.POST("/api/webhook", p.handleWebhook)
}

func (p *PaymentHandler) handleWebhook(c *gin.Context) {
	// webhook示例代码
	logrus.Info("Got webhook from stripe")
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Infof("Error reading request body: %v\n", err)
		c.JSON(http.StatusServiceUnavailable, err.Error())
		return
	}

	event, err := webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"),
		viper.GetString("ENDPOINT_STRIPE_SECRET"))

	if err != nil {
		logrus.Infof("Error verifying webhook signature: %v\n", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			logrus.Infof("Error unmarshaling event.data.raw in session, err = %v", err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		// 支付成功，发到mq中
		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
			logrus.Infof("payment for checkout session %v success!", session.ID)

			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()

			var items []*orderpb.Item
			_ = json.Unmarshal([]byte(session.Metadata["items"]), &items)

			marshalledOrder, err := json.Marshal(&domain.Order{
				ID:          session.Metadata["orderID"],
				CustomerID:  session.Metadata["customerID"],
				Status:      string(stripe.CheckoutSessionPaymentStatusPaid),
				PaymentLink: session.Metadata["paymentLink"],
				Items:       items,
			})
			if err != nil {
				logrus.Infof("Error marshalling domain.order, err = %v", err)
				return
			}

			tr := otel.Tracer("rabbitmq")
			mqCtx, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.publish", broker.EventOrderPaid))
			defer span.End()

			headers := broker.InjectRabbitMQHeaders(mqCtx)
			_ = p.channel.PublishWithContext(mqCtx, broker.EventOrderPaid, "", false, false, amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Persistent,
				Body:         marshalledOrder,
				Headers:      headers,
			})
			logrus.Infof("message pulished to %s, body: %s", broker.EventOrderPaid, string(marshalledOrder))
		}
	}
	c.JSON(http.StatusOK, nil)
}
