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

	lastBackup, err := service.GetLastBackup()
	if err != nil {
		log.Fatal(err)
	}

	// Allownace Service
	// Allowance table
	if err := service.BackupAllowanceService(lastBackup, auth); err != nil {
		log.Panic(err)
	}

	// service.BackupAuthService(auth)
	// service.BackupGalleryService(auth)
}
