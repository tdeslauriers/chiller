package http_client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

type health struct {
	Status string `json:"status"`
}

// http refresher
func TestHttp(t *testing.T) {

	// get -> health check
	res, err := http.Get("http://localhost:8080/health")
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

	bearer, _ := getBearerToken()
	t.Log(bearer.Access_token)
}
