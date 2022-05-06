// Model package
package model

// Generic json response struct
type GenericJsonResponse struct {
	Type string `json: "type"`
	Message string `json: "message"`
}

// Login json response with the user guid and user id
type JsonLoginResponse struct {
	Type string `json: "type"`
	Message string `json: "message"`
	UserGuid string `json: "userguid"`
	UserId int32 `json: "userid"`
}

// House json response with the data from house table
type HouseJsonResponse struct {
	Type string `json: "type"`
	Message string `json: "message"`
	Data []HouseTable `json: "data"`
}

// User json response with the data from users table
type UserJsonResponse struct {
	Type string `json: "type"`
	Message string `json: "message"`
	Data []UserTable `json: "data"`
}
