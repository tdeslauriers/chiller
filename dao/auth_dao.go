package dao

import (
	"fmt"
	"log"
)

const AUTH_BACKUP_DB = "backup_auth"

func insertUser(user User) (err error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "INSERT INTO user (id, username, password, firstname, lastname, date_created, enabled, account_expired, account_locked, birthday) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	r, err := stmt.Exec(
		user.Id,
		user.Username,
		user.Password,
		user.Firstname,
		user.Lastname,
		fmt.Sprintf("%d-%d-%d", user.DateCreated[0], user.DateCreated[1], user.DateCreated[2]),
		user.Enabled,
		user.AccountExpired,
		user.AccountLocked,
		fmt.Sprintf("%d-%d-%d", user.Birthday[0], user.Birthday[1], user.Birthday[2]))
	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return err
	}

	db.Close()

	log.Printf("User Record %d inserted into backup auth database.", id)
	return err
}

func updateUser(user User) (err error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "UPDATE user SET username = ?, password = ?, firstname = ?, lastname = ?, date_created = ?, enabled = ?, account_expired = ?, account_locked = ?, birthday = ? WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	r, err := stmt.Exec(
		user.Username,
		user.Password,
		user.Firstname,
		user.Lastname,
		fmt.Sprintf("%d-%d-%d", user.DateCreated[0], user.DateCreated[1], user.DateCreated[2]),
		user.Enabled,
		user.AccountExpired,
		user.AccountLocked,
		fmt.Sprintf("%d-%d-%d", user.Birthday[0], user.Birthday[1], user.Birthday[2]),
		user.Id)
	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return err
	}

	db.Close()

	log.Printf("Updated User record %d in the backup auth database.", id)
	return err
}

// returns Id if present => clunky
func findAllUsers() (ids []int64, e error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "SELECT id FROM user"
	rs, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	all := make([]int64, 0, 100)
	for rs.Next() {

		var id int64
		err := rs.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		all = append(all, id)
	}

	db.Close()
	return all, e
}

func insertRole(role Role) (err error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "INSERT INTO role (id, role, title, description) VALUES (?, ?, ?, ?);"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	r, err := stmt.Exec(
		role.Id,
		role.Role,
		role.Title,
		role.Description)
	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return err
	}

	db.Close()

	log.Printf("Role Record %d inserted into backup auth database.", id)
	return err
}
