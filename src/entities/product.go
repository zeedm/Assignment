package entities

type Product struct {
	Id            int64
	Name          string
	Price         float64
	Quantity      int64
	VendorId      int64
	ProductTypeId int64
}
