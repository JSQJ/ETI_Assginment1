package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	handlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type driverinfo struct {
	DriverID      string `json:"driverid"`
	FirstName     string `json:"firstname"`
	LastName      string `json:"lastname"`
	MobileNumber  string `json:"mobilenumber"`
	EmailAddress  string `json:"emailaddress"`
	LicenseNumber string `json:"licensenumber"`
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "RideShare Driver API")
}
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

//var driver map[string]driverinfo
//function to check if driver exists
func CheckDriver(db *sql.DB, email string) bool {
	query := fmt.Sprintf("Select * FROM Driver WHERE EmailAddress= '%s'", email)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var driver driverinfo
	for results.Next() {
		err = results.Scan(&driver.DriverID, &driver.FirstName,
			&driver.LastName, &driver.MobileNumber, &driver.EmailAddress, &driver.LicenseNumber)
		if err != nil {
			panic(err.Error())
		} else if driver.EmailAddress == email {
			return true
		}
	}
	return false
}

func GetDriver(db *sql.DB, email string) driverinfo {
	query := fmt.Sprintf("Select * FROM Driver WHERE EmailAddress= '%s'", email)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var driver driverinfo
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&driver.DriverID, &driver.FirstName,
			&driver.LastName, &driver.MobileNumber, &driver.EmailAddress, &driver.LicenseNumber)
		if err != nil {
			panic(err.Error())
		}
	}
	fmt.Println(driver.DriverID, driver.FirstName,
		driver.LastName, driver.MobileNumber, driver.EmailAddress, driver.LicenseNumber)
	return driver
}

func CreateDriver(db *sql.DB, driver driverinfo) bool {
	query := fmt.Sprintf("INSERT INTO Driver VALUES ('%s', '%s', '%s', '%s','%s','%s')",
		driver.DriverID, driver.FirstName, driver.LastName, driver.MobileNumber, driver.EmailAddress, driver.LicenseNumber)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())

	}
	return true
}
func EditDriver(db *sql.DB, driver driverinfo) bool {
	if driver.DriverID == "" {
		return false
	}
	query := fmt.Sprintf("UPDATE Driver SET FirstName = '%s', LastName = '%s', MobileNumber= '%s', LicenseNumber = '%s', EmailAddress = '%s' WHERE DriverID = '%s';",
		driver.FirstName, driver.LastName, driver.MobileNumber, driver.LicenseNumber, driver.EmailAddress, driver.DriverID)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	return true
}

func SearchAvailDriver(db *sql.DB) string {
	//GET All driver ID
	query := "SELECT DriverID FROM driver"
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var DriverIDs string
	for results.Next() {
		var TempID string
		err = results.Scan(&TempID)
		if err != nil {
			panic(err.Error())
		}
		DriverIDs += TempID + ","
	}
	return DriverIDs
}
func GetAllDrivers(w http.ResponseWriter, r *http.Request) {
	//Database
	db, err := sql.Open("mysql", "root:12N28c02@tcp(127.0.0.1:3306)/assignment_db") //Connecting to database
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(SearchAvailDriver(db)))
}
func drivers(w http.ResponseWriter, r *http.Request) {
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
		println("Can't delete Driver")
	} else if r.Method == "GET" {
		if params["emailaddress"] == " " {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Please provide email address"))
			return
		}
		//GET Driver using email address
		DriverInformation := GetDriver(db, params["emailaddress"])
		if err != nil {
			fmt.Printf("There was an error encoding the json. err = %s", err)
		} else if DriverInformation.EmailAddress == "" { // Check if data is empty
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Driver does not exists"))
			return
		} else {
			json.NewEncoder(w).Encode(GetDriver(db, DriverInformation.EmailAddress))
			w.WriteHeader(http.StatusAccepted)
			return
		}
	}
	if r.Header.Get("Content-type") == "application/json" {
		if err != nil {
			fmt.Println(err)
		}
		// POST is for creating new course
		if r.Method == "POST" {

			var newDriver driverinfo
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newDriver)

				if newDriver.EmailAddress == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please enter the required information " +
							"in JSON format"))
					return
				} else {
					// check if driver already exists by email; add only if
					// driver does not exist
					if !CheckDriver(db, newDriver.EmailAddress) {
						CreateDriver(db, newDriver)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("Driver created successfully"))
						return
					} else {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("Email is already in use"))
						return
					}
				}
			}
		}
		//---PUT is for creating or updating
		// existing course---
		if r.Method == "PUT" {
			//---PUT is for creating or updating
			// existing driver---
			var driver driverinfo
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &driver)

				if driver.EmailAddress == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply driver " +
							" information " +
							"in JSON format"))
					return
				} else {
					// check if Driver does not exists; update only if
					// driver does exist
					if !CheckDriver(db, driver.EmailAddress) {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("There is no exsiting driver with " + driver.EmailAddress))
						return
					} else {
						//To update driver details
						EditDriver(db, driver)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("Driver updated successfully"))
						return
					}

				}
			}
		}

	}
}

func GetDriverByID(w http.ResponseWriter, r *http.Request) {
	//Database
	db, err := sql.Open("mysql", "root:12N28c02@tcp(127.0.0.1:3306)/assignment_db") //Connecting to database
	if err != nil {
		fmt.Println(err)
	}
	params := mux.Vars(r)
	if params["DriverID"] == "" {
		w.WriteHeader(
			http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply Driver ID"))
		return
	} else {
		println(params["DriverID"])
		query := fmt.Sprintf("Select LicenseNumber FROM Driver WHERE DriverID= '%s'", params["DriverID"])
		results, err := db.Query(query)
		if err != nil {
			panic(err.Error())
		}
		var DriverNumPlate string
		for results.Next() {
			// map this type to the record in the table
			err = results.Scan(&DriverNumPlate)
			if err != nil {
				panic(err.Error())
			}
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(DriverNumPlate))
		return
	}
}
func main() {
	router := mux.NewRouter()
	//Web Front-end CORS
	headers := handlers.AllowedHeaders([]string{"X-REQUESTED-With", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/api/v1/driver", home)
	router.HandleFunc("/api/v1/driver/GetLicensePlate/{driverid}", GetDriverByID)
	router.HandleFunc("/api/v1/driver/getalldrivers", GetAllDrivers)
	router.HandleFunc("/api/v1/driver/router/{emailaddress}", drivers).Methods(
		"GET", "PUT", "POST")
	fmt.Println("Listening at port 5002")
	log.Fatal(http.ListenAndServe(":5002", handlers.CORS(headers, methods, origins)(router)))
}
