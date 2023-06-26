package dao

import "time"

type Backup struct {
	Id     int64     `db:"id"`
	Backup time.Time `db:"backup"`
}

// nested structs have been pulled appart for re-use

// auth persistance objects
type User struct {
	Id             int64  `json:"id" db:"id"`
	Username       string `json:"username" db:"username"`
	Password       string `json:"password" db:"password"`
	Firstname      string `json:"firstname" db:"firstname"`
	Lastname       string `json:"lastname" db:"lastname"`
	DateCreated    string `json:"dateCreated" db:"date_created"`
	Enabled        bool   `json:"enabled" db:"enabled"`
	AccountExpired bool   `json:"accountExpired" db:"account_expired"`
	AccountLocked  bool   `json:"accountLocked" db:"account_locked"`
	Uuid           string `json:"uuid" db:"uuid"`
	Birthday       string `json:"birthday" db:"birthday"`
}

type UserRoles struct {
	Id   int64 `json:"id"`
	Role Role  `json:"role"`
}

type UrXref struct {
	Id      int64 `json:"id" db:"id"`
	User_id int64 `json:"userId" db:"user_id"`
	Role_id int64 `json:"roleId" db:"role_id"`
}

type Role struct {
	Id          int64  `json:"id" db:"id"`
	Role        string `json:"role" db:"role"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
}

type UserAddresses struct {
	Id      int64   `json:"id"`
	Address Address `json:"address"`
}

type UaXref struct {
	Id         int64
	User_id    int64
	Address_id int64
}

type Address struct {
	Id      int64  `json:"id"`
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
}

type UserPhones struct {
	Id    int64 `json:"id"`
	Phone Phone `json:"phone"`
}

type UpXref struct {
	Id       int64
	User_id  int64
	Phone_id int64
}

type Phone struct {
	Id    int64  `json:"id"`
	Phone string `json:"phone"`
	Type  string `json:"type"`
}

// gallery persistence objects
type Image struct {
	Id          int64  `json:"id" db:"id"`
	Filename    string `json:"filename" db:"filename"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Date        string `json:"date" db:"date"`
	Published   bool   `json:"published" db:"published"`
	// Thumbnail    []byte `json:"thumbnail" db:"thumbnail"`
	// Presentation []byte `json:"presentation" db:"presentation"`
	// Image        []byte `json:"image" db:"image"`
}

type Album struct {
	Id    int64  `json:"id" db:"id"`
	Album string `json:"album" db:"album"`
}

type AlbumImages struct {
	Id    int64 `json:"id"`
	Album Album `json:"album"`
	Image Image
}

type AiXref struct {
	Id       int64 `json:"id" db:"id"`
	Album_id int64 `json:"album_id" db:"album_id"`
	Image_id int64 `json:"image_id" db:"image_id"`
}

// allowance persistence objects
type Allowance struct {
	Id       int64  `json:"userId" db:"id"`
	Balance  string `json:"balance" db:"balance"`
	UserUuid string `json:"userUuid" db:"user_uuid"`
}

type Task struct {
	Id           int64  `json:"id" db:"id"`
	Date         string `json:"date" db:"date"`
	Complete     string `json:"complete" db:"complete"`
	Satisfactory string `json:"satisfactory" db:"satisfactory"`
	TasktypeId   string `json:"tasktypeId" db:"tasktype_id"`
}

type TaskAllowance struct {
	Id          int64  `json:"id" db:"id"`
	TaskId      string `json:"taskId" db:"task_id"`
	AllowanceId string `json:"allowanceId" db:"allowance_id"`
}

type Tasktype struct {
	Id       int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Cadence  string `json:"cadence" db:"cadence"`
	Category string `json:"category" db:"category"`
	Archived string `json:"archived" db:"archived"`
}

type TasktypeAllowance struct {
	Id          int64  `json:"id" db:"id"`
	TasktypeId  string `json:"tasktypeId" db:"tasktype_id"`
	AllowanceId string `json:"allowanceId" db:"allowance_id"`
}

type CleanupAllowance struct {
	AllowanceIds         []int64 `json:"allowanceIds"`
	TasktypeIds          []int64 `json:"tasktypeIds"`
	TaskIds              []int64 `json:"taskIds"`
	TasktypeAllowanceIds []int64 `json:"tasktypeAllowanceIds"`
	TaskAllowanceId      []int64 `json:"taskAllowanceIds"`
}
