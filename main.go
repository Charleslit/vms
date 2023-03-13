package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Vehicle struct {
	ID    int
	Make  string
	Model string
	Year  int
	VIN   string
}

func main() {
	// Open a database connection
	db, err := sql.Open("mysql", "user:password@tcp(host:port)/vehicle_management")
	if err != nil {
		fmt.Printf("Error opening database connection: %v\n", err)
		return
	}
	defer db.Close()

	// Create a new vehicle
	newVehicle := Vehicle{
		Make:  "Toyota",
		Model: "Camry",
		Year:  2022,
		VIN:   "345678901234567",
	}

	// Insert the new vehicle into the database
	stmt, err := db.Prepare("INSERT INTO vehicles (make, model, year, vin) VALUES (?, ?, ?, ?)")
	if err != nil {
		fmt.Printf("Error preparing statement: %v\n", err)
		return
	}
	res, err := stmt.Exec(newVehicle.Make, newVehicle.Model, newVehicle.Year, newVehicle.VIN)
	if err != nil {
		fmt.Printf("Error executing statement: %v\n", err)
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("Error getting last insert ID: %v\n", err)
		return
	}
	fmt.Println("Inserted vehicle with ID:", lastID)

	// Query the database to retrieve the vehicle information
	vehicle := Vehicle{}
	err = db.QueryRow("SELECT * FROM vehicles WHERE id=?", lastID).Scan(&vehicle.ID, &vehicle.Make, &vehicle.Model, &vehicle.Year, &vehicle.VIN)
	if err !=
