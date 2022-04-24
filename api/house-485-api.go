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
	"github.com/google/uuid"
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

	// Handle api requests
	router.HandleFunc("/readUser", ReadUserTable).Methods(http.MethodPost)
	router.HandleFunc("/home", mwCheck(ReadHouseTable)).Methods(http.MethodPost)
	router.HandleFunc("/login", Login).Methods(http.MethodPost)
	router.HandleFunc("/logout", Logout).Methods(http.MethodPost)
	router.HandleFunc("/register", Register).Methods(http.MethodPost)
	router.HandleFunc("/favorite", HouseFavorites).Methods(http.MethodPost)
	router.HandleFunc("/updateFavorite", UpdateFavorite).Methods(http.MethodPost)

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
		if !validateUser(r) {
			http.Error(w, "Unauthorized user", http.StatusForbidden)
		} else {
			f(w, r)
		}
	}
}

func validateUser(r *http.Request) bool {
	ctx := context.Background()
	
	// Verify database is running
	err := db.PingContext(ctx)
	if err != nil {
		return false
	}

	var s model.Session
	err = json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		return false
	}

	if s.UserGuid == "" {
		return false
	}

	tsqlQuery := fmt.Sprintf("SELECT SessionId FROM Session WHERE UserGuid='%s' AND IsActive=1;", s.UserGuid)

	row := db.QueryRowContext(ctx, tsqlQuery)

	var sid int32
	if err = row.Scan(&sid); err != nil {
		return false
	}

	return true
}

