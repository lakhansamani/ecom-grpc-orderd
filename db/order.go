package db

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	order "github.com/lakhansamani/ecom-grpc-apis/order/v1"
)

// Order represents the Order model in DB
type Order struct {
	ID        string `gorm:"primaryKey"`
	UserID    string
	Product   string
	Quantity  int32
	UnitPrice float64
}

// AsAPIOrder converts the Order model to API Order
func (o *Order) AsAPIOrder() *order.Order {
	return &order.Order{
		Id:        o.ID,
		UserId:    o.UserID,
		Product:   o.Product,
		Quantity:  o.Quantity,
		UnitPrice: o.UnitPrice,
	}
}

// Before create
func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate UUID
	o.ID = uuid.NewString()
	return
}

// CreateOrder creates a new order in the database
func (p *provider) CreateOrder(ctx context.Context, o *Order) (*Order, error) {
	err := p.db.WithContext(ctx).Create(o).Error
	return o, err
}

// GetOrderById fetches a order by ID from the database
func (p *provider) GetOrderById(ctx context.Context, id string) (*Order, error) {
	var o Order
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&o).Error
	return &o, err
}
