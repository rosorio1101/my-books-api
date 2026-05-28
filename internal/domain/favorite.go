package domain

import (
	"fmt"
	"time"
)

type Favorite struct {
	ID      string
	BookID  string
	Book    *Book
	AddedAt time.Time
}

func NewFavorite(book *Book) (*Favorite, error) {
	if book == nil {
		return nil, &ValidationError{
			Field:   "book",
			Message: "el libro no puede ser nil",
		}
	}
	if book.ID == "" {
		return nil, &ValidationError{
			Field:   "book.ID",
			Message: "el libro debe tener un ID válido",
		}
	}
	now := time.Now()
	return &Favorite{
		ID:      generateFavoriteID(now),
		BookID:  book.ID,
		Book:    book,
		AddedAt: now,
	}, nil
}

func generateFavoriteID(t time.Time) string {
	// UnixNano retorna los nanosegundos desde epoch (1 de enero de 1970 UTC).
	// Es suficientemente único para una aplicación single-instance.
	return "fav-" + fmt.Sprintf("%d", t.UnixNano())
}
