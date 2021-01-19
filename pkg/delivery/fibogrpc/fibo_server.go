package fibogrpc

import (
	"context"
	"math"

	"github.com/dvaganov/fiboser/pkg/fibonacci"
)

type (
	Server struct {
		UnimplementedFibonacciServer
		service fibonacci.Service
	}
)

func NewFiboServer(service fibonacci.Service) *Server {
	return &Server{
		UnimplementedFibonacciServer{},
		service,
	}
}

func (s *Server) GetRange(ctx context.Context, r *FibonacciRangeRequest) (*FibonacciRangeResponse, error) {
	if r.GetFrom() > math.MaxUint8 ||
		r.GetTo() > math.MaxUint8 ||
		r.GetFrom() > r.GetTo() {
		return nil, fibonacci.ErrInvalidRange
	}

	res, err := s.service.GetRange(ctx, uint8(r.GetFrom()), uint8(r.GetTo()))
	if err != nil {
		return nil, err
	}

	list := make([]*FibonacciNumber, len(res))
	for i, val := range res {
		list[i] = &FibonacciNumber{
			N:     r.GetFrom() + uint32(i),
			Value: val,
		}
	}

	return &FibonacciRangeResponse{List: list}, nil
}
