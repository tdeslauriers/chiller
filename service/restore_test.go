package service

import (
	"chiller/http_client"
	"testing"
)

func TestRestoreAuthService(t *testing.T) {

	auth, _ := http_client.GetBearerToken()

	err := RestoreAuthService(auth)
	t.Logf("Restored Auth User Service Data: %v", err)

}

func TestRestoreAllowanceService(t *testing.T) {

	auth, _ := http_client.GetBearerToken()

	err := RestoreAllowanceService(auth)
	t.Logf("Restored Auth User Service Data: %v", err)
}