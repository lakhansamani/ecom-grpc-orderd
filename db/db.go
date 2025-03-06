package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Provider defines the interface for the database provider
type Provider interface {
	CreateOrder(order *Order) (*Order, error)
	GetOrderById(id string) (*Order, error)
}

// provider implements the Provider interface
type provider struct {
	db *gorm.DB
}

// New creates new database provider
// connects to db and returns the provider
func New(dbURL string) Provider {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate User model
	db.AutoMigrate(&Order{})

	return &provider{db}
}
