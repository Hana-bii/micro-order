package adapters

import (
	"context"
	domain "github.com/Hana-bii/gorder-v2/order/domain/order"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

type MemoryOrderRepository struct {
	lock  *sync.RWMutex
	store []*domain.Order
}

func (m *MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*domain.Order, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	newOrder := &domain.Order{
		ID:          strconv.FormatInt(time.Now().UnixNano(), 10),
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
		Items:       order.Items,
	}
	m.store = append(m.store, newOrder)
	logrus.WithFields(logrus.Fields{
		"input_order":        order,
		"store_after_create": m.store,
	}).Info("memory_order_repo_create")
	return newOrder, nil
}

func (m *MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*domain.Order, error) {
	logrus.Infof("store: %v", m.store)
	for i, v := range m.store {
		logrus.Infof("m.store[%d]: %v", i, v)
	}
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, o := range m.store {
		if o.ID == id && customerID == o.CustomerID {
			logrus.Infof("memory_order_repo_get || found || id=%s || customerID=%s || res=%+v", id, customerID, *o)
			return o, nil
		}
	}
	return nil, domain.NotFoundError{OrderID: id}
}

func (m *MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	found := false
	defer func() {
		logrus.Infof("memory_order_repo || orderID=%s || found=%v", order.ID, found)
	}()
	for i, o := range m.store {
		if o.ID == order.ID && o.CustomerID == order.CustomerID {
			found = true
			// 找到之后去更新
			updatedOrder, err := updateFn(ctx, order)
			if err != nil {
				return err
			}
			m.store[i] = updatedOrder
		}
	}
	if !found {
		return domain.NotFoundError{OrderID: order.ID}
	}
	return nil
}

func NewMemoryOrderRepository() *MemoryOrderRepository {
	s := make([]*domain.Order, 0)
	s = append(s, &domain.Order{
		ID:          "fake-ID",
		CustomerID:  "fake-customer-id",
		Status:      "fake-status",
		PaymentLink: "fake-payment-link",
		Items:       nil,
	})
	return &MemoryOrderRepository{
		store: s,
		lock:  &sync.RWMutex{},
	}
}
