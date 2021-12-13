package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	handlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type passengerinfo struct {
	PassengerID  int    `json:"passengerid"`
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	MobileNumber string `json:"mobilenumber"`
	EmailAddress string `json:"emailaddress"`
}

//var passenger map[string]passengerinfo

func validKey(r *http.Request) bool {
	v := r.URL.Query()
	if key, ok := v["key"]; ok {
		if key[0] == "2c78afaf-97da-4816-bbee-9ad239abb296" {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Rideshare passenger API")
}

//function to check if passenger exists
func CheckPassenger(db *sql.DB, email string) bool {
	query := fmt.Sprintf("Select * FROM Passenger WHERE EmailAddress= '%s'", email)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var passenger passengerinfo
	for results.Next() {
		err = results.Scan(&passenger.PassengerID, &passenger.FirstName,
			&passenger.LastName, &passenger.MobileNumber, &passenger.EmailAddress)
		if err != nil {
			panic(err.Error())
		} else if passenger.EmailAddress == email {
			return true
		}
	}
	return false
}

//funtion for getting passenger
func GetPassenger(db *sql.DB, emailAddr string) passengerinfo {
	query := fmt.Sprintf("Select * FROM Passenger WHERE EmailAddress='%s'", emailAddr)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var passenger passengerinfo
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&passenger.PassengerID, &passenger.FirstName,
			&passenger.LastName, &passenger.MobileNumber, &passenger.EmailAddress)
		if err != nil {
			panic(err.Error())
		}
	}
	return passenger
}

//function for creating passenger
func CreatePassenger(db *sql.DB, passenger passengerinfo) bool {
	query := fmt.Sprintf("INSERT INTO Passenger(PassengerID, FirstName, LastName, MobileNumber, EmailAddress) VALUES (%d,'%s', '%s', '%s','%s')",
		getlastid(db)+1,
		passenger.FirstName,
		passenger.LastName,
		passenger.MobileNumber,
		passenger.EmailAddress)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	return true
}

//function for editing passenger
func EditPassenger(db *sql.DB, passenger passengerinfo) bool {
	query := fmt.Sprintf(
		"UPDATE Passenger SET FirstName='%s', LastName='%s', MobileNumber='%s', EmailAddress='%s' WHERE PassengerID=%d",
		passenger.FirstName, passenger.LastName, passenger.MobileNumber, passenger.EmailAddress, passenger.PassengerID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

//gets the last passengerid from db
func getlastid(db *sql.DB) int { //Gets the last id of passengers
	query1 := "SELECT COUNT(*) FROM Passenger"
	query2 := "SELECT PassengerID FROM Passenger ORDER BY PassengerID DESC LIMIT 1"
	var passengerCount int
	results, err := db.Query(query1) //Run Query
	if err != nil {
		panic(err.Error())
	}
	if results.Next() {
		results.Scan(&passengerCount)
	}
	if passengerCount > 0 {
		results, err := db.Query(query2) //Run Query
		var ID int
		if err != nil {
			panic(err.Error())
		}
		if results.Next() {
			results.Scan(&ID)
		}
		return ID
	} else {
		return 0
	}

}
func passengers(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}
	db, err := sql.Open("mysql", "root:12N28c02@tcp(127.0.0.1:3306)/assignment_db") //connect to database
	if err != nil {
		fmt.Println(err)
	}
	params := mux.Vars(r)
	if r.Method == "DELETE" {
		println("Can't delete passenger")
	} else if r.Method == "GET" {
		if params["emailaddress"] == " " {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Please provide email address"))
			return
		}
		//GET passenger using email address
		PassengerInformation := GetPassenger(db, params["emailaddress"])
		if err != nil {
			fmt.Printf("There was an error encoding the json. err = %s", err)
		} else if PassengerInformation.EmailAddress == "" { // Check if data is empty
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Passenger does not exists"))
			return
		} else {
			json.NewEncoder(w).Encode(GetPassenger(db, PassengerInformation.EmailAddress))
			w.WriteHeader(http.StatusAccepted)
			return
		}
	}

	if r.Header.Get("Content-type") == "application/json" {

		// POST is for creating new passenger
		if r.Method == "POST" {

			// read the string sent to the service
			var newPassenger passengerinfo
			reqBody, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newPassenger)

				if newPassenger.EmailAddress == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please enter the required information " +
							"in JSON format"))
					return
				} else {
					// check if user already exists by email; add only if
					// user does not exist
					if !CheckPassenger(db, newPassenger.EmailAddress) {
						CreatePassenger(db, newPassenger)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("Passenger created successfully"))
						return
					} else {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("Email is already in use"))
						return
					}
				}
			}
		} else if r.Method == "PUT" {
			//---PUT is for creating or updating
			// existing passenger---
			var passenger passengerinfo
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &passenger)

				if passenger.EmailAddress == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply passenger " +
							" information " +
							"in JSON format"))
					return
				} else {
					// check if passenger does not exists; update only if
					// passenger does exist
					if !CheckPassenger(db, passenger.EmailAddress) {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("There is no exsiting passenger with " + passenger.EmailAddress))
						return
					} else {
						//To update user details
						EditPassenger(db, passenger)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("Passenger updated successfully"))
						return
					}

				}
			}
		}

	}
}
func GetPassengerID(db *sql.DB, email string) int {
	query := fmt.Sprintf("Select * FROM Customer WHERE EmailAddress= '%s'", email)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var passenger passengerinfo
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&passenger.PassengerID, &passenger.FirstName,
			&passenger.LastName, &passenger.MobileNumber, &passenger.EmailAddress)
		if err != nil {
			panic(err.Error())
		}
	}
	return passenger.PassengerID
}
func CheckPassengerEmail(w http.ResponseWriter, r *http.Request) {
	//Database
	db, err := sql.Open("mysql", "root:12N28c02@tcp(127.0.0.1:3306)/assignment_db") //Connecting to database
	if err != nil {
		fmt.Println(err)
	}
	params := mux.Vars(r)
	if params["UserEmail"] == "" {
		w.WriteHeader(
			http.StatusUnprocessableEntity)
		w.Write([]byte(
			"422 - Please supply email information"))
		return
	} else if CheckPassenger(db, params["UserEmail"]) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(strconv.Itoa(GetPassengerID(db, params["UserEmail"]))))
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}

func main() {
	// instantiate courses
	//passenger = make(map[string]passengerinfo)

	router := mux.NewRouter()
	//Web Front-end CORS
	headers := handlers.AllowedHeaders([]string{"X-REQUESTED-With", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/api/v1/passenger", home)
	router.HandleFunc("/api/v1/passenger/CheckPassenger/{PassengerEmail}", CheckPassengerEmail)
	router.HandleFunc("/api/v1/passenger/router/{emailaddress}", passengers).Methods(
		"GET", "PUT", "POST")
	fmt.Println("Listening at port 5001")
	log.Fatal(http.ListenAndServe(":5001", handlers.CORS(headers, methods, origins)(router)))
}
