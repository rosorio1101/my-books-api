package domain

import (
	"fmt"
	"strings"
)

type Book struct {
	ID             string
	Title          string
	Author         string
	Description    string
	Price          float64
	FormattedPrice string
	CoverURL       string
	Rating         float64
	RatingCount    int
	Genres         []string
	BookURL        string
}

func NewBook(id, title, author string, price float64) (*Book, error) {
	if strings.TrimSpace(id) == "" {
		return nil, &ValidationError{
			Field:   "id",
			Message: "el ID del libro no puede estar vacío",
		}
	}
	if strings.TrimSpace(title) == "" {
		return nil, &ValidationError{
			Field:   "title",
			Message: "el título del libro no puede estar vacío",
		}
	}
	if strings.TrimSpace(author) == "" {
		return nil, &ValidationError{
			Field:   "author",
			Message: "el autor del libro no puede estar vacío",
		}
	}
	if price < 0 {
		return nil, &ValidationError{
			Field:   "price",
			Message: "el precio no puede ser negativo",
		}
	}
	return &Book{
		ID:     id,
		Title:  title,
		Author: author,
		Price:  price,
	}, nil
}

func (b *Book) String() string {
	return fmt.Sprintf("%s by %s ($%.2f)", b.Title, b.Author, b.Price)
}

func (b *Book) HasRating() bool {
	return b.RatingCount > 0
}

func (b *Book) IsFree() bool {
	return b.Price == 0
}

func (b *Book) GenresString() string {
	return strings.Join(b.Genres, ", ")
}

func (b *Book) WithDescription(description string) *Book {
	b.Description = description
	return b
}

func (b *Book) WithCoverURL(url string) *Book {
	b.CoverURL = url
	return b
}

func (b *Book) WithRating(rating float64, count int) *Book {
	b.Rating = rating
	b.RatingCount = count
	return b
}

func (b *Book) WithGenres(genres []string) *Book {
	b.Genres = make([]string, len(genres))
	copy(b.Genres, genres)
	return b
}

func (b *Book) WithBookURL(url string) *Book {
	b.BookURL = url
	return b
}

func (b *Book) WithFormattedPrice(formatted string) *Book {
	b.FormattedPrice = formatted
	return b
}
