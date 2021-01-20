package fiborest

import (
	"math"
	"net/http"
	"strconv"

	"github.com/dvaganov/fiboser/pkg/fibonacci"
)

type (
	Controller struct {
		service fibonacci.Service
	}

	FibonacciNumber struct {
		N   uint8  `json:"n"`
		Val string `json:"value"`
	}

	FibonacciResponse struct {
		List []FibonacciNumber `json:"list"`
	}
)

func NewController(service fibonacci.Service) Controller {
	return Controller{service: service}
}

const (
	queryFrom = "from"
	queryTo   = "to"
)

func (c *Controller) GetFibonacciRange(w http.ResponseWriter, r *http.Request) {
	NewHttpRequestHandler(c.getFibonacciRange)(w, r)
}

func (c *Controller) getFibonacciRange(r *http.Request) (interface{}, error) {
	from, err := strconv.Atoi(r.URL.Query().Get(queryFrom))
	if err != nil {
		return nil, NewRequestError()
	}

	to, err := strconv.Atoi(r.URL.Query().Get(queryTo))
	if err != nil {
		return nil, NewRequestError()
	}

	if from < 0 || from > math.MaxUint8 || from > to || to < 0 || to > math.MaxUint8 {
		return nil, NewRequestError()
	}

	res, err := c.service.GetRange(r.Context(), uint8(from), uint8(to))
	if err != nil {
		return nil, err
	}

	resp := FibonacciResponse{
		List: make([]FibonacciNumber, len(res)),
	}

	for i, val := range res {
		resp.List[i].N = uint8(from + i)
		resp.List[i].Val = val
	}

	return resp, nil
}
