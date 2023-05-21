package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sagar23sj/go-ecommerce/internal/app/order/mocks"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/apperrors"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/dto"
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

func (suite *OrderAPITestSuite) TestGetOrderHandler() {
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
