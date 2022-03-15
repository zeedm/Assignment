package main

import (
	"api/assignment/src/controllers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func InitializeRoute() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/login", controllers.Login).Methods(http.MethodPost)
	myRouter.HandleFunc("/signup", controllers.Signup).Methods(http.MethodPost)
	subRouter := myRouter.Methods(http.MethodGet, http.MethodPost).Subrouter()
	subRouter.Use(IsAuthenticated)
	subRouter.HandleFunc("/products", controllers.GetProducts).Methods(http.MethodGet)
	subRouter.HandleFunc("/products/addProductToCart", controllers.AddProductToSessionCart).Methods(http.MethodPost)
	subRouter.HandleFunc("/products/viewCart", controllers.ViewCart).Methods(http.MethodGet)
	subRouter.HandleFunc("/checkout", controllers.Checkout).Methods(http.MethodGet)
	subRouter.HandleFunc("/inventorys", controllers.GetInventorys).Methods(http.MethodGet)
	subRouter.HandleFunc("/inventorys/addInventory", controllers.AddInventory).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	InitializeRoute()
}
