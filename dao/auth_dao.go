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

type row interface {
	User | Role | Address | Phone | UserRoles | UserAddresses | UserPhones | UrXref | UaXref | UpXref
}

type Record[T row] struct {
	Id int64
}

type XrefRecord[T row] struct {
	Id   int64
	Fk_1 int64
	Fk_2 int64
}

// Xref inserts + Find alls
const (
	INSERT_UR  = "INSERT INTO user_role (id, user_id, role_id) VALUES (?, ?, ?);"
	INSERT_UA  = "INSERT INTO user_address (id, user_id, address_id) VALUES (?, ?, ?);"
	INSERT_UP  = "INSERT INTO user_phone (id, user_id, phone_id) VALUES (?, ?, ?);"
	FINDALL_UR = "SELECT id, user_id, role_id FROM user_role"
	FINDALL_UA = "SELECT id, user_id, address_id FROM user_address"
	FINDALL_UP = "SELECT id, user_id, phone_id FROM user_phone"
)

func InsertXrefRecord[T row](r XrefRecord[T], query string) (err error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(
		r.Id,
		r.Fk_1,
		r.Fk_2)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	db.Close()

	log.Printf("%T xref record %d inserted into backup auth database.", r, id)
	return err
}

func FindAllXrefs[T row](query string) (xrefs []XrefRecord[T], e error) {
	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	rs, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rs.Next() {

		var xr XrefRecord[T]
		err := rs.Scan(&xr.Id, &xr.Fk_1, &xr.Fk_2)
		if err != nil {
			log.Fatal(err)
		}
		xrefs = append(xrefs, xr)
	}

	db.Close()
	return xrefs, e
}

// Delete Records:
const (
	DELETE_USER    = "DELETE FROM user WHERE id = ?"
	DELETE_ROLE    = "DELETE FROM role WHERE id = ?"
	DELETE_PHONE   = "DELETE FROM phone WHERE id = ?"
	DELETE_ADDRESS = "DELETE FROM address WHERE id = ?"
	DELETE_UR      = "DELETE FROM user_role WHERE id = ?"
	DELETE_UA      = "DELETE FROM user_address WHERE id = ?"
	DELETE_UP      = "DELETE FROM user_phone WHERE id = ?"
)

func DeleteRecord[T row](record Record[T], query string) (e error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(record.Id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count > 0 {
		log.Printf("Deletion of %T record %d successful; %d row(s) affected", record, record.Id, count)
	}
	return err
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
	if count > 0 {
		log.Printf("Updated role record %d in the backup auth database. %d rows affected. ", role.Id, count)
	}
	return err
}

func FindAllRoles() (roles []Role, e error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "SELECT id, role, title, description FROM role"
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

// address crud
func InsertAddress(address Address) (err error) {

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

func UpdateAddress(address Address) (err error) {

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
	if count > 0 {
		log.Printf("Updated address record %d in the backup auth database. %d rows affected. ", address.Id, count)
	}
	return err
}

func FindAllAddresses() (addresses []Address, e error) {

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

// phone crud
func InsertPhone(phone Phone) (err error) {

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

func UpdatePhone(phone Phone) (err error) {

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
	if count > 0 {
		log.Printf("Updated phone record %d in the backup auth database. %d rows affected. ", phone.Id, count)
	}
	return err
}

func FindAllPhones() (phones []Phone, e error) {

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
