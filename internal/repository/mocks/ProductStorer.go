// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	repository "github.com/sagar23sj/go-ecommerce/internal/repository"
	mock "github.com/stretchr/testify/mock"
)

// ProductStorer is an autogenerated mock type for the ProductStorer type
type ProductStorer struct {
	mock.Mock
}

// BeginTx provides a mock function with given fields: ctx
func (_m *ProductStorer) BeginTx(ctx context.Context) (repository.Transaction, error) {
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

// GetProductByID provides a mock function with given fields: ctx, tx, productID
func (_m *ProductStorer) GetProductByID(ctx context.Context, tx repository.Transaction, productID int64) (repository.Product, error) {
	ret := _m.Called(ctx, tx, productID)

	var r0 repository.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.Transaction, int64) (repository.Product, error)); ok {
		return rf(ctx, tx, productID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.Transaction, int64) repository.Product); ok {
		r0 = rf(ctx, tx, productID)
	} else {
		r0 = ret.Get(0).(repository.Product)
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.Transaction, int64) error); ok {
		r1 = rf(ctx, tx, productID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HandleTransaction provides a mock function with given fields: ctx, tx, incomingErr
func (_m *ProductStorer) HandleTransaction(ctx context.Context, tx repository.Transaction, incomingErr error) error {
	ret := _m.Called(ctx, tx, incomingErr)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.Transaction, error) error); ok {
		r0 = rf(ctx, tx, incomingErr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListProducts provides a mock function with given fields: ctx, tx
func (_m *ProductStorer) ListProducts(ctx context.Context, tx repository.Transaction) ([]repository.Product, error) {
	ret := _m.Called(ctx, tx)

	var r0 []repository.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.Transaction) ([]repository.Product, error)); ok {
		return rf(ctx, tx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.Transaction) []repository.Product); ok {
		r0 = rf(ctx, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]repository.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.Transaction) error); ok {
		r1 = rf(ctx, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProductQuantity provides a mock function with given fields: ctx, tx, productsQuantityMap
func (_m *ProductStorer) UpdateProductQuantity(ctx context.Context, tx repository.Transaction, productsQuantityMap map[int64]int64) error {
	ret := _m.Called(ctx, tx, productsQuantityMap)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.Transaction, map[int64]int64) error); ok {
		r0 = rf(ctx, tx, productsQuantityMap)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewProductStorer interface {
	mock.TestingT
	Cleanup(func())
}

// NewProductStorer creates a new instance of ProductStorer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProductStorer(t mockConstructorTestingTNewProductStorer) *ProductStorer {
	mock := &ProductStorer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
