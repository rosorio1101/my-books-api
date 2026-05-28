package ports

import (
	"context"
	"my-books-api/internal/domain"
	"time"
)

type Event struct {
	Type       string
	Payload    any
	OccurredAt time.Time
}

type BookRepository interface {
	SearchBook(ctx context.Context, query string, limit int) ([]*domain.Book, error)
	GetBookById(ctx context.Context, id string) (*domain.Book, error)
	SaveBook(ctx context.Context, book *domain.Book) (*domain.Book, error)
}

type CartRepository interface {
	GetCart(ctx context.Context, id string) (*domain.Cart, error)
	SaveCart(ctx context.Context, cart *domain.Cart) error
	DeleteCart(ctx context.Context, id string) error
}

type FavoriteRepository interface {
	GetFavorites(ctx context.Context, userID string) ([]*domain.Favorite, error)
	AddFavorite(ctx context.Context, userID string, favorite *domain.Favorite) error
	RemoveFavorite(ctx context.Context, userID string, bookID string) error
	IsFavorite(ctx context.Context, userID string, bookID string) (bool, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
}
