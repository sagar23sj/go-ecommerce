package product

import (
	"context"
	"errors"
	"testing"

	"github.com/sagar23sj/go-ecommerce/internal/pkg/apperrors"
	"github.com/sagar23sj/go-ecommerce/internal/pkg/dto"
	"github.com/sagar23sj/go-ecommerce/internal/repository"
	"github.com/sagar23sj/go-ecommerce/internal/repository/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ProductServiceTestSuite struct {
	suite.Suite
	service     Service
	productRepo *mocks.ProductStorer
}

func TestProductServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceTestSuite))
}

// this function executes before the test suite begins execution
func (suite *ProductServiceTestSuite) SetupTest() {
	suite.productRepo = &mocks.ProductStorer{}

	suite.service = NewService(suite.productRepo)
}

// this function executes after all tests executed
func (suite *ProductServiceTestSuite) TearDownTest() {
	suite.productRepo.AssertExpectations(suite.T())
}

func (suite *ProductServiceTestSuite) TestGetProductByID() {

	testCases := []struct {
		name           string
		input          int64
		setup          func()
		expectedOutput dto.Product
		expectedErr    error
	}{
		{
			name:  "Success",
			input: 1,
			setup: func() {
				suite.productRepo.On("GetProductByID", mock.Anything, mock.Anything, int64(1)).Return(repository.Product{
					ID:       1,
					Name:     "XYZ",
					Category: "Premium",
					Price:    100.0,
					Quantity: 10,
				}, nil)
			},
			expectedOutput: dto.Product{
				ID:       1,
				Name:     "XYZ",
				Category: "Premium",
				Price:    100.0,
				Quantity: 10,
			},
			expectedErr: nil,
		},
		{
			name:  "Fail Because Product Not Found",
			input: 1,
			setup: func() {
				suite.productRepo.On("GetProductByID", mock.Anything, mock.Anything, int64(1)).Return(repository.Product{}, nil)
			},
			expectedOutput: dto.Product{},
			expectedErr:    apperrors.ProductNotFound{ID: 1},
		},
	}

	for _, test := range testCases {
		suite.SetupTest()
		suite.Run(test.name, func() {
			test.setup()

			product, err := suite.service.GetProductByID(context.Background(), nil, test.input)
			suite.Equal(test.expectedErr, err)
			suite.Equal(test.expectedOutput.ID, product.ID)
			suite.Equal(test.expectedOutput.Name, product.Name)
			suite.Equal(test.expectedOutput.Category, product.Category)
			suite.Equal(test.expectedOutput.Price, product.Price)
			suite.Equal(test.expectedOutput.Quantity, product.Quantity)
		})
		suite.TearDownTest()
	}
}

func (suite *ProductServiceTestSuite) TestListProducts() {

	testCases := []struct {
		name           string
		setup          func()
		expectedOutput []dto.Product
		expectedErr    error
	}{
		{
			name: "Success",
			setup: func() {
				suite.productRepo.On("ListProducts", mock.Anything, mock.Anything).Return([]repository.Product{{
					ID:       1,
					Name:     "XYZ",
					Category: "Premium",
					Price:    100.0,
					Quantity: 10,
				},
				}, nil)
			},
			expectedOutput: []dto.Product{{
				ID:       1,
				Name:     "XYZ",
				Category: "Premium",
				Price:    100.0,
				Quantity: 10,
			}},
			expectedErr: nil,
		},
		{
			name: "Fail Because DB Query Failed",
			setup: func() {
				suite.productRepo.On("ListProducts", mock.Anything, mock.Anything).Return([]repository.Product{}, errors.New("Something went wrong in db"))
			},
			expectedOutput: []dto.Product{},
			expectedErr:    errors.New("Something went wrong in db"),
		},
	}

	for _, test := range testCases {
		suite.SetupTest()
		suite.Run(test.name, func() {
			test.setup()

			product, err := suite.service.ListProducts(context.Background())
			suite.Equal(test.expectedErr, err)
			suite.Equal(len(test.expectedOutput), len(product))
		})
		suite.TearDownTest()
	}
}

func (suite *ProductServiceTestSuite) TestUpdateProductQuantity() {

	testCases := []struct {
		name        string
		input       map[int64]int64
		setup       func()
		expectedErr error
	}{
		{
			name:  "Success",
			input: map[int64]int64{},
			setup: func() {
				suite.productRepo.On("UpdateProductQuantity", mock.Anything, mock.Anything, map[int64]int64{1: 2}).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:  "Fail Because DB Query Failed",
			input: map[int64]int64{1: 2},
			setup: func() {
				suite.productRepo.On("UpdateProductQuantity", mock.Anything, mock.Anything, map[int64]int64{1: 2}).Return(errors.New("Something went wrong in db"))
			},
			expectedErr: errors.New("Something went wrong in db"),
		},
	}

	for _, test := range testCases {
		suite.SetupTest()
		suite.Run(test.name, func() {
			test.setup()

			err := suite.service.UpdateProductQuantity(context.Background(), nil, map[int64]int64{1: 2})
			suite.Equal(test.expectedErr, err)
		})
		suite.TearDownTest()
	}
}
