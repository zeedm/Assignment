package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

func IsAuthenticated(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			json.NewEncoder(w).Encode("No Token Found")
			return
		}

		var mySigningKey = []byte("test")

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return mySigningKey, nil
		})

		if err != nil {
			json.NewEncoder(w).Encode(err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["userId"].(float64)
			r.Header.Set("UserId", strconv.Itoa(int(id)))
			if claims["isVendor"] == "true" {

				r.Header.Set("IsVendor", "true")
				handler.ServeHTTP(w, r)

			} else if claims["isVendor"] == "false" {

				r.Header.Set("IsVendor", "false")
				handler.ServeHTTP(w, r)
			}
		} else {
			json.NewEncoder(w).Encode("401  Unauthorized")
		}
	})
}
