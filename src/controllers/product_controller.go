package controllers

import (
	"api/assignment/src/config"
	"api/assignment/src/entities"
	"api/assignment/src/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
)

func GetProducts(responseWriter http.ResponseWriter, request *http.Request) {
	db, errorMessage := config.GetDB()
	defer config.CloseDB(db)
	if errorMessage != nil {
		fmt.Println(errorMessage)
	} else {
		productModel := models.ProductModel{
			Db: db,
		}
		if IsAuthorized(request, false) {
			products, errorMessageProduct := productModel.GetAllProducts()
			if errorMessageProduct != nil {
				fmt.Println(errorMessageProduct)
			} else {
				json.NewEncoder(responseWriter).Encode(products)
			}
		} else if IsAuthorized(request, true) {
			intId, errorParseInt := strconv.Atoi(request.Header.Get("UserId"))
			if errorParseInt != nil {
				fmt.Println(errorParseInt)
			}
			products, errorMessageProduct := productModel.FindProductByVendorId(intId)
			if errorMessageProduct != nil {
				fmt.Println(errorMessageProduct)
			} else {
				json.NewEncoder(responseWriter).Encode(products)
			}
		}
	}
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func AddProductToSessionCart(responseWriter http.ResponseWriter, request *http.Request) {
	if !IsAuthorized(request, false) {
		json.NewEncoder(responseWriter).Encode("401 Unauthorized")
		return
	}
	cartInSession, errorMessageSession := GetCartInSession(request)
	if errorMessageSession != nil {
		fmt.Println(errorMessageSession)
		return
	}
	session, errorAddToCart := SetCartInSession(request, &cartInSession)
	if errorAddToCart != nil {
		fmt.Println(errorAddToCart)
		return
	}
	errorSetCart := session.Save(request, responseWriter)
	if errorSetCart != nil {
		json.NewEncoder(responseWriter).Encode("Unable to add product to cart")
	} else {
		json.NewEncoder(responseWriter).Encode("Add product to cart successfully")
	}
}

func ViewCart(responseWriter http.ResponseWriter, request *http.Request) {
	if !IsAuthorized(request, false) {
		json.NewEncoder(responseWriter).Encode("401 Unauthorized")
		return
	}

	cartInSession, _ := GetCartInSession(request)
	json.NewEncoder(responseWriter).Encode(cartInSession)
}

func GetCartInSession(request *http.Request) ([]entities.ProductInCart, error) {
	session, _ := store.Get(request, request.Header["Token"][0])
	var cartInSession []entities.ProductInCart

	json.Unmarshal([]byte(fmt.Sprintf("%s", session.Values["listOfProducts"])), &cartInSession)
	return cartInSession, nil
}

func SetCartInSession(request *http.Request, cartInSession *[]entities.ProductInCart) (sessions.Session, error) {
	session, _ := store.Get(request, request.Header["Token"][0])
	var productInCart entities.ProductInCart
	errorMessageDecode := json.NewDecoder(request.Body).Decode(&productInCart)
	if errorMessageDecode != nil {
		fmt.Println(errorMessageDecode)
		return sessions.Session{}, errorMessageDecode
	} else {
		isContain := false
		for i, p := range *cartInSession {
			if p.Id == productInCart.Id {
				isContain = true
				(*cartInSession)[i].QuantityInCart += productInCart.QuantityInCart
				break
			}
		}
		if isContain == false {
			*cartInSession = append(*cartInSession, productInCart)
		}
	}
	jsonCartInSession, errorJsonCartInSession := json.Marshal(cartInSession)
	if errorJsonCartInSession != nil {
		fmt.Println(errorJsonCartInSession)
		return sessions.Session{}, errorJsonCartInSession
	} else {
		session.Values["listOfProducts"] = jsonCartInSession
		return *session, nil
	}
}
