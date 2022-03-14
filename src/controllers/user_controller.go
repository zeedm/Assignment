package controllers

import (
	"api/assignment/src/config"
	"api/assignment/src/entities"
	"api/assignment/src/models"
	"api/assignment/src/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func Login(responseWriter http.ResponseWriter, request *http.Request) {
	db, errorMessageDB := config.GetDB()
	defer config.CloseDB(db)
	if errorMessageDB != nil {
		json.NewEncoder(responseWriter).Encode(errorMessageDB)
	} else {
		var authdetails entities.UserForm
		errorMessageDecode := json.NewDecoder(request.Body).Decode(&authdetails)
		if errorMessageDecode != nil {
			json.NewEncoder(responseWriter).Encode(errorMessageDecode)
		}
		userModel := models.UserModel{
			Db: db,
		}
		authUser, errorMessageProduct := userModel.GetUserByEmailAndPassword(authdetails.Email, authdetails.Password)
		if errorMessageProduct != nil {
			json.NewEncoder(responseWriter).Encode(errorMessageProduct)
		} else {
			if authUser == (entities.User{}) {
				responseWriter.Header().Set("Content-Type", "application/json")
				json.NewEncoder(responseWriter).Encode("Cannot find user")
				return
			} else {
				validToken, errorMessageAuthUser := utils.GenerateJWT(authUser.Email, fmt.Sprintf("%t", authUser.IsVendor), authUser.Id)
				if errorMessageAuthUser != nil {
					json.NewEncoder(responseWriter).Encode("Unable to generate JWT")
				}
				var token entities.Token
				token.Email = authUser.Email
				token.IsVendor = authUser.IsVendor
				token.TokenString = validToken
				responseWriter.Header().Set("Content-Type", "application/json")
				json.NewEncoder(responseWriter).Encode(token)
				return
			}
		}
	}
	http.Error(responseWriter, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func Signup(responseWriter http.ResponseWriter, request *http.Request) {
	db, errorMessageDB := config.GetDB()
	defer config.CloseDB(db)
	if errorMessageDB != nil {
		fmt.Println(errorMessageDB)
	} else {
		var authdetails entities.UserForm
		errorMessageDecode := json.NewDecoder(request.Body).Decode(&authdetails)
		if errorMessageDecode != nil {
			fmt.Println(errorMessageDecode)
		}
		userModel := models.UserModel{
			Db: db,
		}
		authUser, errorMessageProduct := userModel.AddUser(authdetails.Email, authdetails.Password, authdetails.IsVendor)
		if errorMessageProduct != nil {
			json.NewEncoder(responseWriter).Encode(errorMessageProduct)
		} else {
			if authUser == (entities.User{}) {
				responseWriter.Header().Set("Content-Type", "application/json")
				json.NewEncoder(responseWriter).Encode("Cannot create user")
				return
			} else {
				validToken, errorMessageAuthUser := utils.GenerateJWT(authUser.Email, fmt.Sprintf("%t", authUser.IsVendor), authUser.Id)
				if errorMessageAuthUser != nil {
					json.NewEncoder(responseWriter).Encode("Unable to generate JWT")
				}
				var token entities.Token
				token.Email = authUser.Email
				token.IsVendor = authUser.IsVendor
				token.TokenString = validToken
				responseWriter.Header().Set("Content-Type", "application/json")
				json.NewEncoder(responseWriter).Encode(token)
				return
			}
		}
	}
	http.Error(responseWriter, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
