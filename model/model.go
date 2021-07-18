package model

import "time"

type User struct {
	Email    string `db:"email" json:"email,omitempty"`
	FullName string `db:"full_name" json:"fullName,omitempty"`
	Password string `db:"password" json:"password,omitempty"`
}

type Order struct {
	ID         int    `db:"id" json:"id"`
	Email      string `db:"email" json:"email,omitempty"`
	ClientName string `db:"client_name" json:"clientName,omitempty"`
	Price      int64  `db:"price" json:"price,omitempty"`
	CreatedAt  time.Time
	Products   []Product `json:"products,omitempty"`
}

// Product is
type Product struct {
	ID          int    `db:"id" json:"id,omitempty"`
	Name        string `json:"name,omitempty" db:"name"`
	Price       int64  `json:"price,omitempty" db:"price"`
	Description string `json:"description,omitempty" db:"description"`
	Amount      int    `json:"amount,omitempty" db:"amount"`
}

// OrderProductRelation is
type OrderProductRelation struct {
	ID        int `db:"id" json:"id,omitempty"`
	ProductID int `json:"productID,omitempty" db:"product_id"`
	OrderID   int `json:"orderID,omitempty" db:"order_id"`
	Amount    int `json:"amount,omitempty" db:"amount"`
}
