package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("📚 Bookstore API - Go Interview Lab")
	fmt.Println("====================================")
	fmt.Printf("Arquitectura: Hexagonal (Ports & Adapters)\n")
	fmt.Printf("Go version: solo librería estándar\n")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Puerto configurado: %s\n", port)
}
