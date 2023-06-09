package order

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/asdine/storm/v3"
	productMock "github.com/sagar23sj/go-ecommerce/internal/app/product/mocks"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/apperrors"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/dto"
	"github.com/sagar23sj/go-ecommerce/internal/repository"
	"github.com/sagar23sj/go-ecommerce/internal/repository/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type OrderServiceTestSuite struct {
	suite.Suite
	service        Service
	orderRepo      *mocks.OrderStorer
	orderItemRepo  *mocks.OrderItemStorer
	productService *productMock.Service
}

func TestOrderServiceTestSuite(t *testing.T) {
	suite.Run(t, new(OrderServiceTestSuite))
}

// this function executes before the test suite begins execution
func (suite *OrderServiceTestSuite) SetupTest() {
	suite.orderRepo = &mocks.OrderStorer{}
	suite.orderItemRepo = &mocks.OrderItemStorer{}
	suite.productService = &productMock.Service{}

	suite.service = NewService(suite.orderRepo, suite.orderItemRepo, suite.productService)
}

// this function executes after all tests executed
func (suite *OrderServiceTestSuite) TearDownTest() {
	suite.orderRepo.AssertExpectations(suite.T())
	suite.orderItemRepo.AssertExpectations(suite.T())
	suite.productService.AssertExpectations(suite.T())
}