func ReadHouseTable(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Verify database is running
	err := db.PingContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tsqlQuery := "SELECT HouseId, Price, HouseLocation, Distance, UserId FROM House;"

	// Execute query
	rows, err := db.QueryContext(ctx, tsqlQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var houses []model.HouseTable
	for rows.Next() {
		var house model.HouseTable
		rows.Scan(&house.HouseId, &house.Price, &house.HouseLocation, &house.Distance, &house.UserId)
		houses = append(houses, house)
	}
	
	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(houses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ReadUserTable(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Verify database is running
	err := db.PingContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tsqlQuery := "SELECT UserId, Username, Name, Password FROM Users;"

	// Execute query
	rows, err := db.QueryContext(ctx, tsqlQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var users []model.UserTable
	for rows.Next() {
		var user model.UserTable
		rows.Scan(&user.UserId, &user.Username, &user.Name, &user.Password)
		users = append(users, user)
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Verify Database is running
	err := db.PingContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var u model.UserTable

	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check for valid login
	tsqlQuery := fmt.Sprintf("SELECT UserId FROM Users WHERE Username='%s' AND Name='%s' AND Password='%s';", u.Username, u.Name, u.Password)

	row := db.QueryRowContext(ctx, tsqlQuery)

	var uId int32
	if err = row.Scan(&uId); err != nil {
		http.Error(w, "No Login Found. Please register an account!", http.StatusUnauthorized)
	}

	// query session with user id, check if isActive is 1. if that's true then you get a single row then you send a message saying you are already logged in
	tsqlQuery = fmt.Sprintf("SELECT SessionId, IsActive FROM Session WHERE UserId='%d';", uId)
	aActiveRows, err := db.QueryContext(ctx, tsqlQuery)

	defer aActiveRows.Close()

	var aSess []model.Session
	for aActiveRows.Next() {
		var aS model.Session
		aActiveRows.Scan(&aS.SessionId, &aS.IsActive)
		if aS.IsActive == true {
			http.Error(w, "You are already logged in!", http.StatusForbidden)
			return
		}
		aSess = append(aSess, aS)
	}

	// Log user & create session
	guid := uuid.New()
	tsqlQuery = fmt.Sprintf("INSERT INTO Session VALUES(%d, '%s', 1)", uId, guid)
	result, err := db.ExecContext(ctx, tsqlQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	count, err := result.RowsAffected()
	if err != nil || count != 1 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := model.JsonLoginResponse{ Message: "Logged In", Type: "Success", UserGuid: guid.String() }
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Verify database is running
	err := db.PingContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if there is a valid guid with the current session
	var sess model.Session
	err = json.NewDecoder(r.Body).Decode(&sess)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if sess.UserGuid == "" {
		http.Error(w, "No User Guid provided with the given session.", http.StatusBadRequest)
		return
	}

	tsqlQuery := fmt.Sprintf("UPDATE Session SET IsActive=0 WHERE UserGuid='%s';", sess.UserGuid)

	// Execute query
	result, err := db.ExecContext(ctx, tsqlQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	count, err := result.RowsAffected()
	if err != nil || count != 1 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := model.GenericJsonResponse{ Message: "Successfully logged out of account!", Type: "Success" }
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Verify database is running
	err := db.PingContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var u model.UserTable

	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check user info isn't in User Table
	tsqlQuery := fmt.Sprintf("SELECT UserId FROM Users WHERE Username='%s' AND Name='%s' AND Password='%s';", u.Username, u.Name, u.Password)

	row := db.QueryRowContext(ctx, tsqlQuery)

	var uId int32
	if err = row.Scan(&uId); err == nil {
		http.Error(w, "Login information found. Please log into your account!", http.StatusUnauthorized)
		return
	}

	// Create User in User Table
	tsqlQuery = fmt.Sprintf("INSERT INTO Users VALUES('%s', '%s', '%s');", u.Username, u.Name, u.Password)

	result, err := db.ExecContext(ctx, tsqlQuery)
	if err != nil {
		http.Error(w, "Trouble creating your account. Please try again later!", http.StatusInternalServerError)
		return
	}

	count, err := result.RowsAffected()
	if err != nil || count != 1 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := model.GenericJsonResponse{ Message: "Successfully registered your account!", Type: "Success" }
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HouseFavorites(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Verify database is running
	err := db.PingContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Grab all houses in House Table
	tsqlQuery := "SELECT HouseId, Price, HouseLocation, Distance, UserId FROM House;"
	hRows, err := db.QueryContext(ctx, tsqlQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer hRows.Close()

	// Decode json body
	var houseF model.HouseFavorite
	err = json.NewDecoder(r.Body).Decode(&houseF)

	// Grab houses attached to passed in userId
	tsqlQuery = fmt.Sprintf("SELECT HouseId, Price, HouseLocation, Distance, UserId FROM House WHERE UserId=%d;", houseF.UserId)
	favoriteRows, err := db.QueryContext(ctx, tsqlQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer favoriteRows.Close()

	// Create favorites model
	var favorites []model.HouseTable
	for favoriteRows.Next() {
		var favorite model.HouseTable
		favoriteRows.Scan(&favorite.HouseId, &favorite.Price, &favorite.HouseLocation, &favorite.Distance, &favorite.UserId)
		favorites = append(favorites, favorite)
	}

	if len(favorites) == 0 {
		http.Error(w, "No bookmarks attached to your account!", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := model.HouseJsonResponse{ Message: "", Type: "Success", Data: favorites }
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdateFavorite(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Verify database is running
	err := db.PingContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Decode json body
	var houseF model.HouseFavorite

	err = json.NewDecoder(r.Body).Decode(&houseF)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Decide if house needs to be added or deleted to House Table
	if houseF.Favorite == false {
		// delete the current house from the house table
		tsqlQuery := fmt.Sprintf("DELETE FROM House WHERE UserId=%d AND HouseLocation='%s';", houseF.UserId, houseF.HouseLocation)
		dRows, err := db.ExecContext(ctx, tsqlQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		count, err := dRows.RowsAffected()
		if err != nil || count != 1 {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if houseF.Favorite == true {
		// add the current house to the House Table
		tsqlQuery := fmt.Sprintf("INSERT INTO House VALUES(%f, '%s', %f, %d);", houseF.Price, houseF.HouseLocation, houseF.Distance, houseF.UserId)
		aRows, err := db.ExecContext(ctx, tsqlQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		count, err := aRows.RowsAffected()
		if err != nil || count != 1 {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Favorite value not sent in request.", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := model.GenericJsonResponse{ Message: "Successfully updated your bookmark!", Type: "Success" }
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
