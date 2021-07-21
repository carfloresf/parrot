package model

import "time"

type GenerateReportRequest struct {
	From time.Time `json:"from,omitempty"`
	To   time.Time `json:"to,omitempty"`
}

type GenerateReportResponse struct {
	Data []Data `json:"data,omitempty"`
}

type Data struct {
	Name        string `json:"name,omitempty"`
	TotalAmount int    `json:"totalAmount,omitempty"`
	TotalPrice  int64  `json:"totalPrice,omitempty"`
}
