package models

import (
	"api/assignment/src/entities"
	"database/sql"
)

type ProductModel struct {
	Db *sql.DB
}

func (productModel ProductModel) GetAllProducts() ([]entities.Product, error) {
	rows, errorMessage := productModel.Db.Query(`select p.*, sum(i.Quantity) as Quantity from Product p
	left join Inventory i on p.Id = i.ProductId
	where Quantity > 0
	group by p.Id, p.Name, p.Price, p.ProductTypeId, p.VendorId`)
	if errorMessage != nil {
		return nil, errorMessage
	} else {
		var products []entities.Product
		for rows.Next() {
			var id int64
			var name string
			var quantity int64
			var price float64
			var vendorId int64
			var productTypeId int64
			errorMessageProduct := rows.Scan(&id, &name, &price, &vendorId, &productTypeId, &quantity)
			if errorMessageProduct != nil {
				return nil, errorMessageProduct
			} else {
				product := entities.Product{
					Id:            id,
					Name:          name,
					Quantity:      quantity,
					Price:         price,
					VendorId:      vendorId,
					ProductTypeId: productTypeId,
				}
				products = append(products, product)
			}
		}
		return products, nil
	}
}

func (productModel ProductModel) FindProductByVendorId(vendorId int) ([]entities.Product, error) {
	rows, errorMessage := productModel.Db.Query("select * from Product where vendorId = $1", vendorId)
	if errorMessage != nil {
		return nil, errorMessage
	} else {
		var products []entities.Product
		for rows.Next() {
			var id int64
			var name string
			var price float64
			var vendorId int64
			var productTypeId int64
			errorMessageProduct := rows.Scan(&id, &name, &price, &vendorId, &productTypeId)
			if errorMessageProduct != nil {
				return nil, errorMessageProduct
			} else {
				product := entities.Product{
					Id:            id,
					Name:          name,
					Price:         price,
					VendorId:      vendorId,
					ProductTypeId: productTypeId,
				}
				products = append(products, product)
			}
		}
		return products, nil
	}
}
