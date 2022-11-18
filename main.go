package main

import (
	"chiller/service"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	service.BackupAuthService()
}
