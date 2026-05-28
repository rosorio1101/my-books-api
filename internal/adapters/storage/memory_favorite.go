package storage

import (
	"context"
	"my-books-api/internal/domain"
	"sync"
)

type MemoryFavoriteStore struct {
	mu        sync.RWMutex
	favorites map[string][]*domain.Favorite
}

func NewMemoryFavoriteStore() *MemoryFavoriteStore {
	return &MemoryFavoriteStore{
		favorites: make(map[string][]*domain.Favorite),
	}
}

func (s *MemoryFavoriteStore) GetFavorites(ctx context.Context, userID string) ([]*domain.Favorite, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	favs, ok := s.favorites[userID]
	if !ok || len(favs) == 0 {
		return []*domain.Favorite{}, nil
	}
	result := make([]*domain.Favorite, len(favs))
	copy(result, favs)
	return result, nil
}

func (s *MemoryFavoriteStore) AddFavorite(ctx context.Context, userID string, fav *domain.Favorite) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existing := range s.favorites[userID] {
		if existing.BookID == fav.BookID {
			return domain.ErrAlreadyFavorited
		}
	}

	s.favorites[userID] = append(s.favorites[userID], fav)
	return nil
}

func (s *MemoryFavoriteStore) RemoveFavorite(ctx context.Context, userID string, bookID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	favs, ok := s.favorites[userID]
	if !ok {
		return domain.ErrFavoriteNotFound
	}
	idx := -1
	for i, f := range favs {
		if f.BookID == bookID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return domain.ErrFavoriteNotFound
	}
	s.favorites[userID] = append(favs[:idx], favs[idx+1:]...)
	return nil
}

func (s *MemoryFavoriteStore) IsFavorite(ctx context.Context, userID string, bookID string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	favs, ok := s.favorites[userID]
	if !ok {
		return false, nil
	}

	for _, f := range favs {
		if f.BookID == bookID {
			return true, nil
		}
	}
	return false, nil
}
