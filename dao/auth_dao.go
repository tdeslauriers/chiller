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

// user_role crud
func InsertUserRole(ur UrXref) (err error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "INSERT INTO user_role (id, user_id, role_id) VALUES (?, ?, ?);"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	r, err := stmt.Exec(
		ur.Id,
		ur.User_id,
		ur.Role_id)
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

func FindAllUserroles() (urs []UrXref, e error) {
	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "SELECT id, user_id, role_id FROM user_role"
	rs, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rs.Next() {

		var ur UrXref
		err := rs.Scan(&ur.Id, &ur.User_id, &ur.Role_id)
		if err != nil {
			log.Fatal(err)
		}
		urs = append(urs, ur)
	}

	db.Close()
	return urs, e
}

func FindUserRolesByUserId(id int64) (ur []UrXref, e error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "SELECT id, user_id, role_id FROM user_role WHERE user_id = ?"
	rs, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}

	for rs.Next() {

		var xref UrXref
		err := rs.Scan(&xref.Id, &xref.User_id, &xref.Role_id)
		if err != nil {
			log.Fatal(err)
		}
		ur = append(ur, xref)
	}

	db.Close()
	return ur, e
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

// user_address crud
func FindAllUserAddresses() (uas []UaXref, e error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "SELECT id, user_id, address_id FROM user_address"
	rs, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rs.Next() {

		var ua UaXref
		err := rs.Scan(&ua.Id, &ua.User_id, &ua.Address_id)
		if err != nil {
			log.Fatal(err)
		}
		uas = append(uas, ua)
	}

	db.Close()
	return uas, e
}

func InsertUserAdress(ua UaXref) (err error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "INSERT INTO user_address (id, user_id, address_id) VALUES (?, ?, ?);"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	r, err := stmt.Exec(
		ua.Id,
		ua.User_id,
		ua.Address_id)
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

// user_phone crud
func FindAllUserPhones() (ups []UpXref, e error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "SELECT id, user_id, phone_id FROM user_phone"
	rs, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rs.Next() {

		var up UpXref
		err := rs.Scan(&up.Id, &up.User_id, &up.Phone_id)
		if err != nil {
			log.Fatal(err)
		}
		ups = append(ups, up)
	}

	db.Close()
	return ups, e
}

func InsertUserPhone(up UpXref) (err error) {

	db := dbConn(AUTH_BACKUP_DB)
	defer db.Close()

	query := "INSERT INTO user_phone (id, user_id, phone_id) VALUES (?, ?, ?);"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	r, err := stmt.Exec(
		up.Id,
		up.User_id,
		up.Phone_id)
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

// Delete Records:
type row interface {
	User | Role | Address | Phone | UserRoles | UserAddresses | UserPhones | UrXref | UaXref | UpXref
}

type Record[T row] struct {
	Id int64
}

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
