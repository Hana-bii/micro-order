package convertor

import (
	client "github.com/Hana-bii/gorder-v2/common/client/order"
	"github.com/Hana-bii/gorder-v2/common/genproto/orderpb"
	domain "github.com/Hana-bii/gorder-v2/order/domain/order"
	"github.com/Hana-bii/gorder-v2/order/entity"
)

type OrderConvertor struct {
}

func (c *OrderConvertor) EntityToProto(o *domain.Order) *orderpb.Order {
	c.Check(o)
	return &orderpb.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		Items:       NewItemConvertor().EntitiesToProtos(o.Items),
		PaymentLink: o.PaymentLink,
	}
}

func (c *OrderConvertor) ProtoToEntity(o *orderpb.Order) *domain.Order {
	c.Check(o)
	return &domain.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		Items:       NewItemConvertor().ProtosToEntities(o.Items),
		PaymentLink: o.PaymentLink,
	}
}

func (c *OrderConvertor) ClientToEntity(o *client.Order) *domain.Order {
	c.Check(o)
	return &domain.Order{
		ID:          o.Id,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
		Items:       NewItemConvertor().ClientsToEntities(o.Items),
	}
}

func (c *OrderConvertor) EntityToClient(o *domain.Order) *client.Order {
	c.Check(o)
	return &client.Order{
		Id:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
		Items:       NewItemConvertor().EntitiesToClients(o.Items),
	}
}

func (c *OrderConvertor) Check(o interface{}) {
	if o == nil {
		panic("cannot convert nil order")
	}
}

type ItemConvertor struct {
}

func (c *ItemConvertor) EntitiesToProtos(items []*entity.Item) (res []*orderpb.Item) {
	for _, item := range items {
		res = append(res, c.EntityToProto(item))
	}
	return
}

func (c *ItemConvertor) ProtosToEntities(items []*orderpb.Item) (res []*entity.Item) {
	for _, item := range items {
		res = append(res, c.ProtoToEntity(item))
	}
	return
}

func (c *ItemConvertor) ClientsToEntities(items []client.Item) (res []*entity.Item) {
	for _, item := range items {
		res = append(res, c.ClientToEntity(item))
	}
	return
}

func (c *ItemConvertor) EntitiesToClients(items []*entity.Item) (res []client.Item) {
	for _, item := range items {
		res = append(res, c.EntityToClient(item))
	}
	return
}

func (c *ItemConvertor) EntityToProto(item *entity.Item) *orderpb.Item {
	return &orderpb.Item{
		ID:       item.ID,
		Name:     item.Name,
		Quantity: item.Quantity,
		PriceID:  item.PriceID,
	}
}

func (c *ItemConvertor) ProtoToEntity(item *orderpb.Item) *entity.Item {
	return &entity.Item{
		ID:       item.ID,
		Name:     item.Name,
		Quantity: item.Quantity,
		PriceID:  item.PriceID,
	}
}

func (c *ItemConvertor) ClientToEntity(item client.Item) *entity.Item {
	return &entity.Item{
		ID:       item.Id,
		Name:     item.Name,
		Quantity: item.Quantity,
		PriceID:  item.PriceID,
	}
}

func (c *ItemConvertor) EntityToClient(item *entity.Item) client.Item {
	return client.Item{
		Id:       item.ID,
		Name:     item.Name,
		Quantity: item.Quantity,
		PriceID:  item.PriceID,
	}
}

type ItemWithQuantityConvertor struct {
}

func (c *ItemWithQuantityConvertor) EntitiesToProtos(items []*entity.ItemWithQuantity) (res []*orderpb.ItemWithQuantity) {
	for _, item := range items {
		res = append(res, c.EntityToProto(item))
	}
	return
}

func (c *ItemWithQuantityConvertor) EntityToProto(item *entity.ItemWithQuantity) *orderpb.ItemWithQuantity {
	return &orderpb.ItemWithQuantity{
		ID:       item.ID,
		Quantity: item.Quantity,
	}
}

func (c *ItemWithQuantityConvertor) ProtosToEntities(items []*orderpb.ItemWithQuantity) (res []*entity.ItemWithQuantity) {
	for _, item := range items {
		res = append(res, c.ProtoToEntity(item))
	}
	return
}

func (c *ItemWithQuantityConvertor) ProtoToEntity(item *orderpb.ItemWithQuantity) *entity.ItemWithQuantity {
	return &entity.ItemWithQuantity{
		ID:       item.ID,
		Quantity: item.Quantity,
	}
}

func (c *ItemWithQuantityConvertor) ClientsToEntities(items []client.ItemWithQuantity) (res []*entity.ItemWithQuantity) {
	for _, item := range items {
		res = append(res, c.ClientToEntity(item))
	}
	return
}

func (c *ItemWithQuantityConvertor) ClientToEntity(item client.ItemWithQuantity) *entity.ItemWithQuantity {
	return &entity.ItemWithQuantity{
		ID:       item.Id,
		Quantity: item.Quantity,
	}
}
