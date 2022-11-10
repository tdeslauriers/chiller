package dao

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var u User = User{
	Id:             42,
	Username:       "tom@tom.com",
	Password:       "123456789",
	Firstname:      "Tom",
	Lastname:       "Atomic",
	DateCreated:    []int{2000, 5, 5},
	Enabled:        true,
	AccountExpired: false,
	AccountLocked:  false,
	Birthday:       []int{1971, 8, 8},
}

func TestInsertUser(t *testing.T) {

	t.Log(insertUser(u))
}

func TestUpdateUser(t *testing.T) {

	t.Log(updateUser(u))
}
