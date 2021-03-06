/*
	Written by Dylan Chirigotis
	All API endpoints and backend logic for the House-485-Website
*/

// Main Package
package main

// All api imports
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

// Global database variable
var db *sql.DB

// Database connection parameters
var server = "DESKTOP-K7IIMGF"
var port = 1433
var user = "capstoneapiuser"
var password = "CapstoneApiUser2022!"
var database = "House485Database"

/*
	Main Function
	
	This function creates the connection to the database, creates and hosts the server, creates all
	API endpoints, and handles all of the API endpoints.
*/
func main() {
	// Build connection string
	connectionString  := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", server, user, password, port, database)
	var err error

	// Create connection
	db, err = sql.Open("sqlserver", connectionString)
	if err != nil {
		log.Fatal("Error connecting:", err.Error())
	}

	// Create router
	router := mux.NewRouter()

	// Handle api requests without middleware check
	router.HandleFunc("/readUser", ReadUserTable).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/home", ReadHouseTable).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/login", Login).Methods(http.MethodPost)
	router.HandleFunc("/register", Register).Methods(http.MethodPost)

	// Handle api requests with middleware check
	router.HandleFunc("/logout", mwCheck(Logout)).Methods(http.MethodPost)
	router.HandleFunc("/favorite", mwCheck(HouseFavorites)).Methods(http.MethodPost)
	router.HandleFunc("/updateFavorite", mwCheck(UpdateFavorite)).Methods(http.MethodPost)

	// Create server with router and localhost:8000 address
	srv := &http.Server {
		Addr: ":8000",
		Handler: router,
	}

	// Host and start the server
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

/*
	mwCheck Function

	This function makes sure that the user is validated before running the passed in function.
	
	If the user is unauthorized, it returns an error detailing that the user isn't authorized.
	
	If the user is authorized, it calls the passed in function.
*/

func mwCheck(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	// Return passed in function
	return func(w http.ResponseWriter, r *http.Request) {
		// If user is not authorized, throw unauthorized error
		if !validateUser(r) {
			http.Error(w, "Unauthorized user", http.StatusForbidden)
		} else {
			// Call function since the user is authorized
			f(w, r)
		}
	}
}

/*
	validateUser Function

	This function checks if the UserGuid header is passed into the request.
	
	If the UserGuid is empty, not passed in, or a session was not created with the passed
	in UserGuid, it will return false signaling an unauthorized user.
	
	If there is an active session with the UserGuid, it will return true.
*/

func validateUser(r *http.Request) bool {
	ctx := context.Background()
	
	// Verify database is running
	err := db.PingContext(ctx)
	if err != nil {
		return false
	}

	// Grab userGuid from header
	uGuid := r.Header.Get("userguid")
	if uGuid == "" {
		return false
	}

	tsqlQuery := fmt.Sprintf("SELECT SessionId FROM Session WHERE UserGuid='%s' AND IsActive=1;", uGuid)

	// Execute Query
	row := db.QueryRowContext(ctx, tsqlQuery)
	if err != nil {
		return false
	}

	// Scan for the userId
	var sid int32
	if err = row.Scan(&sid); err != nil {
		return false
	}

	return true
}

