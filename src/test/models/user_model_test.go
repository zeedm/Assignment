package test

import (
	"api/assignment/src/entities"
	"api/assignment/src/models"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestAddUser(t *testing.T) {
	testEmail := "test@email.com"
	testPassword := "testPassword"
	testIsVendor := false
	db, mock := NewMock()
	defer db.Close()

	userModel := models.UserModel{
		Db: db,
	}
	queryCheckExisting := `select * from AppUser where Email = $1`
	rows := sqlmock.NewRows([]string{"id", "email", "password", "isVendor"})
	mock.ExpectQuery(regexp.QuoteMeta(queryCheckExisting)).WithArgs(testEmail).WillReturnRows(rows)

	queryInsert := `insert into AppUser values(?,?,?)`
	result := sqlmock.NewResult(1, 1)
	mock.ExpectExec(regexp.QuoteMeta(queryInsert)).WithArgs(testEmail, testPassword, testIsVendor).WillReturnResult(result)

	user, err := userModel.AddUser(testEmail, testPassword, testIsVendor)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testEmail, user.Email)
	assert.Equal(t, testIsVendor, user.IsVendor)
}

func TestAddExistingUser(t *testing.T) {
	testEmail := "test@email.com"
	testPassword := "testPassword"
	testIsVendor := true
	db, mock := NewMock()
	defer db.Close()

	userModel := models.UserModel{
		Db: db,
	}
	queryCheckExisting := `select * from AppUser where Email = $1`
	mock.ExpectQuery(regexp.QuoteMeta(queryCheckExisting)).WithArgs(testEmail).WillReturnError(fmt.Errorf("User already exists"))

	_, err := userModel.AddUser(testEmail, testPassword, testIsVendor)
	assert.Error(t, err)
}

func TestGetUserByEmailAndPassword(t *testing.T) {
	testEmail := "test@email.com"
	testPassword := "testPassword"
	db, mock := NewMock()

	defer db.Close()
	userModel := models.UserModel{
		Db: db,
	}
	query := `select * from AppUser where Email = $1 and password = $2`
	rows := sqlmock.NewRows([]string{"id", "email", "password", "isVendor"}).AddRow("1", testEmail, testPassword, false)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(testEmail, testPassword).WillReturnRows(rows)
	user, err := userModel.GetUserByEmailAndPassword(testEmail, testPassword)
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.NotNil(t, user)
	assert.Equal(t, testEmail, user.Email)
}

func TestGetUserByEmailAndPasswordNotFound(t *testing.T) {
	testEmail := "test@email.com"
	testPassword := "testPassword"
	db, mock := NewMock()

	userModel := models.UserModel{
		Db: db,
	}
	query := `select * from AppUser where Email = $1 and password = $2`
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(testEmail, testPassword).WillReturnError(fmt.Errorf(""))
	user, err := userModel.GetUserByEmailAndPassword(testEmail, testPassword)
	assert.Error(t, err)
	emptyUser := entities.User{}
	assert.Equal(t, emptyUser, user)
}
