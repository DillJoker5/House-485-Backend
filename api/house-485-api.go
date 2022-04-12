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

func readHouseTable() (int, error) {
	ctx := context.Background()

	// Verify database is running
	err := db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsqlQuery := fmt.Sprintf("SELECT * FROM HouseTable")

	// Execute query
	rows, err := db.QueryContext(ctx, tsqlQuery)
	if err != nil {
		return -1, err
	}

	defer rows.Close()

	var count int
	for rows.Next() {
		var houseid int32
		var price float32
		var houseLocation string
		var distance float32

		err := rows.Scan(&houseid, &price, &houseLocation, &distance)
		if err != nil {
			return -1, err
		}
		// do work here
		count++
	}
	return count, nil
}

func readUserTable() (int, error) {
	return -1, nil
}

func createNewHouse() (model.HouseTable, error) {}

func createNewUser() (model.UserTable, error) {}
