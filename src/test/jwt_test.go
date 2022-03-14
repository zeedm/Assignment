package test

import (
	"api/assignment/src/utils"
	"fmt"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

func TestGenerateJwt(t *testing.T) {
	testEmail := "test@email.com"
	testIsVendor := "true"
	testId := int64(1)
	actualToken, err := utils.GenerateJWT(testEmail, testIsVendor, testId)
	if err != nil {
		t.Errorf(err.Error())
	}

	token, errorParseToken := jwt.Parse(actualToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return []byte("test"), nil
	})
	if errorParseToken != nil {
		t.Errorf(errorParseToken.Error())
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		actualId := claims["userId"].(float64)
		actualEmail := claims["email"]
		actualIsVendor := claims["isVendor"]
		if testEmail != actualEmail {
			t.Errorf("Incorrect Email")
		} else if testIsVendor != actualIsVendor {
			t.Errorf("Incorrect IsVendor")
		} else if testId != int64(actualId) {
			t.Errorf("Incorrect Id")
		}
	} else {
		t.Errorf("401  Unauthorized")
	}
}
