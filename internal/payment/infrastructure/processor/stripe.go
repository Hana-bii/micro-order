package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hana-bii/gorder-v2/common/tracing"

	"github.com/Hana-bii/gorder-v2/common/genproto/orderpb"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

const (
	successURL = "http://localhost:8282/success"
)

type StripeProcessor struct {
	apiKey string
}

func (s StripeProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
	_, span := tracing.Start(ctx, "stripe_processor.create_payment_link")
	defer span.End()

	var items []*stripe.CheckoutSessionLineItemParams
	for _, item := range order.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			// price_ID从order Service向Stock GRPC Service请求查询库存的时候获取，
			// 并由Stock GRPC Service写入持久性存储库中
			Price:    stripe.String(item.PriceID),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}
	marshalledItems, _ := json.Marshal(order.Items)
	metadata := map[string]string{
		"orderID":     order.ID,
		"customerID":  order.CustomerID,
		"status":      order.Status,
		"items":       string(marshalledItems),
		"paymentLink": order.PaymentLink,
	}
	params := &stripe.CheckoutSessionParams{
		Metadata:   metadata,
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(fmt.Sprintf("%s?customerID=%s&orderID=%s", successURL, order.CustomerID, order.ID)),
	}
	result, err := session.New(params)
	if err != nil {
		return "", err
	}
	return result.URL, nil
}

func NewStripeProcessor(apiKey string) *StripeProcessor {
	if apiKey == "" {
		panic("empty api key")
	}
	// Key为全局，只赋值一次
	stripe.Key = apiKey
	return &StripeProcessor{apiKey: apiKey}

}
