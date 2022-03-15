package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(email, isVendor string, id int64) (string, error) {
	var mySigningKey = []byte("test")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["userId"] = id
	claims["authorized"] = true
	claims["email"] = email
	claims["isVendor"] = isVendor
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}
