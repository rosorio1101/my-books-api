package application

import (
	"context"
	"fmt"
	"my-books-api/internal/domain"
	"my-books-api/internal/ports"
)

type FavoriteService struct {
	// repo maneja la persistencia de favoritos
	repo ports.FavoriteRepository
}

func NewFavoriteService(repo ports.FavoriteRepository) *FavoriteService {
	return &FavoriteService{
		repo: repo,
	}
}

func (s *FavoriteService) GetFavorites(ctx context.Context, userID string) ([]*domain.Favorite, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("obteniendo favoritos: %w", ctx.Err())
	default:
	}

	if userID == "" {
		return nil, fmt.Errorf("obteniendo favoritos: userID no puede estar vacío")
	}
	favs, err := s.repo.GetFavorites(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("obteniendo favoritos del usuario %s: %w", userID, err)
	}
	return favs, nil
}

func (s *FavoriteService) AddFavorite(ctx context.Context, userID string, book *domain.Book) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("agregando favorito: %w", ctx.Err())
	default:
	}
	if userID == "" {
		return fmt.Errorf("agregando favorito: userID no puede estar vacío")
	}
	if book == nil {
		return fmt.Errorf("agregando favorito: el libro no puede ser nil")
	}
	fav, err := domain.NewFavorite(book)
	if err != nil {
		return fmt.Errorf("agregando favorito: %w", err)
	}
	if err := s.repo.AddFavorite(ctx, userID, fav); err != nil {
		return fmt.Errorf("agregando libro %s a favoritos: %w", book.ID, err)
	}
	return nil
}

func (s *FavoriteService) RemoveFavorite(ctx context.Context, userID string, bookID string) error {

	select {
	case <-ctx.Done():
		return fmt.Errorf("eliminando favorito: %w", ctx.Err())
	default:
	}
	if userID == "" {
		return fmt.Errorf("eliminando favorito: userID no puede estar vacío")
	}
	if bookID == "" {
		return fmt.Errorf("eliminando favorito: bookID no puede estar vacío")
	}
	if err := s.repo.RemoveFavorite(ctx, userID, bookID); err != nil {
		return fmt.Errorf("eliminando libro %s de favoritos: %w", bookID, err)
	}
	return nil
}