func (suite *OrderServiceTestSuite) TestCreateOrder() {
	type testCaseStruct struct {
		name           string
		input          dto.CreateOrderRequest
		setup          func()
		expectedOutput dto.Order
		expectedErr    error
	}

	testCases := []testCaseStruct{
		{
			name: "Success",
			input: dto.CreateOrderRequest{
				Products: []dto.ProductInfo{
					{
						ProductID: int64(1),
						Quantity:  int64(2),
					},
				},
			},
			setup: func() {
				tx := &storm.DB{}
				suite.orderRepo.On("BeginTx", mock.Anything).Return(tx, nil)
				suite.orderRepo.On("HandleTransaction", mock.Anything, tx, nil).Return(nil)
				suite.productService.On("GetProductByID", mock.Anything, tx, int64(1)).Return(dto.Product{
					ID:       int64(1),
					Name:     "xyz",
					Price:    10.0,
					Category: "Premium",
					Quantity: int64(10),
				}, nil)
				suite.orderRepo.On("CreateOrder", mock.Anything, tx, repository.Order{
					Amount:             20.0,
					DiscountPercentage: 0.0,
					FinalAmount:        20.0,
					Status:             "Placed",
				}).Return(repository.Order{
					ID:                 uint(1),
					Amount:             20.0,
					DiscountPercentage: 0.0,
					FinalAmount:        20.0,
					Status:             "Placed",
				}, nil)
				suite.orderItemRepo.On("StoreOrderItems", mock.Anything, tx, []repository.OrderItem{{
					OrderID:   int64(1),
					ProductID: int64(1),
					Quantity:  int64(2),
				}}).Return(nil)
				suite.productService.On("UpdateProductQuantity", mock.Anything, tx, map[int64]int64{1: 8}).Return(nil)
			},
			expectedOutput: dto.Order{
				ID:                 int64(1),
				Products:           []dto.ProductInfo{{ProductID: 1, Quantity: 2}},
				Amount:             20.0,
				DiscountPercentage: 0.0,
				FinalAmount:        20.0,
				Status:             "Placed",
			},
			expectedErr: nil,
		},
		{
			name: "Success for 3 Premium Products",
			input: dto.CreateOrderRequest{
				Products: []dto.ProductInfo{
					{
						ProductID: int64(1),
						Quantity:  int64(2),
					},
					{
						ProductID: int64(2),
						Quantity:  int64(2),
					},
					{
						ProductID: int64(3),
						Quantity:  int64(2),
					},
				},
			},
			setup: func() {
				tx := &storm.DB{}
				suite.orderRepo.On("BeginTx", mock.Anything).Return(tx, nil)
				suite.orderRepo.On("HandleTransaction", mock.Anything, tx, nil).Return(nil)
				suite.productService.On("GetProductByID", mock.Anything, tx, int64(1)).Return(dto.Product{
					ID:       int64(1),
					Name:     "xyz",
					Price:    10.0,
					Category: "Premium",
					Quantity: int64(10),
				}, nil).Once()
				suite.productService.On("GetProductByID", mock.Anything, tx, int64(2)).Return(dto.Product{
					ID:       int64(2),
					Name:     "xyz",
					Price:    20.0,
					Category: "Premium",
					Quantity: int64(10),
				}, nil).Once()
				suite.productService.On("GetProductByID", mock.Anything, tx, int64(3)).Return(dto.Product{
					ID:       int64(3),
					Name:     "xyz",
					Price:    30.0,
					Category: "Premium",
					Quantity: int64(10),
				}, nil).Once()
				suite.orderRepo.On("CreateOrder", mock.Anything, tx, repository.Order{
					Amount:             120.0,
					DiscountPercentage: 10.0,
					FinalAmount:        108.0,
					Status:             "Placed",
				}).Return(repository.Order{
					ID:                 uint(1),
					Amount:             120.0,
					DiscountPercentage: 10.0,
					FinalAmount:        108.0,
					Status:             "Placed",
				}, nil)
				suite.orderItemRepo.On("StoreOrderItems", mock.Anything, tx, []repository.OrderItem{
					{
						OrderID:   int64(1),
						ProductID: int64(1),
						Quantity:  int64(2),
					},
					{
						OrderID:   int64(1),
						ProductID: int64(2),
						Quantity:  int64(2),
					},
					{
						OrderID:   int64(1),
						ProductID: int64(3),
						Quantity:  int64(2),
					},
				}).Return(nil)
				suite.productService.On("UpdateProductQuantity", mock.Anything, tx, map[int64]int64{1: 8, 2: 8, 3: 8}).Return(nil)
			},
			expectedOutput: dto.Order{
				ID:                 int64(1),
				Products:           []dto.ProductInfo{{ProductID: 1, Quantity: 2}},
				Amount:             120.0,
				DiscountPercentage: 10.0,
				FinalAmount:        108.0,
				Status:             "Placed",
			},
			expectedErr: nil,
		},
		{
			name: "Failed Because Product Limit Exceeded",
			input: dto.CreateOrderRequest{
				Products: []dto.ProductInfo{
					{
						ProductID: int64(1),
						Quantity:  int64(12),
					},
				},
			},
			setup: func() {
				tx := &storm.DB{}
				suite.orderRepo.On("BeginTx", mock.Anything).Return(tx, nil)
				suite.orderRepo.On("HandleTransaction", mock.Anything, tx, mock.Anything).Return(nil)
				suite.productService.On("GetProductByID", mock.Anything, tx, int64(1)).Return(dto.Product{
					ID:       int64(1),
					Name:     "xyz",
					Price:    10.0,
					Category: "Premium",
					Quantity: int64(20),
				}, nil)
			},
			expectedOutput: dto.Order{},
			expectedErr: apperrors.ProductQuantityExceeded{
				ID:            1,
				QuantityLimit: 10,
				QuantityAsked: 12,
			},
		},
		{
			name: "Fail Because Product Quantity Insufficient",
			input: dto.CreateOrderRequest{
				Products: []dto.ProductInfo{
					{
						ProductID: int64(1),
						Quantity:  int64(8),
					},
				},
			},
			setup: func() {
				tx := &storm.DB{}
				suite.orderRepo.On("BeginTx", mock.Anything).Return(tx, nil)
				suite.orderRepo.On("HandleTransaction", mock.Anything, tx, mock.Anything).Return(nil)
				suite.productService.On("GetProductByID", mock.Anything, tx, int64(1)).Return(dto.Product{
					ID:       int64(1),
					Name:     "xyz",
					Price:    10.0,
					Category: "Premium",
					Quantity: int64(6),
				}, nil)
			},
			expectedOutput: dto.Order{},
			expectedErr: apperrors.ProductQuantityInsufficient{
				ID:                1,
				QuantityRemaining: 6,
				QuantityAsked:     8,
			},
		},
	}

	for _, test := range testCases {
		suite.SetupTest()
		suite.Run(test.name, func() {
			test.setup()

			order, err := suite.service.CreateOrder(context.Background(), test.input)

			suite.Equal(test.expectedErr, err)
			suite.Equal(test.expectedOutput.FinalAmount, order.FinalAmount)
			suite.Equal(test.expectedOutput.Status, order.Status)
			suite.Equal(test.expectedOutput.DiscountPercentage, order.DiscountPercentage)
		})
		suite.TearDownTest()
	}
}

