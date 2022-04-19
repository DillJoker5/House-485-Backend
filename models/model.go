package model

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

type HouseTable struct {
	HouseId int32 `json: "houseid"`
	Price float32 `json: "price"`
	HouseLocation string `json: "houselocation"`
	Distance float32 `json: "distance"`
}

type UserTable struct {
	UserId int32 `json: "userid"`
	Username string `json: "username"`
	Name string `json: "name"`
	Password string `json: "password"`
	HouseId int32 `json: "houseid"`
}

type House struct {
	Price float32 `json: "price"`
	HouseLocation string `json: "houselocation"`
	Distance float32 `json: "distance"`
}

type HouseFavorite struct {
	HouseId int32 `json: "houseid"`
	Price float32 `json: "price"`
	HouseLocation string `json: "houselocation"`
	Distance float32 `json: "distance"`
	Favorite bool `json: "favorite"`
}

type User struct {
	Username string `json: "username"`
	Name string `json: "name"`
	Password string `json: "password"`
}
