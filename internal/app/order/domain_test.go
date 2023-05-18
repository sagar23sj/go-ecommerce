package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateUpdateOrderStatusRequest(t *testing.T) {
	testCases := []struct {
		name            string
		requestedStatus string
		currentStatus   string
		expectedOutput  bool
	}{
		{
			name:            "Valid Order Status Request, Dispatch Placed Order",
			requestedStatus: "Dispatched",
			currentStatus:   "Placed",
			expectedOutput:  true,
		},
		{
			name:            "Valid Order Status Request, Cancel Placed Order",
			requestedStatus: "Cancelled",
			currentStatus:   "Placed",
			expectedOutput:  true,
		},
		{
			name:            "Valid Order Status Request, Cancel Dispatched Order",
			requestedStatus: "Cancelled",
			currentStatus:   "Dispatched",
			expectedOutput:  true,
		},
		{
			name:            "Valid Order Status Request, Complete Dispatched Order",
			requestedStatus: "Completed",
			currentStatus:   "Dispatched",
			expectedOutput:  true,
		},
		{
			name:            "Valid Order Status Request, Retun Completed Order",
			requestedStatus: "Returned",
			currentStatus:   "Completed",
			expectedOutput:  true,
		},
		{
			name:            "Incorrect Order Status Request, Cannot jump from Placed to Complete",
			requestedStatus: "Completed",
			currentStatus:   "Placed",
			expectedOutput:  false,
		},
		{
			name:            "Incorrect Order Status Request, Cannot Cancel Completed Order",
			requestedStatus: "Cancelled",
			currentStatus:   "Completed",
			expectedOutput:  false,
		},
		{
			name:            "Incorrect Order Status Request, Cannot Cancel Returned Order",
			requestedStatus: "Cancelled",
			currentStatus:   "Returned",
			expectedOutput:  false,
		},
		{
			name:            "Incorrect Order Status Request, Cannot Place Cancelled Order",
			requestedStatus: "Placed",
			currentStatus:   "Cancelled",
			expectedOutput:  false,
		},
		{
			name:            "Incorrect Order Status Request, Placed order Can't Be Placed Again",
			requestedStatus: "Placed",
			currentStatus:   "Placed",
			expectedOutput:  false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			result := validateUpdateOrderStatusRequest(test.requestedStatus, test.currentStatus)
			assert.Equal(t, test.expectedOutput, result)
		})
	}
}
