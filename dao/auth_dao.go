package dao

import (
	"log"
)

const AUTH_BACKUP_DB = "backup_auth"

// User Crud
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
		user.DateCreated,
		user.Enabled,
		user.AccountExpired,
		user.AccountLocked,
		user.Birthday)
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
		user.DateCreated,
		user.Enabled,
		user.AccountExpired,
		user.AccountLocked,
		user.Birthday,
		user.Id)
	if err != nil {
		return err
	}

	count, err := r.RowsAffected()
	if err != nil {
		return err
	}

	db.Close()

	log.Printf("Updated user record %d in the backup auth database. %d rows effected. ", user.Id, count)
	return err
}

func findAllUsers() (users []User, e error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "SELECT id, username, password, firstname, lastname, date_created, enabled, account_expired, account_locked, birthday FROM user"
	rs, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rs.Next() {

		var user User
		err := rs.Scan(
			&user.Id,
			&user.Username,
			&user.Password,
			&user.Firstname,
			&user.Lastname,
			&user.DateCreated,
			&user.Enabled,
			&user.AccountExpired,
			&user.AccountLocked,
			&user.Birthday,
		)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	db.Close()
	return users, e
}

func deleteUser(user User) (count int64, e error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "DELETE FROM user WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}

	r, err := stmt.Exec(user.Id)
	if err != nil {
		return 0, err
	}

	count, err = r.RowsAffected()
	if err != nil {
		return 0, err
	}

	log.Printf("Deletion of user record %d successful, %d row(s) effected", user.Id, count)
	return count, err
}

// Role Crud
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

func updateRole(role Role) (err error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "UPDATE role SET role = ?, title = ?, description = ? WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	r, err := stmt.Exec(role.Role, role.Title, role.Description, role.Id)
	if err != nil {
		return err
	}

	count, err := r.RowsAffected()
	if err != nil {
		return err
	}

	db.Close()

	log.Printf("Updated role record %d in the backup auth database. %d rows effected. ", role.Id, count)
	return err
}

func findAllRoles() (roles []Role, e error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "SELECT id, role, title, descripiton FROM role"
	rs, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rs.Next() {

		var role Role
		err := rs.Scan(&role.Id, &role.Role, &role.Title, &role.Description)
		if err != nil {
			log.Fatal(err)
		}
		roles = append(roles, role)
	}

	db.Close()
	return roles, e
}

func deleteRole(role Role) (count int64, e error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "DELETE FROM role WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}

	r, err := stmt.Exec(role.Id)
	if err != nil {
		return 0, err
	}

	count, err = r.RowsAffected()
	if err != nil {
		return 0, err
	}

	log.Printf("Deletion of role record %d successful, %d row(s) effected", role.Id, count)
	return count, err
}

// user_role crud
func insertUserRole(userid int64, ur UserRoles) (err error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "INSERT INTO user_role (id, user_id, role_id) VALUES (?, ?, ?);"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	r, err := stmt.Exec(
		ur.Id,
		userid,
		ur.Role.Id)
	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return err
	}

	db.Close()

	log.Printf("UserRole xref Record %d inserted into backup auth database.", id)
	return err
}

func deleteUserRole(ur UserRoles) (count int64, e error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "DELETE FROM user_role WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}

	r, err := stmt.Exec(ur.Id)
	if err != nil {
		return 0, err
	}

	count, err = r.RowsAffected()
	if err != nil {
		return 0, err
	}

	log.Printf("Deletion of user_role record %d successful, %d row(s) effected", ur.Id, count)
	return count, err
}
