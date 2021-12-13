/*package database

import (
	"rest-go-demo/database"

	_ "github.com/jinzhu/gorm/dialects/mysql" //Required for MySQL dialect
)

func main() {
	config :=
		database.Config{
			ServerName: "localhost:3306",
			User:       "root",
			Password:   "root",
			DB:         "learning",
		}

	connectionString := database.GetConnectionString(config)
	err := database.Connect(connectionString)
	if err != nil {
		panic(err.Error())
	}
}
*/ /*
/*package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Person struct { // map this type to the record in the table
	ID        int
	FirstName string
	LastName  string
	Age       int
}

func DeleteRecord(db *sql.DB, ID int) {
	query := fmt.Sprintf(
		"DELETE FROM Persons WHERE ID='%d'", ID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

func EditRecord(db *sql.DB, ID int, FN string, LN string, Age int) {
	query := fmt.Sprintf(
		"UPDATE Persons SET FirstName='%s', LastName='%s', Age=%d WHERE ID=%d",
		FN, LN, Age, ID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

func InsertRecord(db *sql.DB, ID int, FN string, LN string, Age int) {
	query := fmt.Sprintf("INSERT INTO Persons VALUES (%d, '%s', '%s', %d)",
		ID, FN, LN, Age)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func GetRecords(db *sql.DB) {
	results, err := db.Query("Select * FROM my_db.Persons")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		// map this type to the record in the table
		var person Person
		err = results.Scan(&person.ID, &person.FirstName,
			&person.LastName, &person.Age)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println(person.ID, person.FirstName,
			person.LastName, person.Age)
	}
}

func main() {
	// Use mysql as driverName and a valid DSN as dataSourceName:
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")

	// handle error
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened")
	}
	// InsertRecord(db, 2, "Wallace","Tan", 55)
	//EditRecord(db, 2, "Taylor", "Swift", 23)
	DeleteRecord(db, 2)
	GetRecords(db)
	// defer the close till after the main function has finished executing
	defer db.Close()

}
*/
