package http_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"chiller/model"
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

func GetAuthServiceData() (u []model.User, e error) {

	auth, err := getBearerToken()
	if err != nil {
		e = err
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

	var users []model.User
	_ = json.Unmarshal(body, &users)

	return users, e
}
