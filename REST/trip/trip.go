package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type tripinfo struct {
	TripID          int    `json:"tripid"`
	PassengerID     int    `json:"passengerid"`
	DriverID        string `json:"driverIid"`
	PickUpLocation  string `json:"pickuplocation"`
	DropOffLocation string `json:"dropofflocation"`
	PickUpTime      string `json:"pickuptime"`
	DropOffTime     string `json:"dropofftime"`
	Status          string `json:"status"`
}

type Data struct { //Structure to get neccasry input from front end
	PassengerEmail  string `json:"passengeremail"`
	PickUpLocation  string `json:"pickUplocation"`
	DropOffLocation string `json:"dropofflocation"`
	PickUpTime      string `json:"pickuptime"`
	DropOffTime     string `json:"dropofftime"`
}

type AllTrip struct {
	LicensePlate    string `json:"licenseplate"`
	PickUpLocation  string `json:"pickuplocation"`
	DropOffLocation string `json:"dropofflocation"`
	PickUpTime      string `json:"pickuptime"`
	DropOffTime     string `json:"dropofftime"`
	Status          string `json:"status"`
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "RideShare Trip API")
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

func CheckTrip(db *sql.DB, PassengerID int, DriverID string) bool {
	//To check if there is any uncompleted trips
	query := fmt.Sprintf("Select Status FROM Trips WHERE PassengerID = %d OR DriverID = '%s'", PassengerID, DriverID)
	results, err := db.Query(query)
	if err != nil {
		panic("Error here" + err.Error())
	}
	var currentTrip tripinfo
	for results.Next() {
		// map this type to the record in the table and update the object with new data
		err = results.Scan(&currentTrip.Status)
		if err != nil {
			panic(err.Error())
		} else if currentTrip.Status != "Completed" { //To check if trip have finish
			return false
		}
	}
	return true
}
func CheckDriverAvailability(db *sql.DB, IDs []string) string {
	//Check if driver is currently on a trip
	//QueryString := fmt.Sprintf("Select * FROM Trip WHERE Status != 'Completed' AND DriverID = '%s' ", IDs)
	QueryString := "Select DriverID FROM Trips WHERE Status != 'Completed' "

	for i := range IDs {
		if i == 0 {
			QueryString += fmt.Sprintf("AND DriverID = '%s' ", string(IDs[i]))
		} else {
			QueryString += fmt.Sprintf("OR DriverID = '%s' ", string(IDs[i]))
		}

	}
	results, err := db.Query(QueryString)
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		var DriverID string
		// map this type to the record in the table
		err = results.Scan(&DriverID)
		if err != nil {
			panic(err.Error())
		}
		IDs = remove(IDs, DriverID)
	}
	return IDs[0]
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func GetTrip(db *sql.DB, ID int) tripinfo {
	query := fmt.Sprintf("Select * FROM Trips WHERE TripID = %d", ID)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var trip tripinfo
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&trip.TripID, &trip.DriverID, &trip.PassengerID, &trip.PickUpLocation, &trip.DropOffLocation, &trip.PickUpTime, &trip.DropOffTime, &trip.Status)
		if err != nil {
			panic(err.Error())
		}
	}
	return trip
}
func GetTrips(db *sql.DB, PassengerID int) []tripinfo {
	query := fmt.Sprintf("Select * FROM Trips WHERE PassengerID = %d ORDER BY TripID DESC", PassengerID)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var trips []tripinfo
	for results.Next() {
		var trip tripinfo
		// map this type to the record in the table
		err = results.Scan(&trip.TripID, &trip.DriverID, &trip.PassengerID, &trip.PickUpLocation, &trip.DropOffLocation, &trip.PickUpTime, &trip.DropOffTime, &trip.Status)
		if err != nil {
			panic(err.Error())
		}
		trips = append(trips, trip)
	}
	return trips
}

func CreateTrip(db *sql.DB, trip tripinfo) bool {
	query := fmt.Sprintf(
		"INSERT INTO Trip (DriverID, PassengerID, PickupLocation, DropoffLocation, PickUpTime, DropOffTime, Status) VALUES('%s', %d, '%s', '%s', '%s', ' ', 'Pending')",
		trip.DriverID,
		trip.PassengerID,
		trip.PickUpLocation,
		trip.DropOffLocation,
		trip.PickUpTime)

	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	return true
}

func EditTrip(db *sql.DB, trip tripinfo) bool {
	if trip.TripID == 0 {
		return false
	}
	query := fmt.Sprintf("UPDATE Trips SET DriverID = '%s', PassengerID = %d, PickUpLocation = '%s', DropOffLocation='%s', PickUpTime = '%s', DropOffTime = '%s', Status = '%s' WHERE TripID = %d",
		trip.DriverID,
		trip.PassengerID,
		trip.PickUpLocation,
		trip.DropOffLocation,
		trip.PickUpTime,
		trip.DropOffTime,
		trip.Status,
		trip.TripID)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	return true
}

