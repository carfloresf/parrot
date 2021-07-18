package storage

import (
	"github.com/jackc/pgx/v4"

	"github.com/hellerox/parrot/model"
)

// DatabaseStorage with config data
type DatabaseStorage struct {
	conn *pgx.Conn
}

// Storage executes functions on storage resources
type Storage interface {
	InsertUser(c model.User) error
	InsertOrder(m model.Order) (int, error)
	InsertProduct(m model.Product) (int, error)
	InsertOrderProductRelation(cm model.OrderProductRelation) (int, error)
	UpdateOrder(o model.Order) error
}

// NewStorage returns a new DatabaseOperator
func NewStorage(connectionString string) Storage {
	storage := DatabaseStorage{
		conn: connect(connectionString),
	}

	storage.PrepareStatements()

	return &storage
}
