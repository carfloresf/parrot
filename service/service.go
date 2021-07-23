package service

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/hellerox/parrot/model"
	"github.com/hellerox/parrot/storage"
)

const Pepper = "secret-random-string"

// Service controller
type Service struct {
	Storage storage.Storage
}

// CreateUser creates users
func (s *Service) CreateUser(user model.User) (err error) {
	err = hashPassword(&user)
	if err != nil {
		return err
	}

	err = s.Storage.InsertUser(user)
	if err != nil {
		return err
	}

	return nil
}

// CreateOrder creates orders
func (s *Service) CreateOrder(o model.Order) (int, error) {
	var totalOrder int64

	oID, err := s.Storage.InsertOrder(o)
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
			OrderID:   oID,
			Amount:    p.Amount,
		}

		totalProduct := int64(p.Amount) * p.Price
		totalOrder = totalProduct + totalOrder

		_, err = s.Storage.InsertOrderProductRelation(op)
		if err != nil {
			return 0, err
		}
	}

	o.ID = oID
	o.Price = totalOrder

	err = s.Storage.UpdatePriceOrder(o)
	if err != nil {
		return oID, err
	}

	return oID, nil
}

func (s *Service) GenerateReport(r model.GenerateReportRequest) (model.GenerateReportResponse, error) {
	var res model.GenerateReportResponse

	res, err := s.Storage.GetReportData(r)
	if err != nil {
		return res, err
	}

	return res, err
}

func (s *Service) GetUserHash(mail string) (hash string) {
	hash = s.Storage.GetUserHash(mail)
	return hash
}

func hashPassword(user *model.User) error {
	if user.Password == "" {
		return nil
	}

	pwBytes := []byte(user.Password + Pepper)

	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	return nil
}
