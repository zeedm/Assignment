package models

import (
	"api/assignment/src/entities"
	"database/sql"
	"errors"
)

type UserModel struct {
	Db *sql.DB
}

func (userModel UserModel) GetUserByEmailAndPassword(email string, password string) (entities.User, error) {
	var id int64
	var isVendor bool
	var passwordDump string
	errorMessage := userModel.Db.QueryRow(`select * from AppUser where Email = $1 and password = $2`, email, password).Scan(&id, &email, &passwordDump, &isVendor)
	if errorMessage != nil {
		return entities.User{}, errorMessage
	} else {
		user := entities.User{
			Id:       id,
			Email:    email,
			IsVendor: isVendor,
		}
		return user, nil
	}
}

func (userModel UserModel) AddUser(email, password string, isVendor bool) (entities.User, error) {
	errorMessage := userModel.Db.QueryRow(`select * from AppUser where Email = $1`, email).Scan()
	if errorMessage != nil {
		result, errorMessageInsert := userModel.Db.Exec(`insert into AppUser values(?,?,?)`,
			email, password, isVendor)
		if errorMessageInsert != nil {
			return entities.User{}, errorMessageInsert
		} else {
			userId, _ := result.LastInsertId()
			user := entities.User{
				Id:       userId,
				Email:    email,
				IsVendor: isVendor,
			}
			return user, nil
		}
	} else {
		errorResult := errors.New("User already exists")
		return entities.User{}, errorResult
	}
}
