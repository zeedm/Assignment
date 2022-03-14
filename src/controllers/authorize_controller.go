package controllers

import (
	"net/http"
	"strconv"
)

func IsAuthorized(request *http.Request, isVendor bool) bool {
	return request.Header.Get("IsVendor") == strconv.FormatBool(isVendor)
}
