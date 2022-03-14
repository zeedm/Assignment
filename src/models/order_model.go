package models

import (
	"api/assignment/src/entities"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

type OrderModel struct {
	Db *sql.DB
}

func (orderModel OrderModel) CreateOrder(cartInSession []entities.ProductInCart, userId int) error {
	errorCheckInventory := orderModel.CheckQuantityInInventory(cartInSession)
	if errorCheckInventory != nil {
		return errorCheckInventory
	}
	tx, err := orderModel.Db.Begin()
	if err != nil {
		return err
	}

	errorReduceQuantityInInventory := orderModel.ReduceQuantityInInventory(tx, cartInSession)
	if errorReduceQuantityInInventory != nil {
		return errorReduceQuantityInInventory
	}

	errorInsertOrder := orderModel.InsertOrder(tx, cartInSession, userId)
	if errorInsertOrder != nil {
		return errorInsertOrder
	}

	errorCommit := tx.Commit()
	if errorCommit != nil {
		return errorCommit
	}
	return nil
}
func (orderModel OrderModel) CheckQuantityInInventory(cartInSession []entities.ProductInCart) error {
	stringIds := ConvertProductToStringId(&cartInSession)
	for _, a := range stringIds {
		fmt.Println(a)
	}
	rows, errorMessage := orderModel.Db.Query(`select ProductId, sum(quantity) as quantity from Inventory where ProductId in (` + strings.Join(stringIds, ",") + `) group by ProductId`)
	if errorMessage != nil {
		return errorMessage
	}
	for rows.Next() {
		var id int64
		var quantity int64
		errorMessageInventory := rows.Scan(&id, &quantity)
		if errorMessageInventory != nil {
			return errorMessageInventory
		}
		for _, p := range cartInSession {
			if p.Id == id {
				if p.QuantityInCart > quantity {
					return errors.New("Not enough item in Inventory")
				}
			}
		}
	}
	return nil
}

func (orderModel OrderModel) ReduceQuantityInInventory(tx *sql.Tx, cart []entities.ProductInCart) error {
	listOfInventorys, errorGetAllInventory := orderModel.GetAll(cart)
	if errorGetAllInventory != nil {
		return errorGetAllInventory
	}
	sort.SliceStable(listOfInventorys, func(i, j int) bool {
		return listOfInventorys[i].PurchaseDate.Before(listOfInventorys[j].PurchaseDate)
	})
	for _, productInCart := range cart {
		tempQuantity := productInCart.QuantityInCart
		for i2, inventory := range listOfInventorys {
			if inventory.ProductId == productInCart.Id && tempQuantity > 0 {
				if inventory.Quantity >= tempQuantity {
					listOfInventorys[i2].Quantity -= tempQuantity
					tempQuantity = 0
				} else {
					tempQuantity = tempQuantity - listOfInventorys[i2].Quantity
					listOfInventorys[i2].Quantity = 0
				}
			}
			if tempQuantity == 0 {
				break
			}
		}
	}
	statementUpdateQuantity, err := tx.Prepare(`update Inventory set Quantity = $1 where Id = $2`)
	if err != nil {
		return err
	}
	defer statementUpdateQuantity.Close()

	for _, inventory := range listOfInventorys {
		if _, err := statementUpdateQuantity.Exec(inventory.Quantity, inventory.Id); err != nil {
			return err
		}
	}
	return nil
}
func (orderModel OrderModel) GetAll(cart []entities.ProductInCart) ([]entities.Inventory, error) {
	var listOfInventorys []entities.Inventory
	stringIds := ConvertProductToStringId(&cart)
	rows, errorMessage := orderModel.Db.Query(`select * from Inventory where productId in (` + strings.Join(stringIds, ",") + `)`)
	if errorMessage != nil {
		return nil, errorMessage
	}
	for rows.Next() {
		var id int64
		var productId int64
		var quantity int64
		var purchaseDate time.Time
		errorMessageInventory := rows.Scan(&id, &productId, &quantity, &purchaseDate)
		if errorMessageInventory != nil {
			return nil, errorMessageInventory
		}
		inventory := entities.Inventory{
			Id:           id,
			ProductId:    productId,
			Quantity:     quantity,
			PurchaseDate: purchaseDate,
		}
		listOfInventorys = append(listOfInventorys, inventory)
	}
	return listOfInventorys, nil
}

func (orderModel OrderModel) InsertOrder(tx *sql.Tx, cart []entities.ProductInCart, userId int) error {
	var lastInsertId int64
	err := orderModel.Db.QueryRow(`INSERT INTO UserOrder (BuyerId) VALUES($1); select ID = convert(bigint, SCOPE_IDENTITY())`, userId).Scan(&lastInsertId)
	if err != nil {
		return err
	}
	stateMentInsertOrder_Product, err := tx.Prepare(`INSERT INTO Order_Product (OrderId, ProductId, Quantity, OrderDate) VALUES($1, $2, $3, $4)`)
	if err != nil {
		return err
	}
	defer stateMentInsertOrder_Product.Close()

	for _, product := range cart {
		if err != nil {
			return err
		}
		if _, err := stateMentInsertOrder_Product.Exec(lastInsertId, product.Id, product.QuantityInCart, time.Now()); err != nil {
			return err
		}
	}
	return nil
}

func ConvertProductToStringId(cartInSession *[]entities.ProductInCart) []string {
	var tmp []string

	for _, v := range *cartInSession {
		tmp = append(tmp, fmt.Sprint(v.Id))
	}
	return tmp
}
