package controllers

import (
	"api/assignment/src/config"
	"api/assignment/src/entities"
	"api/assignment/src/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func GetInventorys(responseWriter http.ResponseWriter, request *http.Request) {
	if !IsAuthorized(request, true) {
		json.NewEncoder(responseWriter).Encode("401 Unauthorized")
		return
	}
	db, errorMessage := config.GetDB()
	defer config.CloseDB(db)
	if errorMessage != nil {
		fmt.Println(errorMessage)
	} else {
		inventoryModel := models.InventoryModel{
			Db: db,
		}
		intId, errorParseInt := strconv.Atoi(request.Header.Get("UserId"))
		if errorParseInt != nil {
			fmt.Println(errorParseInt)
		}
		inventorys, errorMessageInventory := inventoryModel.GetInventorysByVendorId(int64(intId))
		if errorMessageInventory != nil {
			fmt.Println(errorMessageInventory)
		} else {
			json.NewEncoder(responseWriter).Encode(inventorys)
		}
	}
}

func AddInventory(responseWriter http.ResponseWriter, request *http.Request) {
	if !IsAuthorized(request, true) {
		json.NewEncoder(responseWriter).Encode("401 Unauthorized")
		return
	}
	db, errorMessage := config.GetDB()
	defer config.CloseDB(db)
	if errorMessage != nil {
		fmt.Println(errorMessage)
	} else {
		inventoryModel := models.InventoryModel{
			Db: db,
		}
		var inventory entities.Inventory
		errorMessageDecode := json.NewDecoder(request.Body).Decode(&inventory)
		if errorMessageDecode != nil {
			fmt.Println(errorMessageDecode)
			return
		}
		errorMessageInventory := inventoryModel.AddInventory(inventory)
		if errorMessageInventory != nil {
			fmt.Println(errorMessageInventory)
		} else {
			json.NewEncoder(responseWriter).Encode("Add Inventory successfully")
		}
	}
}