/*
	ReadHouseTable Function
	
	This function reads all of the houses from the database and returns them.
	
	If there is any error in connecting to the database, executing the query,
	or returning the response, this function will return an error.
	
	If no error is run into, this function will return an array of all of the
	houses in the database.
*/

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

	// Close rows to prevent a memory leak
	defer rows.Close()

	// Create array of houses and append each scanned house to the array
	var houses []model.HouseTable
	for rows.Next() {
		var house model.HouseTable
		rows.Scan(&house.HouseId, &house.Price, &house.HouseLocation, &house.Distance, &house.UserId)
		houses = append(houses, house)
	}
	
	// Add in additional headers for localhost and okay status
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
	w.WriteHeader(http.StatusOK)

	// Return response
	err = json.NewEncoder(w).Encode(houses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
	ReadUserTable Function
	
	This function reads all of the users from the database and returns them.
	
	If there is any error in connecting to the database, executing the query,
	or returning the response, this function will return an error.
	
	If no error is run into, this function will return an array of all of the
	users in the database.
*/

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

	// Close rows to prevent a memory leak
	defer rows.Close()

	// Create array of users and append each scanned user to the array
	var users []model.UserTable
	for rows.Next() {
		var user model.UserTable
		rows.Scan(&user.UserId, &user.Username, &user.Name, &user.Password)
		users = append(users, user)
	}

	// Add in additional headers for localhost and okay status
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
	w.WriteHeader(http.StatusOK)

	// Return response
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
	Login Function

	This function reads a user's username and password and returns the UserGuid and UserId
	upon a successful login.
	
	If there is any error in connecting to the database, decoding the request body, executing
	the query, no login information was found in the users table, the user is already logged
	in, creating the session, and returning the response, this function will return an error.

	If no error is run into, this function will return a UserGuid and UserId to the front-end.
*/

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Verify Database is running
	err := db.PingContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var u model.User

	// Decode request body into the variable u
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check for valid login
	tsqlQuery := fmt.Sprintf("SELECT UserId FROM Users WHERE Username='%s' AND Password='%s';", u.Username, u.Password)

	// Execute query
	row := db.QueryRowContext(ctx, tsqlQuery)

	// If no userid is found, throw error and tell the user to register an account
	var uId int32
	if err = row.Scan(&uId); err != nil {
		http.Error(w, "No Login Found. Please register an account!", http.StatusUnauthorized)
		return
	}

	tsqlQuery = fmt.Sprintf("SELECT SessionId FROM Session WHERE UserId='%d' AND IsActive=1;", uId)
	
	// Execute query
	aActiveRow := db.QueryRowContext(ctx, tsqlQuery)

	// Check is user is already logged in
	var sId int32
	err = aActiveRow.Scan(&sId)
	if sId > 0 {
		http.Error(w, "You are already logged into your account!", http.StatusForbidden)
		return
	}

	// Login user & create session
	guid := uuid.New()
	tsqlQuery = fmt.Sprintf("INSERT INTO Session VALUES(%d, '%s', 1)", uId, guid)
	result, err := db.ExecContext(ctx, tsqlQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify that the session was created
	count, err := result.RowsAffected()
	if err != nil || count != 1 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add in additional headers for localhost and okay status
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
	w.WriteHeader(http.StatusOK)
	
	// Return response
	response := model.JsonLoginResponse{ Message: "Logged In", Type: "Success", UserGuid: guid.String(), UserId: uId }
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
	Logout Function

	This function will logout the passed in user by setting their active status to false.

	If there is any error in connecting to the database, decoding the request body, no
	UserGuid is present, updating the session in the sessions table, and returning the
	response, this function will return an error.

	If no error is run into, this function will return a success message to the front-end.
*/

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
	
	// Make sure that the UserGuid is not empty
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

	// Verify that the session was deleted
	count, err := result.RowsAffected()
	if err != nil || count != 1 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add in additional headers for localhost and okay status
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
	w.WriteHeader(http.StatusOK)

	// Return response
	response := model.GenericJsonResponse{ Message: "Successfully logged out of account!", Type: "Success" }
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
	Register Function

	This function will register a user by creating a new user in the Users table.

	If there is any error in connecting to the database, decoding the request body,
	a user trying to register an already existing account, creating the user in the Users
	table, and returing the response, this function will return an error.

	If no error is run into, this function will return a success message to the front-end.
*/

func Register(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Verify database is running
	err := db.PingContext(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var u model.RegisterUser

	// Decode request body
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check user info isn't in User Table
	tsqlQuery := fmt.Sprintf("SELECT UserId FROM Users WHERE Username='%s' AND Name='%s' AND Password='%s';", u.Username, u.Name, u.Password)

	row := db.QueryRowContext(ctx, tsqlQuery)

	// The user account already exists
	var uId int32
	if err = row.Scan(&uId); err == nil {
		http.Error(w, "Login information found. Please log into your account!", http.StatusUnauthorized)
		return
	}

	// Create User in User Table
	tsqlQuery = fmt.Sprintf("INSERT INTO Users VALUES('%s', '%s', '%s');", u.Username, u.Name, u.Password)

	// Error in creating the user
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

	// Add in additional headers for localhost and okay status
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
	w.WriteHeader(http.StatusOK)

	// Return response
	response := model.GenericJsonResponse{ Message: "Successfully registered your account!", Type: "Success" }
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
	HouseFavorites Function

	This function will return all houses that are associated with a passed in userId.

	If there is any error in connecting to the database, grabbing all of the houses from
	the house table, decoding the request body, no houses are associated with that account, and returing the response, this
	function will return an error.

	If no error is run into, this function will return an array of all of the houses
	associated with the request's userId.
*/

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
	
	// Execute query
	hRows, err := db.QueryContext(ctx, tsqlQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Close rows to prevent a memory leak
	defer hRows.Close()

	// Decode json body
	var houseF model.HouseFavorite
	err = json.NewDecoder(r.Body).Decode(&houseF)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Grab houses attached to passed in userId
	tsqlQuery = fmt.Sprintf("SELECT HouseId, Price, HouseLocation, Distance, UserId FROM House WHERE UserId=%d;", houseF.UserId)
	
	// Execute query
	favoriteRows, err := db.QueryContext(ctx, tsqlQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Close favoriteRows to prevent a memory leak
	defer favoriteRows.Close()

	// Create favorites model
	var favorites []model.HouseTable
	for favoriteRows.Next() {
		var favorite model.HouseTable
		favoriteRows.Scan(&favorite.HouseId, &favorite.Price, &favorite.HouseLocation, &favorite.Distance, &favorite.UserId)
		favorites = append(favorites, favorite)
	}

	// Throw error if the favorites model is empty
	if len(favorites) == 0 {
		http.Error(w, "No bookmarks attached to your account!", http.StatusInternalServerError)
		return
	}

	// Add in additional headers for localhost and okay status
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
	w.WriteHeader(http.StatusOK)
	
	// Return response
	response := model.HouseJsonResponse{ Message: "", Type: "Success", Data: favorites }
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
	UpdateFavorite Function

	This function will either create or delete a house that an associated user wants to bookmark.

	If there is any error in connecting to the database, decoding the request body, if deleting the
	given house from the house table malfunctions, if adding the given house to the house table malfunctions
	, a favorite value was not sent in the request body, and returning the response, this function
	will return an error.

	If no error is run into, this function will do one of two things: it will create the house in the
	house table that is associated with the currently logged in user or it will delete the house in
	the house table that is associated with the currently logged in user.
*/

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
		// Delete the current house from the house table
		tsqlQuery := fmt.Sprintf("DELETE FROM House WHERE UserId=%d AND HouseLocation='%s';", houseF.UserId, houseF.HouseLocation)
		
		// Execute query
		dRows, err := db.ExecContext(ctx, tsqlQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		// Verify that the house was deleted from the house table
		count, err := dRows.RowsAffected()
		if err != nil || count != 1 {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if houseF.Favorite == true {
		// Add the current house to the House Table
		tsqlQuery := fmt.Sprintf("INSERT INTO House VALUES(%f, '%s', %f, %d);", houseF.Price, houseF.HouseLocation, houseF.Distance, houseF.UserId)
		
		// Execute query
		aRows, err := db.ExecContext(ctx, tsqlQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Verify that the house was added to the house table
		count, err := aRows.RowsAffected()
		if err != nil || count != 1 {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Throw an error due to the favorite value not being sent in the request body
		http.Error(w, "Favorite value not sent in request.", http.StatusInternalServerError)
		return
	}

	// Add in additional headers for localhost and okay status
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
	w.WriteHeader(http.StatusOK)
	
	// Return response
	response := model.GenericJsonResponse{ Message: "Successfully updated your bookmark!", Type: "Success" }
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
