package http_client

import (
	"bytes"
	"chiller/dao"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	username           = os.Getenv("CHILLER_LOGIN_USERNAME")
	password           = os.Getenv("CHILLER_LOGIN_PASSWORD")
	auth_url           = os.Getenv("CHILLER_AUTH_URL")
	backup_auth_url    = os.Getenv("CHILLER_BACKUP_AUTH_URL")
	backup_gallery_url = os.Getenv("CHILLER_BACKUP_GALLERY_URL")
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

var transCfg *http.Transport = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
var client http.Client = http.Client{
	Transport: transCfg,
}

func getBearerToken() (brr bearer, e error) {

	login := creds{
		Username: username,
		Password: password,
	}

	rb, _ := json.Marshal(login)

	res, err := client.Post(auth_url, "application/json", bytes.NewBuffer(rb))
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

func GetAuthServiceData() (u []dao.User, e error) {

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

	res, err := client.Do(req)
	if err != nil {
		e = err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		e = err
	}

	var users []dao.User
	_ = json.Unmarshal(body, &users)

	return users, e
}

func GetGalleryServiceData() (g []dao.Image, e error) {

	auth, err := getBearerToken()
	if err != nil {
		e = err
	}
	bearer := fmt.Sprintf("Bearer %s", auth.Access_token)

	req, err := http.NewRequest("GET", backup_gallery_url, nil)
	if err != nil {
		e = err
	}
	req.Header.Add("Authorization", bearer)

	res, err := client.Do(req)
	if err != nil {
		e = err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		e = err
	}

	var images []dao.Image
	_ = json.Unmarshal(body, &images)

	return images, e
}
