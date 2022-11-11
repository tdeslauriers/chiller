package dao

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var u User = User{
	Id:             44,
	Username:       "tom@tom.com",
	Password:       "123456789",
	Firstname:      "Tom",
	Lastname:       "Atomic",
	DateCreated:    "2007-02-02",
	Enabled:        false,
	AccountExpired: false,
	AccountLocked:  false,
	Birthday:       "1969-12-08",
}

func TestInsertUser(t *testing.T) {

	t.Log(insertUser(u))
}

func TestUpdateUser(t *testing.T) {

	t.Log(updateUser(u))
}

func TestFindAllUsers(t *testing.T) {

	users, _ := findAllUsers()

	for _, v := range users {
		t.Logf("dob: %v", v)
	}
	if len(users) < 1 {
		t.Log("No users returned.")
		t.Fail()
	}
}
