package ports

import (
	"context"
	"my-books-api/internal/domain"
)

type BookService interface {
	SearchBooks(ctx context.Context, query string, limit int) ([]*domain.Book, error)
	GetBookById(ctx context.Context, id string) (*domain.Book, error)
}

type CartService interface {
	GetCart(ctx context.Context, id string) (*domain.Cart, error)
	AddItem(ctx context.Context, cartID string, book *domain.Book, qty int) error
	RemoveItem(ctx context.Context, cartID string, bookID string) error
	UpdateItemQuantity(ctx context.Context, cartID string, bookID string, quantity int) error
	ClearCart(ctx context.Context, cartID string) error
	Checkout(ctx context.Context, cartID string) (*domain.Cart, error)
}

type FavoriteService interface {
	GetFavorites(ctx context.Context, userID string) ([]*domain.Favorite, error)
	AddFavorite(ctx context.Context, userID string, book *domain.Book) error
	RemoveFavorite(ctx context.Context, userID string, bookID string) error
}
