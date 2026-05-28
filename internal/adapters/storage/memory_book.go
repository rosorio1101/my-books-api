package storage

import (
	"context"
	"my-books-api/internal/domain"
	"strings"
	"sync"
)

type MemoryBookStore struct {
	mu          sync.RWMutex
	books       map[string]*domain.Book
	searchMu    sync.RWMutex
	searchCache map[string][]*domain.Book
}

func NewMemoryBookStore() *MemoryBookStore {
	return &MemoryBookStore{
		books:       make(map[string]*domain.Book),
		searchCache: make(map[string][]*domain.Book),
	}
}

func (s *MemoryBookStore) SearchBooks(ctx context.Context, query string, limit int) ([]*domain.Book, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	normalizedQuery := strings.ToLower(query)
	if normalizedQuery == "" {
		return []*domain.Book{}, nil
	}

	s.searchMu.RLock()
	if cached, ok := s.searchCache[normalizedQuery]; ok {
		s.searchMu.RUnlock()
		if limit > 0 && len(cached) >= limit {
			return cached[:limit], nil
		}
		return cached, nil
	}
	s.searchMu.RUnlock()

	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*domain.Book
	for _, book := range s.books {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		titleMatch := strings.Contains(strings.ToLower(book.Title), normalizedQuery)
		authorMatch := strings.Contains(strings.ToLower(book.Author), normalizedQuery)

		if titleMatch || authorMatch {
			results = append(results, book)
		}

		if limit > 0 && len(results) >= limit {
			break
		}
	}

	if results == nil {
		results = []*domain.Book{}
	}

	return results, nil
}

func (s *MemoryBookStore) GetBookByID(ctx context.Context, id string) (*domain.Book, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	book, ok := s.books[id]
	if !ok {
		return nil, domain.ErrBookNotFound
	}
	return book, nil
}

func (s *MemoryBookStore) SaveBook(ctx context.Context, book *domain.Book) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.books[book.ID] = book
	return nil
}

func (s *MemoryBookStore) CacheSearchResults(query string, books []*domain.Book) {
	normalizedQuery := strings.ToLower(strings.TrimSpace(query))
	if normalizedQuery == "" {
		return
	}
	s.searchMu.Lock()
	defer s.searchMu.Unlock()
	s.searchCache[normalizedQuery] = books
}
