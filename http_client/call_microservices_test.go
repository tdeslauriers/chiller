package http_client

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type health struct {
	Status string `json:"status"`
}

// http refresher
func TestHttp(t *testing.T) {

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Timeout:   5 * time.Second,
		Transport: transCfg,
	}

	// get -> health check
	res, err := client.Get("https://deslauriers.world/api/v1/gateway/health")
	if err != nil {
		t.Log(err)
	}
	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		t.Log(readErr)
	}

	h := health{}
	jsonErr := json.Unmarshal(body, &h)
	if jsonErr != nil {
		t.Log(jsonErr)
	}

	t.Logf("\nHttp Response Code: %d\nBody: service status --> %s", res.StatusCode, h.Status)
}

func TestGetBearerToken(t *testing.T) {

	t.Logf("username: %s", username)

	bearer, _ := GetBearerToken()
	t.Log(bearer.Access_token)
	t.Log(bearer.Roles)
}

func TestGetAuthServiceData(t *testing.T) {

	auth, _ := GetBearerToken()
	// Testing against known data set
	users, _ := GetAuthServiceData(auth)
	for _, v := range users {
		t.Log(v)
	}

	if users[0].Lastname != "Skywalker" {
		t.Fail()
		t.Logf("Expected %s; Actual: %s", "Skywalker", users[0].Lastname)
	}

	if users[0].UserRoles[0].Role.Title != "General Admission" {
		t.Fail()
		t.Logf("Expected %s; Actual: %s", "GeneralAdmission", users[0].UserRoles[0].Role.Title)
	}

	if users[0].UserAddresses[0].Address.City != "Hoth" {
		t.Fail()
		t.Logf("Expected %s; Actual: %s", "Hoth", users[0].UserAddresses[0].Address.City)
	}

}

func TestGetImageIds(t *testing.T) {

	auth, _ := GetBearerToken()

	gids, _ := GetGalleryImageIds(auth)
	for _, v := range gids {
		t.Log(v)
	}
}
