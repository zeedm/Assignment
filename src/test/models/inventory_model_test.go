package test

import (
	"api/assignment/src/entities"
	"api/assignment/src/models"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func CreateMockInventorys() []entities.Inventory {
	return []entities.Inventory{
		{
			Id:           1,
			ProductId:    1,
			Quantity:     5,
			PurchaseDate: time.Now(),
		},
		{
			Id:           2,
			ProductId:    1,
			Quantity:     1,
			PurchaseDate: time.Now(),
		},
		{
			Id:           3,
			ProductId:    2,
			Quantity:     2,
			PurchaseDate: time.Now(),
		},
		{
			Id:           4,
			ProductId:    2,
			Quantity:     10,
			PurchaseDate: time.Now(),
		},
	}
}

func TestAddInventory(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	inventoryModel := models.InventoryModel{
		Db: db,
	}
	mockInventory := entities.Inventory{
		Id:           1,
		ProductId:    1,
		Quantity:     1,
		PurchaseDate: time.Now(),
	}

	queryInsert := `INSERT INTO Inventory (ProductId,Quantity,PurchaseDate) VALUES($1,$2,$3))`
	rows := sqlmock.NewRows([]string{"id", "productId", "quantity", "purchaseId"})
	mock.ExpectQuery(regexp.QuoteMeta(queryInsert)).WithArgs(mockInventory.ProductId, mockInventory.Quantity, mockInventory.PurchaseDate).WillReturnRows(rows)

	err := inventoryModel.AddInventory(mockInventory)
	assert.Nil(t, err)
}

func TestAddInventoryError(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	inventoryModel := models.InventoryModel{
		Db: db,
	}
	mockInventory := entities.Inventory{
		Id:           1,
		ProductId:    1,
		Quantity:     1,
		PurchaseDate: time.Now(),
	}

	queryInsert := `INSERT INTO Inventory (ProductId,Quantity,PurchaseDate) VALUES($1,$2,$3))`
	mock.ExpectExec(regexp.QuoteMeta(queryInsert)).WithArgs(mockInventory.ProductId, mockInventory.Quantity, mockInventory.PurchaseDate).WillReturnError(fmt.Errorf(""))

	err := inventoryModel.AddInventory(mockInventory)
	assert.NotNil(t, err)
}

func TestGetInventory(t *testing.T) {
	testVendorId := 1
	db, mock := NewMock()
	defer db.Close()

	inventoryModel := models.InventoryModel{
		Db: db,
	}
	query := `select p.*, VendorId from Inventory i
	left join Product p on p.Id = i.ProductId
	where VendorId = $1`

	rows := sqlmock.NewRows([]string{"id", "productId", "quantity", "purchaseDate", "vendorId"})

	mockInventorys := CreateMockInventorys()
	for _, p := range mockInventorys {
		rows.AddRow(p.Id, p.ProductId, p.Quantity, p.PurchaseDate, testVendorId)
	}

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(testVendorId).WillReturnRows(rows)
	actualProducts, _ := inventoryModel.GetInventorysByVendorId(int64(testVendorId))
	assert.NotEmpty(t, actualProducts)
	assert.Equal(t, 2, len(actualProducts))
	for i, p := range mockInventorys {
		assert.Equal(t, p.Id, actualProducts[i].Id)
		assert.Equal(t, p.ProductId, actualProducts[i].ProductId)
		assert.Equal(t, p.Quantity, actualProducts[i].Quantity)
		assert.Equal(t, p.PurchaseDate, actualProducts[i].PurchaseDate)
	}
}

func TestGetInventoryError(t *testing.T) {
	testVendorId := 1
	db, mock := NewMock()
	defer db.Close()

	inventoryModel := models.InventoryModel{
		Db: db,
	}
	query := `select p.*, VendorId from Inventory i
	left join Product p on p.Id = i.ProductId
	where VendorId = $1`

	rows := sqlmock.NewRows([]string{"id", "productId", "quantity", "purchaseDate", "vendorId"})

	mockInventorys := CreateMockInventorys()
	for _, p := range mockInventorys {
		rows.AddRow(p.Id, p.ProductId, p.Quantity, p.PurchaseDate, testVendorId)
	}

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(testVendorId).WillReturnError(fmt.Errorf(""))
	actualInventorys, err := inventoryModel.GetInventorysByVendorId(int64(testVendorId))
	assert.Nil(t, actualInventorys)
	assert.Error(t, err)
}
