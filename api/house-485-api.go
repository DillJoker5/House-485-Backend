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
var server = ""
var port = 1433
var user = ""
var password = ""
var database = ""

func main() {
	// Build connection string
	connectionString  := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", server, user, password, port, database)
	var err error

	// Create connection
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		log.Fatal("Error connecting:", err.Error())
	}

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected to the database!\n")

	// Read from House Table

	// Read from User Table

	// Create new House in House Table

	// Create new User in User Table
}

func mwCheck(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func readHouseTable() (int, error) {}

func readUserTable() (int, error) {}

func createNewHouse() (model.HouseTable, error) {}

func createNewUser() (model.UserTable, error) {}
