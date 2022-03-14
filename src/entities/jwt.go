package entities

type UserForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	IsVendor bool   `json:"isVendor"`
}

type Token struct {
	IsVendor    bool   `json:"isVendor"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}
