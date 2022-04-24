package model

type UserTable struct {
	UserId int32 `json: "userid"`
	Username string `json: "username"`
	Name string `json: "name"`
	Password string `json: "password"`
}

type Session struct {
	SessionId int32 `json: "sessionid"`
	UserId int32 `json: "userid"`
	UserGuid string `json: "userguid"`
	IsActive bool `json: "isactive"`
}
