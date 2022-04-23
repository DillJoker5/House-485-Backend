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
	router.HandleFunc("/readUser", mwCheck(ReadUserTable)).Methods(http.MethodPost)
	router.HandleFunc("/home", mwCheck(ReadHouseTable)).Methods(http.MethodPost)
	router.HandleFunc("/login", mwCheck(Login)).Methods(http.MethodPost)
	router.HandleFunc("/logout", mwCheck(Logout)).Methods(http.MethodPost)
	router.HandleFunc("/register", Register).Methods(http.MethodPost)
	//router.HandleFunc("/favorite", mwCheck(HouseFavorites)).Methods(http.MethodPost)

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

	tsqlQuery := fmt.Sprintf("SELECT SessionId FROM Sessions WHERE UserGuid='%s' AND IsActive=1;", s.UserGuid)

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

	tsqlQuery := "SELECT HouseId, Price, HouseLocation, Distance FROM House;"

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
		rows.Scan(&house.HouseId, &house.Price, &house.HouseLocation, &house.Distance)
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

	tsqlQuery := "SELECT UserId, Username, Name, Password, HouseId FROM Users;"

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
		rows.Scan(&user.UserId, &user.Username, &user.Name, &user.Password, &user.HouseId)
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
	tsqlQuery = fmt.Sprintf("INSERT INTO Session VALUES(%d, '%s', 1)", uId, guid) // finish
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
	tsqlQuery = fmt.Sprintf("INSERT INTO Users VALUES('%s', '%s', '%s', %d);", u.Username, u.Name, u.Password, u.HouseId)

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

/*func HouseFavorites(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Verify database is running
	err := db.PingContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Grab all houses in House Table
	tsqlQuery := "SELECT HouseId, Address, Location, Distance FROM House;"
	hRows, err := db.QueryContext(ctx, tsqlQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer hRows.Close()

	// finish

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := model.HouseJsonResponse{ Message: "", Type: "Success", Data: [] }
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}*/
