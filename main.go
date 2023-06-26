package main

import (
	"chiller/http_client"
	"chiller/service"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	auth, err := http_client.GetBearerToken()
	if err != nil {
		log.Fatal(err)
	}

	// service.BackupAuthService(auth)
	// service.BackupGalleryService(auth)

	// service.RestoreAuthService(auth)
	// service.RestoreAllowanceService(auth)
	service.RestoreGalleryService(auth)

}
