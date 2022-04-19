package main

import (
	"encoding/json"
	"log"
	"net/http"
	model "house-485-backend/models"
	"github.com/gorilla/mux"
	"database/sql"
	"context"
	_"github.com/denisenkom/go-mssqldb"
	"fmt"
)

var db *sql.DB
var server = "DESKTOP-K7IIMGF"
var port = 1433
var user = "capstoneapiuser"
var password = "CapstoneApiUser2022!"
var database = "House485Database"

func main() {
	// Build connection string
	connectionString  := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", server, user, password, port, database)
	var err error

	// Create connection
	db, err = sql.Open("sqlserver", connectionString)
	if err != nil {
		log.Fatal("Error connecting:", err.Error())
	}

	router := mux.NewRouter()

	router.HandleFunc("/loginUser", mwCheck(readUserTable)).Methods(http.MethodPost)
	router.HandleFunc("/home", mwCheck(readHouseTable)).Methods(http.MethodPost)
	router.HandleFunc("/newFavorite", mwCheck(createNewHouse)).Methods(http.MethodPost)
	router.HandleFunc("/registerUser", mwCheck(createNewUser)).Methods(http.MethodPost)
	router.HandleFunc("/updateHouses", mwCheck(updateHouseTable)).Methods(http.MethodPost)

	srv := &http.Server {
		Addr: ":8000",
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func mwCheck(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle authentication
		// if valid auth f(w, r)
		// else send_error(w, r)
		f(w, r)
	}
}

func readHouseTable(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Verify database is running
	err := db.PingContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tsqlQuery := fmt.Sprintf("SELECT HouseId, Price, HouseLocation, Distance FROM House")

	// Execute query
	rows, err := db.QueryContext(ctx, tsqlQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	defer rows.Close()

	var houses []model.HouseTable
	for rows.Next() {
		var house model.HouseTable
		rows.Scan(&house.HouseId, &house.Price, &house.HouseLocation, &house.Distance)
		houses = append(houses, house)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var response = model.HouseJsonResponse{ Type: "Success", Data: houses }
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func readUserTable(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Verify database is running
	err := db.PingContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tsqlQuery := fmt.Sprintf("SELECT UserId, Username, Name, Password, HouseId FROM Users")

	// Execute query
	rows, err := db.QueryContext(ctx, tsqlQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	defer rows.Close()

	var users []model.UserTable
	for rows.Next() {
		var user model.UserTable
		rows.Scan(&user.UserId, &user.Username, &user.Name, &user.Password, &user.HouseId)
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var response = model.UserJsonResponse{ Type: "Success", Data: users}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func getHouseInfo(w http.ResponseWriter, r *http.Request) (model.House) {
	var house model.House
	r.ParseForm()
	if r.Form == nil || r.Form["Price"] == nil || r.Form["HouseLocation"] == nil || r.Form["Distance"] == nil {
		http.Error(w, "Cannot Access House Information from Form", http.StatusInternalServerError)
	} else {
		house.Price = r.Form["Price"]
		house.HouseLocation = r.Form["HouseLocation"]
		house.Distance = r.Form["Distance"]
	}
	return house
}

func getUserInfo(w http.ResponseWriter, r *http.Request) (model.User) {
	var user model.User
	r.ParseForm()
	if r.Form == nil || r.Form["Username"] == nil || r.Form["Name"] == nil || r.Form["Password"] == nil {
		http.Error(w, "Cannot Access User Information from Form", http.StatusInternalServerError)
	} else {
		user.Username = r.Form["Username"]
		user.Name = r.Form["Name"]
		user.Password = r.Form["Password"]
	}
	return user
}

func createNewHouse(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Verify database is running
	err := db.PingContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Build new house
	var newHouse model.House = getHouseInfo(w, r)

	tsqlMutation := fmt.Sprintf("INSERT INTO House VALUES(%d, %s, %d)", &newHouse.Price, &newHouse.HouseLocation, &newHouse.Distance)

	// Execute query
	newHouseRow, err := db.QueryContext(ctx, tsqlMutation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	defer newHouseRow.Close()

	// Grab new house
	var rows []model.HouseTable
	for newHouseRow.Next() {
		var row model.HouseTable
		newHouseRow.Scan(&row.HouseId, &row.Price, &row.HouseLocation, &row.Distance)
		rows = append(rows, row)
	}

	// Throw error if house could not be grabbed
	if rows == nil {
		http.Error(w, "Cannot find newly created House", http.StatusBadRequest)
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var response = model.HouseJsonResponse{ Type: "Success", Data: rows }
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func createNewUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Verify database is running
	err := db.PingContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
	// Build new user from register form
	var newUser model.User = getUserInfo(w, r)

	tsqlMutation := fmt.Sprintf("INSERT INTO Users VALUES(%s, %s, %s)", &newUser.Username, &newUser.Name, &newUser.Password)

	// Execute query
	newUserRow, err := db.QueryContext(ctx, tsqlMutation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	defer newUserRow.Close()

	// Grab new user
	var rows []model.UserTable
	for newUserRow.Next() {
		var row model.UserTable
		newUserRow.Scan(&row.UserId, &row.Username, &row.Name, &row.Password, &row.HouseId)
		rows = append(rows, row)
	}

	// Throw error if user could not be grabbed
	if rows == nil {
		http.Error(w, "Cannot register your account with current information", http.StatusBadRequest)
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var response = model.UserJsonResponse{ Type: "Success", Data: rows }
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func updateHouseTable(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Verify database is running
	err := db.PingContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Grab house that needs to be updated
	var updateHouse model.HouseTable = getUpdateHouseInfo(w, r)

	var tsqlMutation string

	if updateHouse.favorite {
		// create new House entry
	} else {
		// delete House from House Table
	}

	// Execute query
	houseRow, err := db.QueryContext(ctx, tsqlMutation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	defer houseRow.Close()
}
