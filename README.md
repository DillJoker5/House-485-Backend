# House-485-Backend
Complete api backend for the House-485-Website. Backend/APIs are completely built off of Go.

## Description
This repository holds all of the backend logic, api endpoints, and route endpoints for the [House-485-Website](https://github.com/DillJoker5/House-485-Website).

## Required Imports
* `encoding/json`
* `log`
* `net/http`
* `house-485-backend/models`
* `github.com/gorilla/mux`
* `database/sql`
* `context`
* `github.com/denisenkom/go-mssqldb`
* `fmt`
* `github.com/google/uuid`

## Installation
You will need the following things installed and working on your device before running this project.
* Please note that you will need to install and run both the [House-485-Website](https://github.com/DillJoker5/House-485-Website) and [House-485-Database](https://github.com/DillJoker5/House-485-Database) to fully use this backend project. Please see those installation guides for how to install them.
* [Go](https://go.dev/doc/install)
* Your favorite code editor for Go
* Optional: Postman

## Running
* Create a folder where you will store this project
* Clone the repository in the newly created folder
* Run the [House-485-Database](https://github.com/DillJoker5/House-485-Database) and make sure that you are able to log in with the user. See [here](https://github.com/DillJoker5/House-485-Database) for help
* Open a new terminal in the directory of the repository
* Before running the repository, run the following command: go mod tidy. This command will install all of the required imports for this repository.
* In the same terminal, run the following command to start the server: go run api/house-485-api.go
* Optional: test an api endpoint in Postman using the URL https://localhost:8000/

## Code Examples
Code examples of functions, handlers, and other important information will be shown here.
### All route handlers

	router.HandleFunc("/readUser", ReadUserTable).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/home", ReadHouseTable).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/login", Login).Methods(http.MethodPost)
	router.HandleFunc("/register", Register).Methods(http.MethodPost)

	router.HandleFunc("/logout", mwCheck(Logout)).Methods(http.MethodPost)
	router.HandleFunc("/favorite", mwCheck(HouseFavorites)).Methods(http.MethodPost)
	router.HandleFunc("/updateFavorite", mwCheck(UpdateFavorite)).Methods(http.MethodPost)


### Register Function

	func Register(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		err := db.PingContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var u model.RegisterUser

		err = json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tsqlQuery := fmt.Sprintf("SELECT UserId FROM Users WHERE Username='%s' AND Name='%s' AND Password='%s';", u.Username, u.Name, u.Password)

		row := db.QueryRowContext(ctx, tsqlQuery)

		var uId int32
		if err = row.Scan(&uId); err == nil {
			http.Error(w, "Login information found. Please log into your account!", http.StatusUnauthorized)
			return
		}

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

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
		w.WriteHeader(http.StatusOK)
		response := model.GenericJsonResponse{ Message: "Successfully registered your account!", Type: "Success" }
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

### UpdateFavorite Function

	func UpdateFavorite(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		err := db.PingContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var houseF model.HouseFavorite

		err = json.NewDecoder(r.Body).Decode(&houseF)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if houseF.Favorite == false {
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

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
		w.WriteHeader(http.StatusOK)
		response := model.GenericJsonResponse{ Message: "Successfully updated your bookmark!", Type: "Success" }
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
