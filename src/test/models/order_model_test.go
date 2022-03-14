package test

import (
	"api/assignment/src/entities"
	"api/assignment/src/models"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func CreateMockProductInCarts() []entities.ProductInCart {
	return []entities.ProductInCart{
		{
			Id:             1,
			QuantityInCart: 2,
		},
		{
			Id:             2,
			QuantityInCart: 3,
		},
	}
}

func TestInsertOrder(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	orderModel := models.OrderModel{
		Db: db,
	}
	mockProductInCarts := CreateMockProductInCarts()
	mock.MatchExpectationsInOrder(false)
	testOrderId := 1
	testBuyerId := 1
	testPurchaseDate := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

	monkey.Patch(time.Now, func() time.Time {
		return testPurchaseDate
	})
	mock.ExpectBegin()
	tx, _ := db.Begin()
	rows := sqlmock.NewRows([]string{"id"}).AddRow(testOrderId)
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO UserOrder (BuyerId) VALUES($1); select ID = convert(bigint, SCOPE_IDENTITY())`)).
		WillReturnRows(rows)

	for _, p := range mockProductInCarts {
		mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO Order_Product (OrderId, ProductId, Quantity, OrderDate) VALUES($1, $2, $3, $4)`)).
			ExpectExec().
			WithArgs(testOrderId, p.Id, p.QuantityInCart, testPurchaseDate).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	err := orderModel.InsertOrder(tx, mockProductInCarts, testBuyerId)
	assert.Nil(t, err)
}

func TestInsertOrderError(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	orderModel := models.OrderModel{
		Db: db,
	}
	mockProductInCarts := CreateMockProductInCarts()
	testOrderId := 1
	testBuyerId := 1
	mock.ExpectBegin()
	tx, _ := db.Begin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO UserOrder (BuyerId) VALUES($1); select ID = convert(bigint, SCOPE_IDENTITY())`)).
		WillReturnError(fmt.Errorf(""))

	for _, p := range mockProductInCarts {
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO Order_Product (OrderId, ProductId, Quantity, OrderDate) VALUES($1, $2, $3, $4)`)).
			WithArgs(testOrderId, p.Id, p.QuantityInCart, time.Now).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	err := orderModel.InsertOrder(tx, mockProductInCarts, testBuyerId)
	assert.Error(t, err)
}

func TestInsertOrderErrorInsertOrder_Product(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	orderModel := models.OrderModel{
		Db: db,
	}
	mockProductInCarts := CreateMockProductInCarts()
	mock.MatchExpectationsInOrder(false)
	testOrderId := 1
	testBuyerId := 1
	testPurchaseDate := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

	monkey.Patch(time.Now, func() time.Time {
		return testPurchaseDate
	})
	mock.ExpectBegin()
	tx, _ := db.Begin()
	rows := sqlmock.NewRows([]string{"id"}).AddRow(testOrderId)
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO UserOrder (BuyerId) VALUES($1); select ID = convert(bigint, SCOPE_IDENTITY())`)).
		WillReturnRows(rows)

	for _, p := range mockProductInCarts {
		mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO Order_Product (OrderId, ProductId, Quantity, OrderDate) VALUES($1, $2, $3, $4)`)).
			ExpectExec().
			WithArgs(testOrderId, p.Id, p.QuantityInCart, testPurchaseDate).
			WillReturnError(fmt.Errorf(""))
	}

	err := orderModel.InsertOrder(tx, mockProductInCarts, testBuyerId)
	assert.Error(t, err)
}

func TestCheckQuantityInInventory(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()
	mockProductInCarts := CreateMockProductInCarts()
	stringIds := models.ConvertProductToStringId(&mockProductInCarts)
	orderModel := models.OrderModel{
		Db: db,
	}
	rows := sqlmock.NewRows([]string{"id", "quantityInCart"})
	for _, p := range mockProductInCarts {
		rows.AddRow(p.Id, 100)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`select ProductId, sum(quantity) as quantity from Inventory where ProductId in (` + strings.Join(stringIds, ",") + `) group by ProductId`)).
		WillReturnRows(rows)
	err := orderModel.CheckQuantityInInventory(mockProductInCarts)
	assert.Nil(t, err)
}

func TestCheckQuantityInInventoryGetInventoryErrorMessage(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()
	mockProductInCarts := CreateMockProductInCarts()
	stringIds := models.ConvertProductToStringId(&mockProductInCarts)
	orderModel := models.OrderModel{
		Db: db,
	}
	rows := sqlmock.NewRows([]string{"id", "quantityInCart"})
	for _, p := range mockProductInCarts {
		rows.AddRow(p.Id, 1)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`select ProductId, sum(quantity) as quantity from Inventory where ProductId in (` + strings.Join(stringIds, ",") + `) group by ProductId`)).
		WillReturnError(fmt.Errorf(""))
	err := orderModel.CheckQuantityInInventory(mockProductInCarts)
	assert.Error(t, err)
}

func TestCheckQuantityInInventoryNotEnoughInInventory(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()
	mockProductInCarts := CreateMockProductInCarts()
	stringIds := models.ConvertProductToStringId(&mockProductInCarts)
	orderModel := models.OrderModel{
		Db: db,
	}
	rows := sqlmock.NewRows([]string{"id", "quantityInCart"})
	for _, p := range mockProductInCarts {
		rows.AddRow(p.Id, 1)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`select ProductId, sum(quantity) as quantity from Inventory where ProductId in (` + strings.Join(stringIds, ",") + `) group by ProductId`)).
		WillReturnRows(rows)
	err := orderModel.CheckQuantityInInventory(mockProductInCarts)
	assert.EqualError(t, err, "Not enough item in Inventory")
}

