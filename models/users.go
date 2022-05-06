// Model package
package model

// Basic user struct
type User struct {
	Username string `json: "username"`
	Password string `json: "password"`
}

// Basic register user struct
type RegisterUser struct {
	Username string `json: "username"`
	Name string `json: "name"`
	Password string `json: "password"`
}

// User table struct modeled against the database
type UserTable struct {
	UserId int32 `json: "userid"`
	Username string `json: "username"`
	Name string `json: "name"`
	Password string `json: "password"`
}

// Session struct modeled against the database
type Session struct {
	SessionId int32 `json: "sessionid"`
	UserId int32 `json: "userid"`
	UserGuid string `json: "userguid"`
	IsActive bool `json: "isactive"`
}
