package fibonacci

import (
	"context"
	"errors"
)

//mockgen -source=./pkg/fibonacci/fibonacci.go -destination=./pkg/fibonacci/fibonaccimock/mock.go -package=fibonaccimock
type (
	Service interface {
		GetRange(ctx context.Context, from, to uint8) ([]string, error)
	}

	Cache interface {
		Get(ctx context.Context, n uint8) (string, error)
		Save(ctx context.Context, n uint8, val string) error
	}
)

var (
	ErrInvalidRange = errors.New("invalid range")
)
