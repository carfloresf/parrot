package service

import (
	"github.com/hellerox/parrot/model"
	"github.com/hellerox/parrot/storage"
)

// Service controller
type Service struct {
	Storage storage.Storage
}

// CreateUser creates users
func (s *Service) CreateUser(ms model.User) (err error) {
	err = s.Storage.InsertUser(ms)
	if err != nil {
		return err
	}

	return nil
}

// CreateOrder creates orders
func (s *Service) CreateOrder(o model.Order) (int, error) {
	var totalOrder int64

	oId, err := s.Storage.InsertOrder(o)
	if err != nil {
		return 0, err
	}

	for _, p := range o.Products {
		idp, err := s.Storage.InsertProduct(p)
		if err != nil {
			return 0, err
		}

		op := model.OrderProductRelation{
			ProductID: idp,
			OrderID:   oId,
			Amount:    p.Amount,
		}

		totalProduct := int64(p.Amount) * p.Price
		totalOrder = totalProduct + totalOrder

		_, err = s.Storage.InsertOrderProductRelation(op)
		if err != nil {
			return 0, err
		}

	}

	o.ID = oId
	o.Price = totalOrder
	err = s.Storage.UpdateOrder(o)
	if err != nil {
		return oId, err
	}

	return oId, nil
}
