package model

type GenericJsonResponse struct {
	Type string `json: "type"`
	Message string `json: "message"`
}

type JsonLoginResponse struct {
	Type string `json: "type"`
	Message string `json: "message"`
	UserGuid string `json: "userguid"`
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
