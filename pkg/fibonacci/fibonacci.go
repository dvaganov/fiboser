package fibonacci

import (
	"context"
	"errors"
)

type (
	Service interface {
		GetRange(ctx context.Context, from, to uint8) ([]string, error)
	}
)

var (
	ErrInvalidRange = errors.New("invalid range")
)
