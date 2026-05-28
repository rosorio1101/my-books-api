package storage

import (
	"context"
	"my-books-api/internal/domain"
	"sync"
)

type MemoryCartStore struct {
	mu    sync.RWMutex
	carts map[string]*domain.Cart
}

func NewMemoryCartStore() *MemoryCartStore {
	return &MemoryCartStore{
		carts: make(map[string]*domain.Cart),
	}
}

func (s *MemoryCartStore) GetCart(ctx context.Context, cartID string) (*domain.Cart, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	cart, ok := s.carts[cartID]
	if !ok {
		return nil, domain.ErrCartNotFound
	}
	return cart, nil
}

func (s *MemoryCartStore) SaveCart(ctx context.Context, cart *domain.Cart) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.carts[cart.ID] = cart
	return nil
}

func (s *MemoryCartStore) DeleteCart(ctx context.Context, cartID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.carts, cartID)
	return nil
}
