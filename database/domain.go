package database

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `db:"id" primary:"true"`
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

type Inventory struct {
	ID        uuid.UUID  `db:"id" primary:"true"`
	Name      string     `db:"name"`
	Quantity  int        `db:"quantity"`
	Price     float64    `db:"price"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
