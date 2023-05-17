package order

import (
	"context"
	"testing"

	productMock "github.com/sagar23sj/go-ecommerce/internal/app/product/mocks"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/apperrors"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/dto"
	"github.com/sagar23sj/go-ecommerce/internal/repository"
	"github.com/sagar23sj/go-ecommerce/internal/repository/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
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
				tx := &gorm.DB{}
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
			name: "Fail Product Limit Exceeded",
			input: dto.CreateOrderRequest{
				Products: []dto.ProductInfo{
					{
						ProductID: int64(1),
						Quantity:  int64(12),
					},
				},
			},
			setup: func() {
				tx := &gorm.DB{}
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
			name: "Fail Product Quantity Insufficient",
			input: dto.CreateOrderRequest{
				Products: []dto.ProductInfo{
					{
						ProductID: int64(1),
						Quantity:  int64(12),
					},
				},
			},
			setup: func() {
				tx := &gorm.DB{}
				suite.orderRepo.On("BeginTx", mock.Anything).Return(tx, nil)
				suite.orderRepo.On("HandleTransaction", mock.Anything, tx, mock.Anything).Return(nil)
				suite.productService.On("GetProductByID", mock.Anything, tx, int64(1)).Return(dto.Product{
					ID:       int64(1),
					Name:     "xyz",
					Price:    10.0,
					Category: "Premium",
					Quantity: int64(10),
				}, nil)
			},
			expectedOutput: dto.Order{},
			expectedErr: apperrors.ProductQuantityInsufficient{
				ID:                1,
				QuantityRemaining: 10,
				QuantityAsked:     12,
			},
		},
	}

	for _, test := range testCases {
		suite.SetupTest()
		suite.Run(test.name, func() {
			test.setup()

			_, err := suite.service.CreateOrder(context.Background(), dto.CreateOrderRequest{
				Products: []dto.ProductInfo{{ProductID: test.input.Products[0].ProductID, Quantity: test.input.Products[0].Quantity}},
			})

			suite.Equal(test.expectedErr, err)
		})
		suite.TearDownTest()
	}
}
