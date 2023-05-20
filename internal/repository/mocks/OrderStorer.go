// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	repository "github.com/sagar23sj/go-ecommerce/internal/repository"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// OrderStorer is an autogenerated mock type for the OrderStorer type
type OrderStorer struct {
	mock.Mock
}

// BeginTx provides a mock function with given fields: ctx
func (_m *OrderStorer) BeginTx(ctx context.Context) (repository.Transaction, error) {
	ret := _m.Called(ctx)

	var r0 repository.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (repository.Transaction, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) repository.Transaction); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repository.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateOrder provides a mock function with given fields: ctx, tx, order
func (_m *OrderStorer) CreateOrder(ctx context.Context, tx repository.Transaction, order repository.Order) (repository.Order, error) {
	ret := _m.Called(ctx, tx, order)

	var r0 repository.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.Transaction, repository.Order) (repository.Order, error)); ok {
		return rf(ctx, tx, order)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.Transaction, repository.Order) repository.Order); ok {
		r0 = rf(ctx, tx, order)
	} else {
		r0 = ret.Get(0).(repository.Order)
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.Transaction, repository.Order) error); ok {
		r1 = rf(ctx, tx, order)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrderByID provides a mock function with given fields: ctx, tx, orderID
func (_m *OrderStorer) GetOrderByID(ctx context.Context, tx repository.Transaction, orderID int64) (repository.Order, error) {
	ret := _m.Called(ctx, tx, orderID)

	var r0 repository.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.Transaction, int64) (repository.Order, error)); ok {
		return rf(ctx, tx, orderID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.Transaction, int64) repository.Order); ok {
		r0 = rf(ctx, tx, orderID)
	} else {
		r0 = ret.Get(0).(repository.Order)
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.Transaction, int64) error); ok {
		r1 = rf(ctx, tx, orderID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HandleTransaction provides a mock function with given fields: ctx, tx, incomingErr
func (_m *OrderStorer) HandleTransaction(ctx context.Context, tx repository.Transaction, incomingErr error) error {
	ret := _m.Called(ctx, tx, incomingErr)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.Transaction, error) error); ok {
		r0 = rf(ctx, tx, incomingErr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListOrders provides a mock function with given fields: ctx, tx
func (_m *OrderStorer) ListOrders(ctx context.Context, tx repository.Transaction) ([]repository.Order, error) {
	ret := _m.Called(ctx, tx)

	var r0 []repository.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.Transaction) ([]repository.Order, error)); ok {
		return rf(ctx, tx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.Transaction) []repository.Order); ok {
		r0 = rf(ctx, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]repository.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.Transaction) error); ok {
		r1 = rf(ctx, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOrderDispatchDate provides a mock function with given fields: ctx, tx, orderID, dispatchedAt
func (_m *OrderStorer) UpdateOrderDispatchDate(ctx context.Context, tx repository.Transaction, orderID int64, dispatchedAt time.Time) error {
	ret := _m.Called(ctx, tx, orderID, dispatchedAt)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.Transaction, int64, time.Time) error); ok {
		r0 = rf(ctx, tx, orderID, dispatchedAt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateOrderStatus provides a mock function with given fields: ctx, tx, orderID, status
func (_m *OrderStorer) UpdateOrderStatus(ctx context.Context, tx repository.Transaction, orderID int64, status string) error {
	ret := _m.Called(ctx, tx, orderID, status)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.Transaction, int64, string) error); ok {
		r0 = rf(ctx, tx, orderID, status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewOrderStorer interface {
	mock.TestingT
	Cleanup(func())
}

// NewOrderStorer creates a new instance of OrderStorer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewOrderStorer(t mockConstructorTestingTNewOrderStorer) *OrderStorer {
	mock := &OrderStorer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