func TestCheckQuantityInInventoryNotEnoughInInventorySecondItem(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()
	mockProductInCarts := CreateMockProductInCarts()
	stringIds := models.ConvertProductToStringId(&mockProductInCarts)
	orderModel := models.OrderModel{
		Db: db,
	}
	rows := sqlmock.NewRows([]string{"id", "quantityInCart"})
	for i, p := range mockProductInCarts {
		if i == 0 {
			rows.AddRow(p.Id, 100)
		} else {
			rows.AddRow(p.Id, 1)
		}
	}

	mock.ExpectQuery(regexp.QuoteMeta(`select ProductId, sum(quantity) as quantity from Inventory where ProductId in (` + strings.Join(stringIds, ",") + `) group by ProductId`)).
		WillReturnRows(rows)
	err := orderModel.CheckQuantityInInventory(mockProductInCarts)
	assert.EqualError(t, err, "Not enough item in Inventory")
}

func CreateMockReducedInventorys() []entities.Inventory {
	return []entities.Inventory{
		{
			Id:           1,
			ProductId:    1,
			Quantity:     3,
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
			Quantity:     0,
			PurchaseDate: time.Now(),
		},
		{
			Id:           4,
			ProductId:    2,
			Quantity:     9,
			PurchaseDate: time.Now(),
		},
	}
}

func TestReduceQuantityInInventory(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	orderModel := models.OrderModel{
		Db: db,
	}
	mockProductInCarts := CreateMockProductInCarts()
	mockInventorys := CreateMockInventorys()
	stringIds := models.ConvertProductToStringId(&mockProductInCarts)
	mock.MatchExpectationsInOrder(false)
	mock.ExpectBegin()
	tx, _ := db.Begin()
	inventoryRows := sqlmock.NewRows([]string{"id", "productId", "quantity", "purchaseDate"})
	for _, inv := range mockInventorys {
		inventoryRows.AddRow(inv.Id, inv.ProductId, inv.Quantity, inv.PurchaseDate)
	}
	mock.ExpectQuery(regexp.QuoteMeta(`select * from Inventory where productId in (` + strings.Join(stringIds, ",") + `)`)).
		WillReturnRows(inventoryRows)
	mockReducedInventorys := CreateMockReducedInventorys()
	for _, inv := range mockReducedInventorys {
		mock.ExpectPrepare(regexp.QuoteMeta(`update Inventory set Quantity = $1 where Id = $2`)).
			ExpectExec().
			WithArgs(inv.Quantity, inv.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	err := orderModel.ReduceQuantityInInventory(tx, mockProductInCarts)
	assert.Nil(t, err)
}

func TestReduceQuantityInInventoryErrorGetInventory(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	orderModel := models.OrderModel{
		Db: db,
	}
	mockProductInCarts := CreateMockProductInCarts()
	stringIds := models.ConvertProductToStringId(&mockProductInCarts)
	mock.MatchExpectationsInOrder(false)
	mock.ExpectBegin()
	tx, _ := db.Begin()
	mock.ExpectQuery(regexp.QuoteMeta(`select * from Inventory where productId in (` + strings.Join(stringIds, ",") + `)`)).
		WillReturnError(fmt.Errorf("Select Error"))
	mockReducedInventorys := CreateMockReducedInventorys()
	for _, inv := range mockReducedInventorys {
		mock.ExpectPrepare(regexp.QuoteMeta(`update Inventory set Quantity = $1 where Id = $2`)).
			ExpectExec().
			WithArgs(inv.Quantity, inv.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	err := orderModel.ReduceQuantityInInventory(tx, mockProductInCarts)
	assert.EqualError(t, err, "Select Error")
}

func TestReduceQuantityInInventoryErrorUpdateInventory(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	orderModel := models.OrderModel{
		Db: db,
	}
	mockProductInCarts := CreateMockProductInCarts()
	mockInventorys := CreateMockInventorys()
	stringIds := models.ConvertProductToStringId(&mockProductInCarts)
	mock.MatchExpectationsInOrder(false)
	mock.ExpectBegin()
	tx, _ := db.Begin()
	inventoryRows := sqlmock.NewRows([]string{"id", "productId", "quantity", "purchaseDate"})
	for _, inv := range mockInventorys {
		inventoryRows.AddRow(inv.Id, inv.ProductId, inv.Quantity, inv.PurchaseDate)
	}
	mock.ExpectQuery(regexp.QuoteMeta(`select * from Inventory where productId in (` + strings.Join(stringIds, ",") + `)`)).
		WillReturnRows(inventoryRows)
	mockReducedInventorys := CreateMockReducedInventorys()
	for _, inv := range mockReducedInventorys {
		mock.ExpectPrepare(regexp.QuoteMeta(`update Inventory set Quantity = $1 where Id = $2`)).
			ExpectExec().
			WithArgs(inv.Quantity, inv.Id).
			WillReturnError(fmt.Errorf("Update Error"))
	}

	err := orderModel.ReduceQuantityInInventory(tx, mockProductInCarts)
	assert.EqualError(t, err, "Update Error")
}
