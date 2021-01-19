package fibogrpc

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/dvaganov/fiboser/pkg/fibonacci"
	"github.com/dvaganov/fiboser/pkg/fibonacci/fibonaccimock"
	"github.com/golang/mock/gomock"
)

const (
	bufSize = 1024 * 1024
)

func TestGrpcGetRange(t *testing.T) {
	ctx := context.Background()

	errInternal := errors.New("internal error")

	testCases := []struct {
		name       string
		request    *FibonacciRangeRequest
		mockRes    []string
		mockErr    error
		expectResp *FibonacciRangeResponse
		expectErr  error
	}{
		{
			name:       "out of range - 'from' param",
			request:    &FibonacciRangeRequest{From: 1000, To: 10},
			mockRes:    nil,
			mockErr:    nil,
			expectResp: nil,
			expectErr:  fibonacci.ErrInvalidRange,
		},
		{
			name:       "out of range - 'to' param",
			request:    &FibonacciRangeRequest{From: 10, To: 1000},
			mockRes:    nil,
			mockErr:    nil,
			expectResp: nil,
			expectErr:  fibonacci.ErrInvalidRange,
		},
		{
			name:       "invalid range",
			request:    &FibonacciRangeRequest{From: 50, To: 10},
			mockRes:    nil,
			mockErr:    nil,
			expectResp: nil,
			expectErr:  fibonacci.ErrInvalidRange,
		},
		{
			name:       "internal error",
			request:    &FibonacciRangeRequest{From: 0, To: 10},
			mockRes:    nil,
			mockErr:    errInternal,
			expectResp: nil,
			expectErr:  errInternal,
		},
		{
			name:    "success",
			request: &FibonacciRangeRequest{From: 1, To: 2},
			mockRes: []string{"10", "10"},
			mockErr: nil,
			expectResp: &FibonacciRangeResponse{List: []*FibonacciNumber{
				{N: 1, Value: "10"},
				{N: 2, Value: "10"},
			}},
			expectErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockService := fibonaccimock.NewMockService(ctrl)
			fibosrv := NewFiboServer(mockService)

			if tc.mockRes != nil || tc.mockErr != nil {
				mockService.EXPECT().
					GetRange(ctx, gomock.Any(), gomock.Any()).
					Return(tc.mockRes, tc.mockErr)
			}

			resp, err := fibosrv.GetRange(ctx, tc.request)

			if tc.expectErr != err {
				t.Errorf("expect %v, got %v", tc.expectErr, err)
			}

			if !reflect.DeepEqual(tc.expectResp, resp) {
				t.Errorf("expect %v, got %v", tc.expectResp, resp)
			}
		})
	}
}
