package model

type GenericJsonResponse struct {
	Type string `json: "type"`
	Message string `json: "message"`
}

type JsonLoginResponse struct {
	Type string `json: "type"`
	Message string `json: "message"`
	UserGuid string `json: "userguid"`
	UserId int32 `json: "userid"`
}

type HouseJsonResponse struct {
	Type string `json: "type"`
	Message string `json: "message"`
	Data []HouseTable `json: "data"`
}

type UserJsonResponse struct {
	Type string `json: "type"`
	Message string `json: "message"`
	Data []UserTable `json: "data"`
}
