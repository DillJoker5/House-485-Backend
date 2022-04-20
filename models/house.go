package model

type HouseTable struct {
	HouseId int32 `json: "houseid"`
	Price float32 `json: "price"`
	HouseLocation string `json: "houselocation"`
	Distance float32 `json: "distance"`
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
