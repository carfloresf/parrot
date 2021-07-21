package storage

import (
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/hellerox/parrot/model"
)

// DatabaseStorage with config data
type DatabaseStorage struct {
	pool pgxpoolmock.PgxPool
}

// Storage executes functions on storage resources
type Storage interface {
	InsertUser(u model.User) error
	GetUserHash(mail string) string
	InsertOrder(o model.Order) (int, error)
	InsertProduct(p model.Product) (int, error)
	InsertOrderProductRelation(opr model.OrderProductRelation) (int, error)
	UpdateOrder(o model.Order) error
	GetReportData(r model.GenerateReportRequest) (model.GenerateReportResponse, error)
}

// NewStorage returns a new DatabaseOperator
func NewStorage(connectionString string) Storage {
	storage := DatabaseStorage{
		pool: connect(connectionString),
	}

	return &storage
}
