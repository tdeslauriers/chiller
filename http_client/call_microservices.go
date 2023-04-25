package http_client

import (
	"bytes"
	"chiller/dao"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	username             = os.Getenv("CHILLER_LOGIN_USERNAME")
	password             = os.Getenv("CHILLER_LOGIN_PASSWORD")
	auth_url             = os.Getenv("CHILLER_AUTH_URL")
	Backup_auth_url      = os.Getenv("CHILLER_BACKUP_AUTH_URL")
	Backup_gallery_url   = os.Getenv("CHILLER_BACKUP_GALLERY_URL")
	Backup_allowance_url = os.Getenv("CHILLER_BACKUP_ALLOWANCE_URL")
)

type creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Bearer struct {
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
	Timeout:   180 * time.Second,
}

func GetBearerToken() (brr Bearer, e error) {

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

		b := Bearer{}
		body, err := io.ReadAll(res.Body)
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

func GetAuthServiceData(t Bearer) (u []dao.User, e error) {

	bearer := fmt.Sprintf("Bearer %s", t.Access_token)

	req, err := http.NewRequest("GET", Backup_auth_url, nil)
	if err != nil {
		e = err
	}
	req.Header.Add("Authorization", bearer)

	res, err := client.Do(req)
	if err != nil {
		e = err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		e = err
	}

	var users []dao.User
	_ = json.Unmarshal(body, &users)

	return users, e
}

func GetGalleryImageIds(t Bearer) (gids []int64, e error) {

	bearer := fmt.Sprintf("Bearer %s", t.Access_token)

	req, err := http.NewRequest("GET", Backup_gallery_url+"/list", nil)
	if err != nil {
		e = err
	}
	req.Header.Add("Authorization", bearer)

	res, err := client.Do(req)
	if err != nil {
		e = err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		e = err
	}

	_ = json.Unmarshal(body, &gids)

	return gids, e
}

func GetGalleryImage(id int64, t Bearer) (image dao.Image, e error) {

	bearer := fmt.Sprintf("Bearer %s", t.Access_token)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%d", Backup_gallery_url, id), nil)
	if err != nil {
		e = err
	}
	req.Header.Add("Authorization", bearer)
	req.Close = true

	res, err := client.Do(req)
	if err != nil {
		return image, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		e = err
	}

	_ = json.Unmarshal(body, &image)

	return image, e

}

func GetBackupTable(endpoint string, t Bearer, v interface{}) error {

	bearer := fmt.Sprintf("Bearer %s", t.Access_token)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", bearer)
	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return err
	}

	return nil
}
