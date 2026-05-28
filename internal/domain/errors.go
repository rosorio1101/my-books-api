package domain

import (
	"errors"
	"fmt"
)

var ErrBookNotFound = errors.New("Book not found")

var ErrCartEmpty = errors.New("Cart is empty")

var ErrCartNotFound = errors.New("Cart not found")

var ErrInvalidQuantity = errors.New("Invalid quantity, must be greater than zero")

var ErrAlreadyFavorited = errors.New("Book is already in favorites")

var ErrFavoriteNotFound = errors.New("Favorite not found")

var ErrItemNotInCart = errors.New("Item not in cart")

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Failed validation on field %s: %s", e.Field, e.Message)
}

type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
}

func isNotFound(err error) bool {
	var notFound *NotFoundError
	return errors.As(err, &notFound)
}

func isValidationError(err error) bool {
	var validErr *ValidationError
	return errors.As(err, &validErr)
}
