package rest

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dvaganov/fiboser/pkg/fibonacci/fibonaccimock"
	"github.com/golang/mock/gomock"
)

func TestGetFibonacciRange(t *testing.T) {
	respBadRequest := `{"message":"bad request"}`
	respInternalErr := `{"message":"internal server error"}`
	respSuccess := `{"list":[{"n":5,"value":"1"},{"n":6,"value":"1"}]}`

	testCases := []struct {
		name     string
		from, to interface{}
		mockRes  []string
		mockErr  error
		status   int
		res      string
	}{
		{"no from param", nil, 10, nil, nil, http.StatusBadRequest, respBadRequest},
		{"no to param", 10, nil, nil, nil, http.StatusBadRequest, respBadRequest},
		{"no query param", nil, nil, nil, nil, http.StatusBadRequest, respBadRequest},
		{"invalid from param", "invalid", 10, nil, nil, http.StatusBadRequest, respBadRequest},
		{"invalid to param", 10, "invalid", nil, nil, http.StatusBadRequest, respBadRequest},
		{"out of range", -10, 300, nil, nil, http.StatusBadRequest, respBadRequest},
		{"invalid range", 50, 10, nil, nil, http.StatusBadRequest, respBadRequest},
		{"internal error", 10, 50, nil, errors.New("internal error"), http.StatusInternalServerError, respInternalErr},
		{"success", 5, 10, []string{"1", "1"}, nil, http.StatusOK, respSuccess},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := fibonaccimock.NewMockService(ctrl)

			fiboctrl := NewController(mock)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			q := r.URL.Query()
			if tc.from != nil {
				q.Add(queryFrom, fmt.Sprint(tc.from))
			}

			if tc.to != nil {
				q.Add(queryTo, fmt.Sprint(tc.to))
			}
			r.URL.RawQuery = q.Encode()

			if tc.mockRes != nil || tc.mockErr != nil {
				mock.EXPECT().
					GetRange(r.Context(), gomock.Any(), gomock.Any()).
					Return(tc.mockRes, tc.mockErr)
			}

			fiboctrl.GetFibonacciRange(w, r)

			if tc.status != w.Code {
				t.Errorf("expect %d, got %d", tc.status, w.Code)
			}

			if tc.res != w.Body.String() {
				t.Errorf("expect %s, got %s", tc.res, w.Body.String())
			}
		})
	}
}
