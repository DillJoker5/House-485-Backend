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
	houseCount, houseErr := readHouseTable()
	if houseErr != nil {
		log.Fatal("Error reading house table:", houseErr.Error())
	}
	fmt.Printf("Read %d row(s) successfully", houseCount)

	// Read from User Table
	userCount, userErr := readUserTable()
	if userErr != nil {
		log.Fatal("Error reading user table:", userErr.Error())
	}
	fmt.Printf("Read %d row(s) successfully", userCount)

	// Create new House in House Table
	newHouse, newHouseErr := createNewHouse()
	if newHouseErr != nil {
		log.Fatal("Error creating new house:", newHouseErr.Error())
	}
	fmt.Printf("Created new house successfully", newHouse)

	// Create new User in User Table
	newUser, newUserErr := createNewUser()
	if newUserErr != nil {
		log.Fatal("Error creating new user:", newUserErr.Error())
	}
	fmt.Printf("Created new user successfully", newUser)
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

	tsqlQuery := fmt.Sprintf("SELECT * FROM House")

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
	ctx := context.Background()

	// Verify database is running
	err := db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsqlQuery := fmt.Sprintf("SELECT * FROM Users")

	// Execute query
	rows, err := db.QueryContext(ctx, tsqlQuery)
	if err != nil {
		return -1, err
	}

	defer rows.Close()

	var count int
	for rows.Next() {
		var userid int32
		var username string
		var name string
		var password string
		var houseid int32

		err := rows.Scan(&userid, &username, &name, &password, &houseid)
		if err != nil {
			return -1, err
		}
		// do work here
		count++
	}
	return count, nil
}

func createNewHouse() (model.HouseTable, error) {
	ctx := context.Background()

	// Verify database is running
	err := db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	tsqlMutation := fmt.Sprintf("INSERT INTO House VALUES()") // finish mutation

	// Execute query
	newHouse, err := db.QueryContext(ctx, tsqlMutation)
	if err != nil {
		return nil, err
	}

	defer newHouse.Close()

	// finish here
}

func createNewUser() (model.UserTable, error) {
	ctx := context.Background()

	// Verify database is running
	err := db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	tsqlMutation := fmt.Sprintf("INSERT INTO Users VALUES()") // finish mutation

	// Execute query
	newUser, err := db.QueryContext(ctx, tsqlMutation)
	if err != nil {
		return nil, err
	}

	defer newUser.Close()

	// finish here
}
