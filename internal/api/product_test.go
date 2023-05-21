package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sagar23sj/go-ecommerce/internal/app/product/mocks"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/apperrors"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/dto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ProductAPITestSuite struct {
	suite.Suite
	productSvc *mocks.Service
	router     chi.Router
}

func TestProductAPITestSuite(t *testing.T) {
	suite.Run(t, new(ProductAPITestSuite))
}

// this function executes before the test suite begins execution
func (suite *ProductAPITestSuite) SetupTest() {
	suite.productSvc = &mocks.Service{}
	suite.router = chi.NewRouter()
}

// this function executes after all tests executed
func (suite *ProductAPITestSuite) TearDownTest() {
	suite.productSvc.AssertExpectations(suite.T())
}

func (suite *ProductAPITestSuite) TestGetProductHandler() {
	t := suite.T()
	testCases := []struct {
		name               string
		productID          interface{}
		setup              func()
		expectedStatusCode int
	}{
		{
			name:      "Success",
			productID: 1,
			setup: func() {
				suite.productSvc.On("GetProductByID", mock.Anything, mock.Anything, int64(1)).Return(dto.Product{
					ID:       1,
					Name:     "XYZ",
					Category: "Premium",
					Price:    100.0,
					Quantity: 10,
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:      "Fail Because Product Not Found",
			productID: 1,
			setup: func() {
				suite.productSvc.On("GetProductByID", mock.Anything, mock.Anything, int64(1)).Return(dto.Product{}, apperrors.ProductNotFound{ID: 1})
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:      "Fail Because Invalid ProductID In Request",
			productID: "w",
			setup: func() {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range testCases {
		suite.SetupTest()
		suite.Run(test.name, func() {
			test.setup()

			suite.router.Get("/products/{id}", getProductHandler(suite.productSvc))
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/products/%v", test.productID), bytes.NewBuffer([]byte(``)))
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

func (suite *ProductAPITestSuite) TestListProductsHandler() {
	t := suite.T()
	testCases := []struct {
		name               string
		setup              func()
		expectedStatusCode int
	}{
		{
			name: "Success",
			setup: func() {
				suite.productSvc.On("ListProducts", mock.Anything).Return([]dto.Product{
					{ID: 1,
						Name:     "XYZ",
						Category: "Premium",
						Price:    100.0,
						Quantity: 10,
					},
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Fail Because Something Went Wrong",
			setup: func() {
				suite.productSvc.On("ListProducts", mock.Anything).Return([]dto.Product{}, errors.New("something went wrong"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range testCases {
		suite.SetupTest()
		suite.Run(test.name, func() {
			test.setup()

			suite.router.Get("/products", listProductHandler(suite.productSvc))
			req, err := http.NewRequest(http.MethodGet, "/products", bytes.NewBuffer([]byte(``)))
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
