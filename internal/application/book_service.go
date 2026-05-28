package application

import (
	"context"
	"fmt"
	"my-books-api/internal/domain"
	"my-books-api/internal/ports"
)

type BookService struct {
	repo ports.BookRepository
}

func NewBookService(repo ports.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) SearchBooks(ctx context.Context, query string, limit int) ([]*domain.Book, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	books, err := s.repo.SearchBook(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("buscando libros con la query %q: %w", query, err)
	}

	return books, nil
}

func (s *BookService) GetBookByID(ctx context.Context, id string) (*domain.Book, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if id == "" {
		return nil, fmt.Errorf("obteniendo libro: %w", ctx.Err())
	}

	book, err := s.repo.GetBookById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("obtenido libro %s: %w", id, err)
	}

	return book, nil
}
