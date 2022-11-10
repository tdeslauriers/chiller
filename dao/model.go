package dao

// nested structs have been pulled appart for re-use

// auth persistance objects
type User struct {
	Id             int64           `json:"id"`
	Username       string          `json:"username"`
	Password       string          `json:"password"`
	Firstname      string          `json:"firstname"`
	Lastname       string          `json:"lastname"`
	DateCreated    []int           `json:"dateCreated"`
	Enabled        bool            `json:"enabled"`
	AccountExpired bool            `json:"accountExpired"`
	AccountLocked  bool            `json:"accountLocked"`
	Birthday       []int           `json:"birthday"`
	UserRoles      []UserRoles     `json:"userRoles"`
	UserAddresses  []UserAddresses `json:"userAddresses"`
	UserPhones     []UserPhones    `json:"userPhones"`
}

type UserRoles struct {
	Id   int64 `json:"id"`
	Role Role  `json:"role"`
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

type Phone struct {
	Id    int64  `json:"id"`
	Phone string `json:"phone"`
	Type  string `json:"type"`
}

// gallery persistence objects
