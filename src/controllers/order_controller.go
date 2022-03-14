package controllers

import (
	"api/assignment/src/config"
	"api/assignment/src/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func Checkout(responseWriter http.ResponseWriter, request *http.Request) {
	if !IsAuthorized(request, false) {
		json.NewEncoder(responseWriter).Encode("401 Unauthorized")
		return
	}
	db, errorMessage := config.GetDB()
	defer config.CloseDB(db)
	if errorMessage != nil {
		fmt.Println(errorMessage)
		return
	}

	cartInSession, errorMessageSession := GetCartInSession(request)
	if errorMessageSession != nil {
		fmt.Println(errorMessageSession)
		return
	}

	userId, errorParseInt := strconv.Atoi(request.Header.Get("UserId"))
	if errorParseInt != nil {
		fmt.Println(errorParseInt)
	}
	orderModel := models.OrderModel{
		Db: db,
	}
	errorMessageCreateOrder := orderModel.CreateOrder(cartInSession, userId)
	if errorMessageCreateOrder != nil {
		fmt.Println(errorMessageCreateOrder)
		json.NewEncoder(responseWriter).Encode("Unable to checkout")
		return
	}
	json.NewEncoder(responseWriter).Encode("Checkout successfully")
}