func (suite *OrderServiceTestSuite) TestUpdateOrderStatus() {
	type testCaseStruct struct {
		name           string
		input          dto.UpdateOrderStatusRequest
		setup          func()
		expectedOutput dto.Order
		expectedErr    error
	}

	now = func() time.Time { return time.Date(2023, 05, 18, 00, 00, 00, 00, time.UTC) }
	timeNow := now()
	testCases := []testCaseStruct{
		{
			name: "Success",
			input: dto.UpdateOrderStatusRequest{
				OrderID: 1,
				Status:  "Dispatched",
			},
			setup: func() {
				tx := &storm.DB{}
				suite.orderRepo.On("BeginTx", mock.Anything).Return(tx, nil)
				suite.orderRepo.On("HandleTransaction", mock.Anything, tx, mock.Anything).Return(nil)
				suite.orderRepo.On("GetOrderByID", mock.Anything, mock.Anything, int64(1)).Return(repository.Order{
					ID:                 uint(1),
					Amount:             20.0,
					DiscountPercentage: 0.0,
					FinalAmount:        20.0,
					Status:             "Placed",
				}, nil).Once()
				suite.orderRepo.On("UpdateOrderStatus", mock.Anything, mock.Anything, int64(1), "Dispatched").Return(nil)
				suite.orderRepo.On("UpdateOrderDispatchDate", mock.Anything, mock.Anything, int64(1), timeNow).Return(nil)
				suite.orderRepo.On("GetOrderByID", mock.Anything, mock.Anything, int64(1)).Return(repository.Order{
					ID:                 uint(1),
					Amount:             20.0,
					DiscountPercentage: 0.0,
					FinalAmount:        20.0,
					Status:             "Dispatched",
					DispatchedAt:       timeNow,
				}, nil).NotBefore()
			},
			expectedOutput: dto.Order{
				ID:                 int64(1),
				Products:           []dto.ProductInfo{},
				Amount:             20.0,
				DiscountPercentage: 0.0,
				FinalAmount:        20.0,
				Status:             "Dispatched",
			},
			expectedErr: nil,
		},
		{
			name: "Success When Order Cancelled",
			input: dto.UpdateOrderStatusRequest{
				OrderID: 1,
				Status:  "Cancelled",
			},
			setup: func() {
				tx := &storm.DB{}
				suite.orderRepo.On("BeginTx", mock.Anything).Return(tx, nil)
				suite.orderRepo.On("HandleTransaction", mock.Anything, tx, mock.Anything).Return(nil)
				suite.orderRepo.On("GetOrderByID", mock.Anything, mock.Anything, int64(1)).Return(repository.Order{
					ID:                 uint(1),
					Amount:             20.0,
					DiscountPercentage: 0.0,
					FinalAmount:        20.0,
					Status:             "Placed",
				}, nil).Once()
				suite.orderRepo.On("UpdateOrderStatus", mock.Anything, mock.Anything, int64(1), "Cancelled").Return(nil)
				suite.orderItemRepo.On("GetOrderItemsByOrderID", mock.Anything, mock.Anything, int64(1)).Return([]repository.OrderItem{
					{
						ID:        uint(1),
						OrderID:   1,
						ProductID: 1,
						Quantity:  2,
					},
				}, nil).Once()
				suite.productService.On("GetProductByID", mock.Anything, tx, int64(1)).Return(dto.Product{
					ID:       int64(1),
					Name:     "xyz",
					Price:    10.0,
					Category: "Premium",
					Quantity: int64(10),
				}, nil)
				suite.productService.On("UpdateProductQuantity", mock.Anything, tx, map[int64]int64{1: 12}).Return(nil)
				suite.orderRepo.On("GetOrderByID", mock.Anything, mock.Anything, int64(1)).Return(repository.Order{
					ID:                 uint(1),
					Amount:             20.0,
					DiscountPercentage: 0.0,
					FinalAmount:        20.0,
					Status:             "Cancelled",
					DispatchedAt:       timeNow,
				}, nil).NotBefore()
			},
			expectedOutput: dto.Order{
				ID:                 int64(1),
				Products:           []dto.ProductInfo{},
				Amount:             20.0,
				DiscountPercentage: 0.0,
				FinalAmount:        20.0,
				Status:             "Cancelled",
			},
			expectedErr: nil,
		},
		{
			name: "Failed Because Order Status Invalid ",
			input: dto.UpdateOrderStatusRequest{
				OrderID: int64(1),
				Status:  "test",
			},
			setup: func() {
				tx := &storm.DB{}
				suite.orderRepo.On("BeginTx", mock.Anything).Return(tx, nil)
				suite.orderRepo.On("HandleTransaction", mock.Anything, tx, mock.Anything).Return(nil)
			},
			expectedOutput: dto.Order{},
			expectedErr:    apperrors.OrderStatusInvalid{ID: int64(1)},
		},
		{
			name: "Failed Because Order Updation Invalid ",
			input: dto.UpdateOrderStatusRequest{
				OrderID: int64(1),
				Status:  "Completed",
			},
			setup: func() {
				tx := &storm.DB{}
				suite.orderRepo.On("BeginTx", mock.Anything).Return(tx, nil)
				suite.orderRepo.On("HandleTransaction", mock.Anything, tx, mock.Anything).Return(nil)
				suite.orderRepo.On("GetOrderByID", mock.Anything, mock.Anything, int64(1)).Return(repository.Order{
					ID:                 uint(1),
					Amount:             20.0,
					DiscountPercentage: 0.0,
					FinalAmount:        20.0,
					Status:             "Placed",
				}, nil)
			},
			expectedOutput: dto.Order{},
			expectedErr: apperrors.OrderUpdationInvalid{
				ID:             int64(1),
				CurrentState:   "Placed",
				RequestedState: "Completed",
			},
		},
		{
			name: "Failed Because Order Not Found ",
			input: dto.UpdateOrderStatusRequest{
				OrderID: int64(1),
				Status:  "Completed",
			},
			setup: func() {
				tx := &storm.DB{}
				suite.orderRepo.On("BeginTx", mock.Anything).Return(tx, nil)
				suite.orderRepo.On("HandleTransaction", mock.Anything, tx, mock.Anything).Return(nil)
				suite.orderRepo.On("GetOrderByID", mock.Anything, mock.Anything, int64(1)).Return(repository.Order{}, nil)
			},
			expectedOutput: dto.Order{},
			expectedErr:    apperrors.OrderNotFound{ID: int64(1)},
		},
		{
			name: "Failed Because Order Updation Failed",
			input: dto.UpdateOrderStatusRequest{
				OrderID: 1,
				Status:  "Dispatched",
			},
			setup: func() {
				tx := &storm.DB{}
				suite.orderRepo.On("BeginTx", mock.Anything).Return(tx, nil)
				suite.orderRepo.On("HandleTransaction", mock.Anything, tx, mock.Anything).Return(nil)
				suite.orderRepo.On("GetOrderByID", mock.Anything, mock.Anything, int64(1)).Return(repository.Order{
					ID:                 uint(1),
					Amount:             20.0,
					DiscountPercentage: 0.0,
					FinalAmount:        20.0,
					Status:             "Placed",
				}, nil).Once()
				suite.orderRepo.On("UpdateOrderStatus", mock.Anything, mock.Anything, int64(1), "Dispatched").Return(errors.New("something went wrong"))
			},
			expectedOutput: dto.Order{},
			expectedErr:    fmt.Errorf("error occured while updating order status: %w", errors.New("something went wrong")),
		},
	}

	for _, test := range testCases {
		suite.SetupTest()
		suite.Run(test.name, func() {
			test.setup()

			order, err := suite.service.UpdateOrderStatus(context.Background(), test.input.OrderID, test.input.Status)
			suite.Equal(test.expectedErr, err)
			suite.Equal(test.expectedOutput.Status, order.Status)
		})
		suite.TearDownTest()
	}
}

