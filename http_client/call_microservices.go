package http_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	username        = os.Getenv("CHILLER_LOGIN_USERNAME")
	password        = os.Getenv("CHILLER_LOGIN_PASSWORD")
	auth_url        = os.Getenv("CHILLER_AUTH_URL")
	backup_auth_url = os.Getenv("CHILLER_BACKUP_AUTH_URL")
)

type creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type bearer struct {
	Username     string   `json:"username"`
	Roles        []string `json:"roles"`
	Access_token string   `json:"access_token"`
	Token_type   string   `json:"token_type"`
	Expires_in   int      `json:"expires_in"`
}

func getBearerToken() (brr bearer, e error) {

	login := creds{
		Username: username,
		Password: password,
	}

	rb, _ := json.Marshal(login)

	res, err := http.Post(auth_url, "application/json", bytes.NewBuffer(rb))
	if err != nil {
		e = err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {

		b := bearer{}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			e = err
		}

		authErr := json.Unmarshal(body, &b)
		if authErr != nil {
			e = authErr
		}
		brr = b
	}

	return brr, e
}

type User struct {
	Id             int64           `json:"id"`
	Username       string          `json:"username"`
	Password       string          `json:"password"`
	Firstname      string          `json:"firstname"`
	Lastname       string          `json:"lastname"`
	DateCreated    []int           `json:"dateCreated"`
	Enabled        bool            `json:"enabled"`
	AccountExpired bool            `json:"accountExpired"`
	AccountLocked  bool            `json:"accountLocked"`
	Birthday       []int           `json:"birthday"`
	UserRoles      []UserRoles     `json:"userRoles"`
	UserAddresses  []UserAddresses `json:"userAddresses"`
	UserPhones     []UserPhones    `json:"userPhones"`
}

type UserRoles struct {
	Id   int64 `json:"id"`
	Role struct {
		Id          int64  `json:"id"`
		Role        string `json:"role"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
}

type UserAddresses struct {
	Id      int64 `json:"id"`
	Address struct {
		Id      int64  `json:"id"`
		Address string `json:"address"`
		City    string `json:"city"`
		State   string `json:"state"`
		Zip     string `json:"zip"`
	}
}

type UserPhones struct {
	Id    int64 `json:"id"`
	Phone struct {
		Id    int64  `json:"id"`
		Phone string `json:"phone"`
		Type  string `json:"type"`
	}
}

func GetAuthServiceData() (u []User, e error) {

	auth, err := getBearerToken()
	if err != nil {
		panic(err)
	}
	bearer := fmt.Sprintf("Bearer %s", auth.Access_token)

	req, err := http.NewRequest("GET", backup_auth_url, nil)
	if err != nil {
		e = err
	}
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		e = err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		e = err
	}

	var users []User
	_ = json.Unmarshal(body, &users)

	return users, e
}
