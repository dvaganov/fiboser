package fibonacci

import (
	"context"
	"math/big"
)

type (
	dynamicService struct {
		cache Cache
	}
)

func NewService(cache Cache) Service {
	return &dynamicService{
		cache: cache,
	}
}

func (m *dynamicService) GetRange(ctx context.Context, from, to uint8) ([]string, error) {
	if from > to {
		return nil, ErrInvalidRange
	}

	fibSlice := make([]string, 0, to-from)

	var f *big.Int
	f1, f2 := big.NewInt(1), big.NewInt(0)
	for i := uint16(0); i <= uint16(to); i++ {
		var ok bool
		if f, ok = m.fromCache(ctx, uint8(i)); !ok {
			f = f2
			if i > 1 {
				f.Add(f1, f2)
			}
			m.toCache(ctx, uint8(i), f)
		}

		if i >= uint16(from) {
			fibSlice = append(fibSlice, f.String())
		}

		f1, f2 = f, f1
	}

	return fibSlice, nil
}

func (m *dynamicService) fromCache(ctx context.Context, n uint8) (*big.Int, bool) {
	if m.cache == nil {
		return nil, false
	}

	val, err := m.cache.Get(ctx, n)
	if err != nil {
		return nil, false
	}

	f := new(big.Int)
	if _, ok := f.SetString(val, 10); !ok {
		return nil, ok
	}

	return f, true
}

func (m *dynamicService) toCache(ctx context.Context, n uint8, f *big.Int) {
	if m.cache == nil {
		return
	}
	m.cache.Save(ctx, n, f.String())
}
