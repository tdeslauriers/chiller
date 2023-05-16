package dao

import "time"

type Backup struct {
	Id     int64     `db:"id"`
	Backup time.Time `db:"backup"`
}

// nested structs have been pulled appart for re-use

// auth persistance objects
type User struct {
	Id             int64           `json:"id"`
	Username       string          `json:"username"`
	Password       string          `json:"password"`
	Firstname      string          `json:"firstname"`
	Lastname       string          `json:"lastname"`
	DateCreated    string          `json:"dateCreated"`
	Enabled        bool            `json:"enabled"`
	AccountExpired bool            `json:"accountExpired"`
	AccountLocked  bool            `json:"accountLocked"`
	Birthday       string          `json:"birthday"`
	Uuid           string          `json:"uuid"`
	UserRoles      []UserRoles     `json:"userRoles"`
	UserAddresses  []UserAddresses `json:"userAddresses"`
	UserPhones     []UserPhones    `json:"userPhones"`
}

type UserRoles struct {
	Id   int64 `json:"id"`
	Role Role  `json:"role"`
}

type UrXref struct {
	Id      int64
	User_id int64
	Role_id int64
}

type Role struct {
	Id          int64  `json:"id"`
	Role        string `json:"role"`
	Title       string `json:"title"`
	Description string `json:"description"`
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
	Id           int64         `json:"id"`
	Filename     string        `json:"filename"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Date         string        `json:"date"`
	Published    bool          `json:"published"`
	Thumbnail    []byte        `json:"thumbnail"`
	Presentation []byte        `json:"presentation"`
	Image        []byte        `json:"image"`
	AlbumImages  []AlbumImages `json:"albumImages"`
}

type Album struct {
	Id    int64  `json:"id"`
	Album string `json:"album"`
}

type AlbumImages struct {
	Id    int64 `json:"id"`
	Album Album `json:"album"`
	Image Image
}

type AiXref struct {
	Id       int64
	Album_id int64
	Image_id int64
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