func (suite *OrderServiceTestSuite) TestGetOrdeDetails() {
	type testCaseStruct struct {
		name           string
		orderID        int64
		setup          func()
		expectedOutput dto.Order
		expectedErr    error
	}

	testCases := []testCaseStruct{
		{
			name:    "Success",
			orderID: int64(1),
			setup: func() {
				suite.orderRepo.On("GetOrderByID", mock.Anything, mock.Anything, int64(1)).Return(repository.Order{
					ID:                 uint(1),
					Amount:             20.0,
					DiscountPercentage: 0.0,
					FinalAmount:        20.0,
					Status:             "Placed",
				}, nil).Once()
				suite.orderItemRepo.On("GetOrderItemsByOrderID", mock.Anything, mock.Anything, int64(1)).Return([]repository.OrderItem{
					{
						ID:        uint(1),
						OrderID:   1,
						ProductID: 1,
						Quantity:  2,
					},
				}, nil).Once()
			},
			expectedOutput: dto.Order{
				ID:                 int64(1),
				Products:           []dto.ProductInfo{{ProductID: 1, Quantity: 2}},
				Amount:             20.0,
				DiscountPercentage: 0.0,
				FinalAmount:        20.0,
				Status:             "Placed",
			},
			expectedErr: nil,
		},
		{
			name:    "Fail Because Something Wrong With Fetching OrderItems",
			orderID: int64(1),
			setup: func() {
				suite.orderRepo.On("GetOrderByID", mock.Anything, mock.Anything, int64(1)).Return(repository.Order{
					ID:                 uint(1),
					Amount:             20.0,
					DiscountPercentage: 0.0,
					FinalAmount:        20.0,
					Status:             "Placed",
				}, nil).Once()
				suite.orderItemRepo.On("GetOrderItemsByOrderID", mock.Anything, mock.Anything, int64(1)).Return([]repository.OrderItem{}, errors.New("error fetching data for OrderItems")).Once()
			},
			expectedOutput: dto.Order{},
			expectedErr:    errors.New("error fetching data for OrderItems"),
		},
		{
			name:    "Fail Because Something Wrong With Fetching Order",
			orderID: int64(1),
			setup: func() {
				suite.orderRepo.On("GetOrderByID", mock.Anything, mock.Anything, int64(1)).Return(repository.Order{
					ID:                 uint(1),
					Amount:             20.0,
					DiscountPercentage: 0.0,
					FinalAmount:        20.0,
					Status:             "Placed",
				}, errors.New("error fetching data for Order")).Once()
			},
			expectedOutput: dto.Order{},
			expectedErr:    errors.New("error fetching data for Order"),
		},
		{
			name:    "Fail Because Order Not Found",
			orderID: int64(1),
			setup: func() {
				suite.orderRepo.On("GetOrderByID", mock.Anything, mock.Anything, int64(1)).Return(repository.Order{}, nil).Once()
			},
			expectedOutput: dto.Order{},
			expectedErr:    apperrors.OrderNotFound{ID: 1},
		},
	}

	for _, test := range testCases {
		suite.SetupTest()
		suite.Run(test.name, func() {
			test.setup()

			order, err := suite.service.GetOrderDetailsByID(context.Background(), test.orderID)
			suite.Equal(test.expectedErr, err)
			suite.Equal(test.expectedOutput.ID, order.ID)
			suite.Equal(test.expectedOutput.Products, order.Products)
			suite.Equal(test.expectedOutput.Amount, order.Amount)
		})
		suite.TearDownTest()
	}
}

