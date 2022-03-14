package models

import (
	"api/assignment/src/entities"
	"database/sql"
	"time"
)

type InventoryModel struct {
	Db *sql.DB
}

func (inventoryModel InventoryModel) GetInventorysByVendorId(vendorId int64) ([]entities.Inventory, error) {
	rows, errorMessage := inventoryModel.Db.Query(`select p.*, VendorId from Inventory i
	left join Product p on p.Id = i.ProductId
	where VendorId = $1`, vendorId)
	if errorMessage != nil {
		return nil, errorMessage
	} else {
		var inventorys []entities.Inventory
		for rows.Next() {
			var id int64
			var productId int64
			var quantity int64
			var dumpVendorId int64
			var purchaseDate time.Time
			errorMessageInventory := rows.Scan(&id, &productId, &quantity, &purchaseDate, &dumpVendorId)
			if errorMessageInventory != nil {
				return nil, errorMessageInventory
			} else {
				inventory := entities.Inventory{
					Id:           id,
					ProductId:    productId,
					Quantity:     quantity,
					PurchaseDate: purchaseDate,
				}
				inventorys = append(inventorys, inventory)
			}
		}
		return inventorys, nil
	}
}

func (inventoryModel InventoryModel) AddInventory(inventory entities.Inventory) error {
	_, err := inventoryModel.Db.Query(`INSERT INTO Inventory (ProductId,Quantity,PurchaseDate) VALUES($1,$2,$3))`, inventory.ProductId, inventory.Quantity, inventory.PurchaseDate)
	if err != nil {
		return err
	}
	return nil
}
