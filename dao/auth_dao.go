package dao

import (
	"database/sql"
	"log"
)

const AUTH_BACKUP_DB = "backup_auth"

// birthday can be empty
func birthdayNullString(bd string) sql.NullString {
	if len(bd) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: bd,
		Valid:  true,
	}
}

// User Crud
func InsertUser(user User) (err error) {

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
		birthdayNullString(user.Birthday))
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

func UpdateUser(user User) (err error) {

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
		birthdayNullString(user.Birthday),
		user.Id)
	if err != nil {
		return err
	}

	count, err := r.RowsAffected()
	if err != nil {
		return err
	}

	db.Close()
	if count > 0 {
		log.Printf("Updated user record %d in the backup auth database.)", user.Id)
	}
	return err
}

func FindAllUsers() (users []User, e error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "SELECT id, username, password, firstname, lastname, date_created, enabled, account_expired, account_locked, COALESCE(birthday, '') FROM user"
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

func DeleteUser(user User) (count int64, e error) {

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
func InsertRole(role Role) (err error) {

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

func UpdateRole(role Role) (err error) {

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

func FindAllRoles() (roles []Role, e error) {

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

func DeleteRole(role Role) (count int64, e error) {

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

// address crud
func insertAddress(address Address) (err error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "INSERT INTO address (id, address, city, state, zip) VALUES (?, ?, ?, ?, ?);"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	r, err := stmt.Exec(
		address.Id,
		address.Address,
		address.City,
		address.State,
		address.Zip)
	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return err
	}

	db.Close()

	log.Printf("Address Record %d inserted into backup auth database.", id)
	return err
}

func updateAddress(address Address) (err error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "UPDATE address SET address = ?, city = ?, state = ?, zip = ? WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	r, err := stmt.Exec(address.Address, address.City, address.State, address.Zip, address.Id)
	if err != nil {
		return err
	}

	count, err := r.RowsAffected()
	if err != nil {
		return err
	}

	db.Close()

	log.Printf("Updated address record %d in the backup auth database. %d rows effected. ", address.Id, count)
	return err
}

func findAllAddresses() (addresses []Address, e error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "SELECT id, address, city, state, zip FROM address"
	rs, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rs.Next() {

		var address Address
		err := rs.Scan(&address.Id, &address.Address, &address.City, &address.State, &address.Zip)
		if err != nil {
			log.Fatal(err)
		}
		addresses = append(addresses, address)
	}

	db.Close()
	return addresses, e
}

func deleteAddress(address Address) (count int64, e error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "DELETE FROM address WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}

	r, err := stmt.Exec(address.Id)
	if err != nil {
		return 0, err
	}

	count, err = r.RowsAffected()
	if err != nil {
		return 0, err
	}

	log.Printf("Deletion of address record %d successful, %d row(s) affected", address.Id, count)
	return count, err
}

// user_address crud
func insertUserAdress(userid int64, ua UserAddresses) (err error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "INSERT INTO user_address (id, user_id, address_id) VALUES (?, ?, ?);"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	r, err := stmt.Exec(
		ua.Id,
		userid,
		ua.Address.Id)
	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return err
	}

	db.Close()

	log.Printf("UserAddress xref Record %d inserted into backup auth database.", id)
	return err
}

func deleteUserAddress(ua UserAddresses) (count int64, e error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "DELETE FROM user_address WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}

	r, err := stmt.Exec(ua.Id)
	if err != nil {
		return 0, err
	}

	count, err = r.RowsAffected()
	if err != nil {
		return 0, err
	}

	log.Printf("Deletion of user_address record %d successful, %d row(s) effected", ua.Id, count)
	return count, err
}

// phone crud
func insertPhone(phone Phone) (err error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "INSERT INTO phone (id, phone, type) VALUES (?, ?, ?);"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	r, err := stmt.Exec(phone.Id, phone.Phone, phone.Type)
	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return err
	}

	db.Close()

	log.Printf("Phone Record %d inserted into backup auth database.", id)
	return err
}

func updatePhone(phone Phone) (err error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "UPDATE phone SET phone = ?, type = ? WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	r, err := stmt.Exec(phone.Phone, phone.Type, phone.Id)
	if err != nil {
		return err
	}

	count, err := r.RowsAffected()
	if err != nil {
		return err
	}

	db.Close()

	log.Printf("Updated phone record %d in the backup auth database. %d rows effected. ", phone.Id, count)
	return err
}

func findAllPhones() (phones []Phone, e error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "SELECT id, phone, type FROM phone"
	rs, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rs.Next() {

		var phone Phone
		err := rs.Scan(&phone.Id, &phone.Phone, &phone.Type)
		if err != nil {
			log.Fatal(err)
		}
		phones = append(phones, phone)
	}

	db.Close()
	return phones, e
}

func deletePhone(phone Phone) (count int64, e error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "DELETE FROM phone WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}

	r, err := stmt.Exec(phone.Id)
	if err != nil {
		return 0, err
	}

	count, err = r.RowsAffected()
	if err != nil {
		return 0, err
	}

	log.Printf("Deletion of phone record %d successful, %d row(s) affected", phone.Id, count)
	return count, err
}

// user_phone crud
func insertUserPhone(userid int64, up UserPhones) (err error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "INSERT INTO user_phone (id, user_id, phone_id) VALUES (?, ?, ?);"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	r, err := stmt.Exec(
		up.Id,
		userid,
		up.Phone.Id)
	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return err
	}

	db.Close()

	log.Printf("UserPhone xref Record %d inserted into backup auth database.", id)
	return err
}

func deleteUserPhone(up UserPhones) (count int64, e error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "DELETE FROM user_phone WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}

	r, err := stmt.Exec(up.Id)
	if err != nil {
		return 0, err
	}

	count, err = r.RowsAffected()
	if err != nil {
		return 0, err
	}

	log.Printf("Deletion of user_phone record %d successful, %d row(s) effected", up.Id, count)
	return count, err
}