func (suite *OrderServiceTestSuite) TestListOrders() {
	type testCaseStruct struct {
		name           string
		setup          func()
		expectedOutput []dto.Order
		expectedErr    error
	}

	testCases := []testCaseStruct{
		{
			name: "Success",
			setup: func() {
				suite.orderRepo.On("ListOrders", mock.Anything, mock.Anything).Return([]repository.Order{
					{
						ID:                 uint(1),
						Amount:             20.0,
						DiscountPercentage: 0.0,
						FinalAmount:        20.0,
						Status:             "Placed",
					},
				}, nil).Once()
			},
			expectedOutput: []dto.Order{
				{
					ID:                 int64(1),
					Products:           []dto.ProductInfo{},
					Amount:             20.0,
					DiscountPercentage: 0.0,
					FinalAmount:        20.0,
					Status:             "Cancelled",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Fail Because Something Wrong With Fetching Orders List",
			setup: func() {
				suite.orderRepo.On("ListOrders", mock.Anything, mock.Anything).Return([]repository.Order{}, errors.New("error fetching data for Orders")).Once()
			},
			expectedOutput: []dto.Order{},
			expectedErr:    errors.New("error fetching data for Orders"),
		},
	}

	for _, test := range testCases {
		suite.SetupTest()
		suite.Run(test.name, func() {
			test.setup()

			orderList, err := suite.service.ListOrders(context.Background())
			suite.Equal(test.expectedErr, err)
			suite.Equal(len(test.expectedOutput), len(orderList))
		})
		suite.TearDownTest()
	}
}
