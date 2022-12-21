package dao

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
	Id          int64         `json:"id"`
	Filename    string        `json:"filename"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Date        string        `json:"date"`
	Published   bool          `json:"published"`
	Thumbnail   string        `json:"thumbnail"`
	Image       string        `json:"image"`
	AlbumImages []AlbumImages `json:"albumImages"`
}

type Album struct {
	Id    int64  `json:"id"`
	Album string `json:"album"`
}

type AlbumImages struct {
	Id    int64 `json:"id"`
	Album Album `json:"album"`
}
