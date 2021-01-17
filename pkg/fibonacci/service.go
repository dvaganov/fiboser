package fibonacci

import (
	"context"
	"math/big"
)

type (
	dynamicService struct{}
)

func NewService() Service {
	return &dynamicService{}
}

func (m *dynamicService) GetRange(ctx context.Context, from, to uint8) ([]string, error) {
	if from > to {
		return nil, ErrInvalidRange
	}

	fibSlice := make([]string, 0, to-from)

	var f *big.Int
	f1, f2 := big.NewInt(1), big.NewInt(0)
	for i := uint16(0); i <= uint16(to); i++ {
		f = f2
		if i > 1 {
			f.Add(f1, f2)
		}

		if i >= uint16(from) {
			fibSlice = append(fibSlice, f.String())
		}

		f1, f2 = f, f1
	}

	return fibSlice, nil
}
