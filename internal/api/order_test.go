package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sagar23sj/go-ecommerce/internal/app/order/mocks"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/apperrors"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/dto"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type OrderAPITestSuite struct {
	suite.Suite
	orderSvc *mocks.Service
	router   chi.Router
}

func TestOrderAPITestSuite(t *testing.T) {
	suite.Run(t, new(OrderAPITestSuite))
}

// this function executes before the test suite begins execution
func (suite *OrderAPITestSuite) SetupTest() {
	suite.orderSvc = &mocks.Service{}
	suite.router = chi.NewRouter()
}

// this function executes after all tests executed
func (suite *OrderAPITestSuite) TearDownTest() {
	suite.orderSvc.AssertExpectations(suite.T())
}

func (suite *OrderAPITestSuite) TestGetOrderDetailsHandler() {
	t := suite.T()
	testCases := []struct {
		name               string
		orderID            interface{}
		setup              func()
		expectedStatusCode int
	}{
		{
			name:    "Success",
			orderID: 1,
			setup: func() {
				suite.orderSvc.On("GetOrderDetailsByID", mock.Anything, int64(1)).Return(dto.Order{
					ID:                 int64(1),
					Products:           []dto.ProductInfo{{ProductID: 1, Quantity: 2}},
					Amount:             20.0,
					DiscountPercentage: 0.0,
					FinalAmount:        20.0,
					Status:             "Placed",
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:    "Fail Because Order Not Found",
			orderID: 1,
			setup: func() {
				suite.orderSvc.On("GetOrderDetailsByID", mock.Anything, int64(1)).Return(dto.Order{}, apperrors.ProductNotFound{ID: 1})
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:    "Fail Because Invalid OrderID In Request",
			orderID: "w",
			setup: func() {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:    "Fail Because Something Went Wrong",
			orderID: 1,
			setup: func() {
				suite.orderSvc.On("GetOrderDetailsByID", mock.Anything, int64(1)).Return(dto.Order{}, errors.New("something went wrong"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range testCases {
		suite.SetupTest()
		suite.Run(test.name, func() {
			test.setup()

			suite.router.Get("/orders/{id}", getOrderDetailsHandler(suite.orderSvc))
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/orders/%v", test.orderID), bytes.NewBuffer([]byte(``)))
			if err != nil {
				t.Errorf("error occured while making http request, error : %v", err.Error())
			}

			recorder := httptest.NewRecorder()
			suite.router.ServeHTTP(recorder, req)

			suite.Equal(test.expectedStatusCode, recorder.Code)
		})
		suite.TearDownTest()
	}
}

func (suite *OrderAPITestSuite) TestListOrdersHandler() {
	t := suite.T()
	testCases := []struct {
		name               string
		setup              func()
		expectedStatusCode int
	}{
		{
			name: "Success",
			setup: func() {
				suite.orderSvc.On("ListOrders", mock.Anything).Return([]dto.Order{
					{
						ID:                 int64(1),
						Products:           []dto.ProductInfo{{ProductID: 1, Quantity: 2}},
						Amount:             20.0,
						DiscountPercentage: 0.0,
						FinalAmount:        20.0,
						Status:             "Placed",
					},
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Fail Because Something Went Wrong",
			setup: func() {
				suite.orderSvc.On("ListOrders", mock.Anything).Return([]dto.Order{}, errors.New("something went wrong"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range testCases {
		suite.SetupTest()
		suite.Run(test.name, func() {
			test.setup()

			suite.router.Get("/orders", listOrdersHandler(suite.orderSvc))
			req, err := http.NewRequest(http.MethodGet, "/orders", bytes.NewBuffer([]byte(``)))
			if err != nil {
				t.Errorf("error occured while making http request, error : %v", err.Error())
			}

			recorder := httptest.NewRecorder()
			suite.router.ServeHTTP(recorder, req)

			suite.Equal(test.expectedStatusCode, recorder.Code)
		})
		suite.TearDownTest()
	}
}

func (suite *OrderAPITestSuite) TestUpdateOrderStatusHandler() {
	t := suite.T()
	testCases := []struct {
		name               string
		input              dto.UpdateOrderStatusRequest
		setup              func()
		expectedStatusCode int
	}{
		{
			name: "Success",
			input: dto.UpdateOrderStatusRequest{
				OrderID: 1,
				Status:  "Dispatched",
			},
			setup: func() {
				suite.orderSvc.On("UpdateOrderStatus", mock.Anything, int64(1), "Dispatched").Return(dto.Order{
					ID:                 int64(1),
					Products:           []dto.ProductInfo{{ProductID: 1, Quantity: 2}},
					Amount:             20.0,
					DiscountPercentage: 0.0,
					FinalAmount:        20.0,
					Status:             "Dispatched",
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Fail Because Order Status Invalid",
			input: dto.UpdateOrderStatusRequest{
				OrderID: 1,
				Status:  "test",
			},
			setup: func() {
				suite.orderSvc.On("UpdateOrderStatus", mock.Anything, int64(1), "test").Return(dto.Order{}, apperrors.OrderStatusInvalid{ID: 1})
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "Fail Because Order Updation Not Allowed",
			input: dto.UpdateOrderStatusRequest{
				OrderID: 1,
				Status:  "Cancelled",
			},
			setup: func() {
				suite.orderSvc.On("UpdateOrderStatus", mock.Anything, int64(1), "Cancelled").Return(dto.Order{}, apperrors.OrderUpdationInvalid{
					ID:             1,
					RequestedState: "Cancelled",
					CurrentState:   "Completed",
				})
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "Fail Because Order Not Found",
			input: dto.UpdateOrderStatusRequest{
				OrderID: 1,
				Status:  "Cancelled",
			},
			setup: func() {
				suite.orderSvc.On("UpdateOrderStatus", mock.Anything, int64(1), "Cancelled").Return(dto.Order{}, apperrors.OrderNotFound{ID: 1})
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "Fail Because Something Went Wrong",
			input: dto.UpdateOrderStatusRequest{
				OrderID: 1,
				Status:  "Cancelled",
			},
			setup: func() {
				suite.orderSvc.On("UpdateOrderStatus", mock.Anything, int64(1), "Cancelled").Return(dto.Order{}, errors.New("something went wrong"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Fail Because OrderID Missing",
			input: dto.UpdateOrderStatusRequest{
				Status: "Cancelled",
			},
			setup: func() {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail Because Status Missing",
			input: dto.UpdateOrderStatusRequest{
				OrderID: 1,
			},
			setup: func() {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range testCases {
		suite.SetupTest()
		suite.Run(test.name, func() {
			test.setup()

			suite.router.Patch("/orders/{id}/status", updateOrderStatusHandler(suite.orderSvc))
			requestObj, err := json.Marshal(test.input)
			if err != nil {
				logger.Errorw(context.Background(), "error occured while marshaling json request")
			}

			req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/orders/%v/status", test.input.OrderID), bytes.NewBuffer(requestObj))
			if err != nil {
				t.Errorf("error occured while making http request, error : %v", err.Error())
			}

			recorder := httptest.NewRecorder()
			suite.router.ServeHTTP(recorder, req)

			suite.Equal(test.expectedStatusCode, recorder.Code)
		})
		suite.TearDownTest()
	}
}

func (suite *OrderAPITestSuite) TestCreateOrderHandler() {
	t := suite.T()
	testCases := []struct {
		name               string
		input              dto.CreateOrderRequest
		setup              func()
		expectedStatusCode int
	}{
		{
			name: "Success",
			input: dto.CreateOrderRequest{
				Products: []dto.ProductInfo{
					{
						ProductID: 1,
						Quantity:  2,
					},
					{
						ProductID: 2,
						Quantity:  2,
					},
				},
			},
			setup: func() {
				suite.orderSvc.On("CreateOrder", mock.Anything, dto.CreateOrderRequest{
					Products: []dto.ProductInfo{
						{
							ProductID: 1,
							Quantity:  2,
						},
						{
							ProductID: 2,
							Quantity:  2,
						},
					},
				}).Return(dto.Order{
					ID:                 int64(1),
					Products:           []dto.ProductInfo{{ProductID: 1, Quantity: 2}, {ProductID: 2, Quantity: 2}},
					Amount:             20.0,
					DiscountPercentage: 0.0,
					FinalAmount:        20.0,
					Status:             "Placed",
				}, nil)
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "Fail Because Missing Products In Order",
			input: dto.CreateOrderRequest{
				Products: []dto.ProductInfo{},
			},
			setup: func() {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail Because Duplicate Products In Request",
			input: dto.CreateOrderRequest{
				Products: []dto.ProductInfo{
					{
						ProductID: 1,
						Quantity:  2,
					},
					{
						ProductID: 1,
						Quantity:  2,
					},
				},
			},
			setup: func() {
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "Fail Because Negative Quantity Of Products",
			input: dto.CreateOrderRequest{
				Products: []dto.ProductInfo{
					{
						ProductID: 1,
						Quantity:  -2,
					},
					{
						ProductID: 2,
						Quantity:  2,
					},
				},
			},
			setup: func() {
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "Fail Because Product Not Found",
			input: dto.CreateOrderRequest{
				Products: []dto.ProductInfo{
					{
						ProductID: 1,
						Quantity:  2,
					},
					{
						ProductID: 2,
						Quantity:  2,
					},
				},
			},
			setup: func() {
				suite.orderSvc.On("CreateOrder", mock.Anything, dto.CreateOrderRequest{
					Products: []dto.ProductInfo{
						{
							ProductID: 1,
							Quantity:  2,
						},
						{
							ProductID: 2,
							Quantity:  2,
						},
					},
				}).Return(dto.Order{}, apperrors.ProductNotFound{ID: 1})
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "Fail Because Product Quantity Limit Exceeded",
			input: dto.CreateOrderRequest{
				Products: []dto.ProductInfo{
					{
						ProductID: 1,
						Quantity:  11,
					},
					{
						ProductID: 2,
						Quantity:  2,
					},
				},
			},
			setup: func() {
				suite.orderSvc.On("CreateOrder", mock.Anything, dto.CreateOrderRequest{
					Products: []dto.ProductInfo{
						{
							ProductID: 1,
							Quantity:  11,
						},
						{
							ProductID: 2,
							Quantity:  2,
						},
					},
				}).Return(dto.Order{}, apperrors.ProductQuantityExceeded{ID: 1, QuantityLimit: 10, QuantityAsked: 11})
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "Fail Because Product Quantity Insufficient",
			input: dto.CreateOrderRequest{
				Products: []dto.ProductInfo{
					{
						ProductID: 1,
						Quantity:  8,
					},
					{
						ProductID: 2,
						Quantity:  2,
					},
				},
			},
			setup: func() {
				suite.orderSvc.On("CreateOrder", mock.Anything, dto.CreateOrderRequest{
					Products: []dto.ProductInfo{
						{
							ProductID: 1,
							Quantity:  8,
						},
						{
							ProductID: 2,
							Quantity:  2,
						},
					},
				}).Return(dto.Order{}, apperrors.ProductQuantityInsufficient{ID: 1, QuantityRemaining: 4, QuantityAsked: 8})
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "Fail Because Something Went Wrong",
			input: dto.CreateOrderRequest{
				Products: []dto.ProductInfo{
					{
						ProductID: 1,
						Quantity:  8,
					},
					{
						ProductID: 2,
						Quantity:  2,
					},
				},
			},
			setup: func() {
				suite.orderSvc.On("CreateOrder", mock.Anything, dto.CreateOrderRequest{
					Products: []dto.ProductInfo{
						{
							ProductID: 1,
							Quantity:  8,
						},
						{
							ProductID: 2,
							Quantity:  2,
						},
					},
				}).Return(dto.Order{}, errors.New("something went wrong"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range testCases {
		suite.SetupTest()
		suite.Run(test.name, func() {
			test.setup()

			suite.router.Post("/orders", createOrderHandler(suite.orderSvc))
			requestObj, err := json.Marshal(test.input)
			if err != nil {
				logger.Errorw(context.Background(), "error occured while marshaling json request")
			}

			req, err := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(requestObj))
			if err != nil {
				t.Errorf("error occured while making http request, error : %v", err.Error())
			}

			recorder := httptest.NewRecorder()
			suite.router.ServeHTTP(recorder, req)

			suite.Equal(test.expectedStatusCode, recorder.Code)
		})
		suite.TearDownTest()
	}
}
