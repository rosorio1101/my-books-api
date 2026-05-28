package domain

import (
	"fmt"
	"time"
)

type CartItem struct {
	BookID   string
	Title    string
	Price    float64
	Quantity int
}

func (ci CartItem) Subtotal() float64 {
	return ci.Price * float64(ci.Quantity)
}

func (ci CartItem) String() string {
	return fmt.Sprintf("%s (x%d) - $%.2f", ci.Title, ci.Quantity, ci.Subtotal())
}

type Cart struct {
	ID        string
	Items     []CartItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCart() *Cart {
	now := time.Now()
	return &Cart{
		// Usamos UnixNano como ID simple.
		// time.Now().UnixNano() da nanosegundos desde epoch — suficientemente único
		// para una app single-instance.
		ID:        fmt.Sprintf("cart-%d", now.UnixNano()),
		Items:     make([]CartItem, 0), // slice vacío (no nil) — importante para JSON
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (c *Cart) AddItem(book *Book, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	for i := range c.Items {
		if c.Items[i].BookID == book.ID {
			c.Items[i].Quantity += quantity
			c.UpdatedAt = time.Now()
			return nil
		}
	}
	c.Items = append(c.Items, CartItem{
		BookID:   book.ID,
		Title:    book.Title,
		Price:    book.Price,
		Quantity: quantity,
	})
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Cart) RemoveItem(bookID string) error {
	// Buscamos el índice del ítem a eliminar
	index := -1
	for i, item := range c.Items {
		if item.BookID == bookID {
			index = i
			break
		}
	}

	if index == -1 {
		return ErrItemNotInCart
	}

	c.Items = append(c.Items[:index], c.Items[index+1:]...)
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Cart) UpdateQuantity(bookID string, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	for i := range c.Items {
		if c.Items[i].BookID == bookID {
			c.Items[i].Quantity = quantity
			c.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrItemNotInCart
}

func (c *Cart) Total() float64 {
	var total float64
	for _, item := range c.Items {
		// Usamos el método Subtotal() del CartItem para mantener
		// la lógica de cálculo en un solo lugar (DRY - Don't Repeat Yourself)
		total += item.Subtotal()
	}
	return total
}

func (c *Cart) ItemCount() int {
	count := 0
	for _, item := range c.Items {
		count += item.Quantity
	}
	return count
}

func (c *Cart) UniqueItemCount() int {
	return len(c.Items)
}

func (c *Cart) IsEmpty() bool {
	return len(c.Items) == 0
}

func (c *Cart) Clear() {
	c.Items = make([]CartItem, 0)
	c.UpdatedAt = time.Now()
}

func (c *Cart) FindItem(bookID string) (*CartItem, bool) {
	for i := range c.Items {
		if c.Items[i].BookID == bookID {
			// Retornamos &c.Items[i] — un puntero al elemento dentro del slice.
			// Esto permite que quien llame pueda modificar el ítem directamente.
			// CUIDADO: si el slice crece y Go reasigna memoria, este puntero
			// podría quedar inválido. En nuestro caso es seguro porque solo
			// lo usamos inmediatamente después de la búsqueda.
			return &c.Items[i], true
		}
	}
	return nil, false
}

func (c *Cart) String() string {
	return fmt.Sprintf("Cart{ID: %s, Items: %d, Total: $%.2f}",
		c.ID, c.ItemCount(), c.Total())
}
