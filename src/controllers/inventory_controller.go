package controllers

import (
	"api/assignment/src/config"
	"api/assignment/src/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func GetInventorys(responseWriter http.ResponseWriter, request *http.Request) {
	if !IsAuthorized(request, false) {
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
