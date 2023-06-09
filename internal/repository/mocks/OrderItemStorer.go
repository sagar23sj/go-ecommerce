// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	repository "github.com/sagar23sj/go-ecommerce/internal/repository"
	mock "github.com/stretchr/testify/mock"
)

// OrderItemStorer is an autogenerated mock type for the OrderItemStorer type
type OrderItemStorer struct {
	mock.Mock
}

// BeginTx provides a mock function with given fields: ctx
func (_m *OrderItemStorer) BeginTx(ctx context.Context) (repository.Transaction, error) {
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

// GetOrderItemsByOrderID provides a mock function with given fields: ctx, tx, orderID
func (_m *OrderItemStorer) GetOrderItemsByOrderID(ctx context.Context, tx repository.Transaction, orderID int64) ([]repository.OrderItem, error) {
	ret := _m.Called(ctx, tx, orderID)

	var r0 []repository.OrderItem
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.Transaction, int64) ([]repository.OrderItem, error)); ok {
		return rf(ctx, tx, orderID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.Transaction, int64) []repository.OrderItem); ok {
		r0 = rf(ctx, tx, orderID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]repository.OrderItem)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.Transaction, int64) error); ok {
		r1 = rf(ctx, tx, orderID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HandleTransaction provides a mock function with given fields: ctx, tx, incomingErr
func (_m *OrderItemStorer) HandleTransaction(ctx context.Context, tx repository.Transaction, incomingErr error) error {
	ret := _m.Called(ctx, tx, incomingErr)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.Transaction, error) error); ok {
		r0 = rf(ctx, tx, incomingErr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StoreOrderItems provides a mock function with given fields: ctx, tx, orderItems
func (_m *OrderItemStorer) StoreOrderItems(ctx context.Context, tx repository.Transaction, orderItems []repository.OrderItem) error {
	ret := _m.Called(ctx, tx, orderItems)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.Transaction, []repository.OrderItem) error); ok {
		r0 = rf(ctx, tx, orderItems)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewOrderItemStorer interface {
	mock.TestingT
	Cleanup(func())
}

// NewOrderItemStorer creates a new instance of OrderItemStorer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewOrderItemStorer(t mockConstructorTestingTNewOrderItemStorer) *OrderItemStorer {
	mock := &OrderItemStorer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
