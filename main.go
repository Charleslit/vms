package main

import (
	"database/sql"
	"fmt"
	"log"

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
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/vehicle_management")
//db, err := sql.Open("sqlite3", "./vehicle_management.db")
// forconnecting to sqlite3 database in the directory 
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
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
		log.Fatalf("Failed to prepare insert statement: %v", err)
	}
	res, err := stmt.Exec(newVehicle.Make, newVehicle.Model, newVehicle.Year, newVehicle.VIN)
	if err != nil {
		log.Fatalf("Failed to insert vehicle into database: %v", err)
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		log.Fatalf("Failed to retrieve last inserted ID: %v", err)
	}
	fmt.Println("Inserted vehicle with ID:", lastID)

	// Query the database to retrieve the vehicle information
	vehicle := Vehicle{}
	err = db.QueryRow("SELECT * FROM vehicles WHERE id=?", lastID).Scan(&vehicle.ID, &vehicle.Make, &vehicle.Model, &vehicle.Year, &vehicle.VIN)
	if err != nil {
		log.Fatalf("Failed to retrieve vehicle from database: %v", err)
	}
	fmt.Println("Retrieved vehicle:", vehicle)

	// Update the vehicle information
	stmt, err = db.Prepare("UPDATE vehicles SET make=?, model=?, year=?, vin=? WHERE id=?")
	if err != nil {
		log.Fatalf("Failed to prepare update statement: %v", err)
	}
	_, err = stmt.Exec("Toyota", "Corolla", 2022, "23456789012345678", lastID)
	if err != nil {
		log.Fatalf("Failed to update vehicle in database: %v", err)
	}
	fmt.Println("Updated vehicle with ID:", lastID)

	// Delete the vehicle from the database
	stmt, err = db.Prepare("DELETE FROM vehicles WHERE id=?")
	if err != nil {
		log.Fatalf("Failed to prepare delete statement: %v", err)
	}
	_, err = stmt.Exec(lastID)
	if err != nil {
		log.Fatalf("Failed to delete vehicle from database: %v", err)
	}
	fmt.Println("Deleted vehicle with ID:", lastID)
}
