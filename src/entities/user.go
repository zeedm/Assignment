package entities

type User struct {
	Id       int64
	Email    string
	Password string
	IsVendor bool
}
