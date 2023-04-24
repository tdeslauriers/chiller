package http_client

import (
	"chiller/dao"

	"crypto/tls"
	"encoding/json"
	"fmt"
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

}

func TestGetImageIds(t *testing.T) {

	auth, _ := GetBearerToken()

	ids, err := GetGalleryImageIds(auth)
	if err != nil {
		t.Log(err)
	}

	t.Log(ids)
}

func TestGetImage(t *testing.T) {

	auth, _ := GetBearerToken()
	ids, _ := GetGalleryImageIds(auth)
	for _, v := range ids {
		img, _ := GetGalleryImage(v, auth)
		t.Logf("%d - %s - %s: %s", img.Id, img.Title, img.Filename, img.Description)
	}

}

func TestGetBackupTable(t *testing.T) {

	auth, _ := GetBearerToken()
	var allowances []dao.Allowance
	t.Log(Backup_allowance_url)
	err := GetBackupTable(fmt.Sprintf("%s%s/%d", Backup_allowance_url, "/allowances", 1682077547), auth, &allowances)
	if err != nil {
		t.Log(err)
	}

	for _, a := range allowances {
		t.Logf("id: %v", a)
	}
}
