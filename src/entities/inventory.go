package entities

import (
	"time"
)

type Inventory struct {
	Id           int64
	ProductId    int64
	Quantity     int64
	PurchaseDate time.Time
}
