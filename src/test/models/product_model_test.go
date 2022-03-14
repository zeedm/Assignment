package test

import (
	"api/assignment/src/entities"
	"api/assignment/src/models"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func CreateMockProducts() []entities.Product {
	return []entities.Product{
		{
			Id:            1,
			Name:          "testName",
			Price:         5,
			Quantity:      5,
			VendorId:      1,
			ProductTypeId: 1,
		},
		{
			Id:            2,
			Name:          "testNameByVendor1",
			Price:         4,
			Quantity:      4,
			VendorId:      1,
			ProductTypeId: 1,
		},
	}
}

func TestGetAllProduct(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	productModel := models.ProductModel{
		Db: db,
	}
	query := `select p.*, sum(i.Quantity) as Quantity from Product p
	left join Inventory i on p.Id = i.ProductId
	where Quantity > 0
	group by p.Id, p.Name, p.Price, p.ProductTypeId, p.VendorId`

	rows := sqlmock.NewRows([]string{"id", "name", "price", "vendorId", "productTypeId", "quantity"})

	mockProducts := CreateMockProducts()
	for _, p := range mockProducts {
		rows.AddRow(p.Id, p.Name, p.Price, p.VendorId, p.ProductTypeId, p.Quantity)
	}

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
	actualProducts, _ := productModel.GetAllProducts()
	assert.NotEmpty(t, actualProducts)
	assert.Equal(t, 2, len(actualProducts))
	for i, p := range mockProducts {
		assert.Equal(t, p.Id, actualProducts[i].Id)
		assert.Equal(t, p.Name, actualProducts[i].Name)
		assert.Equal(t, p.Price, actualProducts[i].Price)
		assert.Equal(t, p.VendorId, actualProducts[i].VendorId)
		assert.Equal(t, p.ProductTypeId, actualProducts[i].ProductTypeId)
		assert.Equal(t, p.Quantity, actualProducts[i].Quantity)
	}
}

func TestFindProductByVendorId(t *testing.T) {
	vendorId := 1
	db, mock := NewMock()
	defer db.Close()

	productModel := models.ProductModel{
		Db: db,
	}
	query := `select * from Product where vendorId = $1`

	rows := sqlmock.NewRows([]string{"id", "name", "price", "vendorId", "productTypeId"})

	mockProducts := CreateMockProducts()
	for _, p := range mockProducts {
		rows.AddRow(p.Id, p.Name, p.Price, p.VendorId, p.ProductTypeId)
	}

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(vendorId).WillReturnRows(rows)
	actualProducts, _ := productModel.FindProductByVendorId(vendorId)
	assert.NotEmpty(t, actualProducts)
	assert.Equal(t, 2, len(actualProducts))
	for i, p := range mockProducts {
		assert.Equal(t, p.Id, actualProducts[i].Id)
		assert.Equal(t, p.Name, actualProducts[i].Name)
		assert.Equal(t, p.Price, actualProducts[i].Price)
		assert.Equal(t, p.VendorId, actualProducts[i].VendorId)
		assert.Equal(t, p.ProductTypeId, actualProducts[i].ProductTypeId)
	}
}

func TestFindProductByVendorIdNotFound(t *testing.T) {
	vendorId := 1
	db, mock := NewMock()
	defer db.Close()

	productModel := models.ProductModel{
		Db: db,
	}
	query := `select * from Product where vendorId = $1`

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(vendorId).WillReturnError(fmt.Errorf(""))
	_, err := productModel.FindProductByVendorId(vendorId)
	assert.Error(t, err)
}

func TestFindProductByVendorIdEmpty(t *testing.T) {
	vendorId := 1
	db, mock := NewMock()
	defer db.Close()

	productModel := models.ProductModel{
		Db: db,
	}
	query := `select * from Product where vendorId = $1`

	rows := sqlmock.NewRows([]string{"id", "name", "price", "vendorId", "productTypeId"})

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(vendorId).WillReturnRows(rows)
	actualProducts, _ := productModel.FindProductByVendorId(vendorId)
	assert.Empty(t, actualProducts)
}
