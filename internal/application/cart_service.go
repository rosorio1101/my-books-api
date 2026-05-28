package application

import (
	"context"
	"errors"
	"fmt"
	"my-books-api/internal/domain"
	"my-books-api/internal/ports"
	"time"
)

type CartService struct {
	repo      ports.CartRepository
	publisher ports.EventPublisher
}

func NewCartService(repo ports.CartRepository, pub ports.EventPublisher) *CartService {
	return &CartService{repo: repo, publisher: pub}
}

func (s *CartService) GetCart(ctx context.Context, cartID string) (*domain.Cart, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("obteniendo carrito: %w", ctx.Err())
	default:
	}

	cart, err := s.repo.GetCart(ctx, cartID)
	if err != nil {
		if errors.Is(err, domain.ErrCartNotFound) {
			newCart := domain.NewCart()
			return newCart, nil
		}

		return nil, fmt.Errorf("obteniendo carrito %s: %w", cartID, err)
	}

	return cart, nil
}

func (s *CartService) AddItem(ctx context.Context, cartID string, book *domain.Book, qty int) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("agregando item al carrito: %w", ctx.Err())
	default:
	}
	cart, err := s.GetCart(ctx, cartID)
	if err != nil {
		return fmt.Errorf("agregando item al carrito: %w", err)
	}

	if err := cart.AddItem(book, qty); err != nil {
		return fmt.Errorf("agregando item al carrito: %w", err)
	}

	if err := s.repo.SaveCart(ctx, cart); err != nil {
		return fmt.Errorf("guardando carrito tras agregar item: %w", err)
	}
	return nil
}

func (s *CartService) RemoveItem(ctx context.Context, cartID string, bookID string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("eliminando item del carrito: %w", ctx.Err())
	default:
	}
	cart, err := s.repo.GetCart(ctx, cartID)
	if err != nil {
		return fmt.Errorf("eliminando item del carrito: %w", err)
	}

	if err := cart.RemoveItem(bookID); err != nil {
		return fmt.Errorf("eliminando item del carrito: %w", err)
	}

	if err := s.repo.SaveCart(ctx, cart); err != nil {
		return fmt.Errorf("guardando carrito tras eliminar item: %w", err)
	}
	return nil
}

func (s *CartService) UpdateItemQuantity(ctx context.Context, cartID string, bookID string, qty int) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("actualizando cantidad del item: %w", ctx.Err())
	default:
	}
	if qty <= 0 {
		return s.RemoveItem(ctx, cartID, bookID)
	}

	cart, err := s.repo.GetCart(ctx, cartID)
	if err != nil {
		return fmt.Errorf("actualizando cantidad del item: %w", err)
	}

	if err := cart.UpdateQuantity(bookID, qty); err != nil {
		return fmt.Errorf("actualizando cantidad del item: %w", err)
	}

	if err := s.repo.SaveCart(ctx, cart); err != nil {
		return fmt.Errorf("guardando carrito tras actualizar cantidad: %w", err)
	}
	return nil
}

func (s *CartService) ClearCart(ctx context.Context, cartID string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("limpiando carrito: %w", ctx.Err())
	default:
	}

	if err := s.repo.DeleteCart(ctx, cartID); err != nil {
		return fmt.Errorf("limpiando carrito %s: %w", cartID, err)
	}
	return nil
}

func (s *CartService) Checkout(ctx context.Context, cartID string) (*domain.Cart, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("procesando checkout: %w", ctx.Err())
	default:
	}
	cart, err := s.repo.GetCart(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("procesando checkout: %w", err)
	}
	if len(cart.Items) == 0 {
		return nil, fmt.Errorf("procesando checkout: %w", domain.ErrCartEmpty)
	}
	checkoutEvent := ports.Event{
		Type:       "cart.checkout",
		Payload:    cart,
		OccurredAt: time.Now(),
	}
	if err := s.publisher.Publish(ctx, checkoutEvent); err != nil {
		return nil, fmt.Errorf("publicando evento de checkout: %w", err)
	}
	if err := s.repo.DeleteCart(ctx, cartID); err != nil {
		return nil, fmt.Errorf("limpiando carrito tras checkout: %w", err)
	}
	return cart, nil
}
