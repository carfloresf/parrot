package storage

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/hellerox/parrot/model"
)

// DatabaseStorage with config data
type DatabaseStorage struct {
	pool *pgxpool.Pool
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