func GetAllDriver() []string {
	response, err := http.Get("http://localhost:5002/api/v1/driver/getalldrivers")
	if err != nil {
		fmt.Print(err.Error())
	}
	if response.StatusCode == http.StatusAccepted {
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			println(err)
		} else {
			IDs := strings.Split(string(responseData), ",")
			replacer := strings.NewReplacer(",", "")
			var newIDs []string
			for i := range IDs {
				newIDs = append(newIDs, replacer.Replace(IDs[i]))
			}
			return newIDs
		}
	}
	return nil
}
func CheckPassenger(Email string) int {
	//Check Passenger exsist in the database
	URL := "http://localhost:5001/api/v1/passenger/CheckPassenger/" + Email

	response, err := http.Get(URL)
	if err != nil {
		fmt.Print(err.Error())
		return 0
	}
	if err != nil {
		log.Fatal(err)
	} else if response.StatusCode == http.StatusCreated {
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			println(err)
		} else {
			data, err := strconv.Atoi(string(responseData))
			if err != nil {
				println(err)
			}
			return data
		}
	}
	return 0
}
func trips(w http.ResponseWriter, r *http.Request) {
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
		println("Can't delete Trip record")
	} else if r.Method == "GET" {
		if params["TripID"] == " " {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Please provide Trip ID"))
			return
		}
		//GET trip using TripID
		TripID, err := strconv.Atoi(params["TripID"])
		if err != nil {
			fmt.Println(err)
		}
		TripInformation := GetTrip(db, TripID)
		if err != nil {
			fmt.Printf("There was an error encoding the json. err = %s", err)
		} else if TripInformation.TripID == 0 { // Check if data is empty
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("There is no trip session"))
			return
		} else {
			json.NewEncoder(w).Encode(GetTrip(db, TripInformation.TripID))
			w.WriteHeader(http.StatusAccepted)
			return
		}
	}
	if r.Header.Get("Content-type") == "application/json" {

		// POST is for creating new course
		if r.Method == "POST" {
			// read the string sent to the service
			var newTripData Data
			reqBody, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newTripData)

				if newTripData.PassengerEmail == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply in JSON format"))
					return
				} else {
					var newTrip tripinfo
					newTrip.PassengerID = CheckPassenger(newTripData.PassengerEmail) //Check if user is in database
					if newTrip.PassengerID == 0 {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("passenger not found"))
						return
					}
					//find available driver
					AvailableDriverID := CheckDriverAvailability(db, GetAllDriver())
					if AvailableDriverID != "" {
						newTrip.DriverID = AvailableDriverID
					} else {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("There is no available driver"))
						return
					}
					//Add trip details
					newTrip.PickUpLocation = newTripData.PickUpLocation
					newTrip.DropOffLocation = newTripData.DropOffLocation
					newTrip.PickUpTime = newTripData.PickUpTime
					newTrip.DropOffTime = newTripData.DropOffTime
					//Check if passenger or driver has uncompleted trips
					if CheckTrip(db, newTrip.PassengerID, newTrip.DriverID) {
						CreateTrip(db, newTrip)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("Trip created successfully"))
						return
					} else {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("The current driver or passenger has an ongoing trip"))
						return
					}
				}
			}
		}
		//---PUT is for creating or updating
		// existing course---
		if r.Method == "PUT" {
			var updatedTrip tripinfo
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				err := json.Unmarshal(reqBody, &updatedTrip)
				if err != nil {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("There was an error encoding the json."))
					return
				}
				if updatedTrip.PassengerID == 0 || updatedTrip.DriverID == "" || updatedTrip.TripID == 0 {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply passenger information"))
					return
				} else {
					if CheckTrip(db, updatedTrip.PassengerID, updatedTrip.DriverID) {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("There is no trips found"))
						return
					} else {
						//To update trip details
						if EditTrip(db, updatedTrip) { //To edit trip
							w.WriteHeader(http.StatusCreated)
							w.Write([]byte("Trip updated successfully"))
							return
						} else {
							w.WriteHeader(http.StatusUnprocessableEntity)
							w.Write([]byte("Trip unable to update"))
							return
						}
					}
				}
			}
		}

	}
}

func GetLicensePlateNumber(DriverID string) string {
	URL := "http://localhost:5002/api/v1/driver/GetLicensePlate/" + DriverID

	response, err := http.Get(URL)
	if err != nil {
		fmt.Print(err.Error())
		return ""
	}

	if err != nil {
		log.Fatal(err)
	} else if response.StatusCode == http.StatusCreated {
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			println(err)
		} else {
			return string(responseData)
		}
	}
	return ""
}

func GetAllTrips(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}
	//Database
	db, err := sql.Open("mysql", "root:12N28c02@tcp(127.0.0.1:3306)/assignment_db") //Connecting to database
	if err != nil {
		fmt.Println(err)
	}
	params := mux.Vars(r)

	if params["Email"] == "" { // Check if data is empty
		w.WriteHeader(
			http.StatusUnprocessableEntity)
		w.Write([]byte("Please provide valid email"))
		return
	} else {
		CustomerID := CheckPassenger(params["Email"])
		if CustomerID != 0 {
			TripData := GetTrips(db, CustomerID)
			var JSONObject []AllTrip
			for _, data := range TripData {
				var TempAllTripData = AllTrip{PickUpLocation: data.PickUpLocation, DropOffLocation: data.DropOffLocation,
					PickUpTime: data.PickUpTime, DropOffTime: data.DropOffTime,
					Status: data.Status, LicensePlate: GetLicensePlateNumber(data.DriverID)}
				JSONObject = append(JSONObject, TempAllTripData)
			}

			json.NewEncoder(w).Encode(JSONObject)
			w.WriteHeader(http.StatusAccepted)
		}

		return
	}
}
func main() {
	//API
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/trips", home)                                                //Test API
	router.HandleFunc("/api/v1/trips/Router/{TripID}", trips).Methods("GET", "PUT", "POST") //API Manipulation
	router.HandleFunc("/api/v1/trips/{Email}", GetAllTrips)
	fmt.Println("Listening at port 5003")
	log.Fatal(http.ListenAndServe(":5003", router))
}
