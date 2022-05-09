// Model package
package model

// House table struct modeled against the database
type HouseTable struct {
	HouseId int32 `json: "houseid"`
	Price float32 `json: "price"`
	HouseLocation string `json: "houselocation"`
	Distance float32 `json: "distance"`
	UserId int32 `json: "userid"`
}

/* House favorite struct modeled against the database and adding
   the favorite value and the user guid*/
type HouseFavorite struct {
	HouseId int32 `json: "houseid"`
	Price float32 `json: "price"`
	HouseLocation string `json: "houselocation"`
	Distance float32 `json: "distance"`
	UserId int32 `json: "userid"`
	Favorite bool `json: "favorite"`
	UserGuid string `json: "userguid"`
}
