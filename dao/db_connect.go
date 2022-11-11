package dao

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var (
	user = os.Getenv("CHILLER_BACKUP_DB_USERNAME")
	pass = os.Getenv("CHILLER_BACKUP_DB_PASSWORD")
	dbIP = os.Getenv("CHILLER_BACKUP_DB_IP")
	// need to pass in db name on creation for different dbs
)

// DBConn is db connector function
func dbConn(name string) *sql.DB {
	var url = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, pass, dbIP, name)
	db, err := sql.Open("mysql", url)
	if err != nil {
		fmt.Printf("Cannot connect to database: %s/%s\n", dbIP, name)
		log.Fatal("Database connection error: ", err)
	} else {
		fmt.Printf("Connected to: %s/%s\n", dbIP, name)
	}
	return db
}
