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
	Birthday:       "",
}

func TestInsertUser(t *testing.T) {

	t.Log(InsertUser(u))
}

func TestUpdateUser(t *testing.T) {

	t.Log(UpdateUser(u))
}

func TestFindAllUsers(t *testing.T) {

	users, _ := FindAllUsers()

	for _, v := range users {
		t.Logf("dob: %v", v)
	}
	if len(users) < 1 {
		t.Log("No users returned.")
		t.Fail()
	}
}
