package http_client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	username = os.Getenv("CHILLER_LOGIN_USERNAME")
	password = os.Getenv("CHILLER_LOGIN_PASSWORD")
	url      = os.Getenv("CHILLER_AUTH_URL")
)

type creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type bearer struct {
	Username     string `json:"username"`
	Roles        []string
	Access_token string `json:"access_token"`
	Token_type   string `json:"token_type"`
	Expires_in   int    `json:"expires_in"`
}

func getBearerToken() (brr bearer, e error) {

	login := creds{
		Username: username,
		Password: password,
	}

	rb, _ := json.Marshal(login)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(rb))
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
