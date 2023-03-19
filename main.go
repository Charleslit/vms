package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)


func main() {
	// Open a database connection
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/vehicle_management")
	//db, err := sql.Open("sqlite3", "./vehicle_management.db")
	// forconnecting to sqlite3 database in the directory 
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Schedule maintenance for all vehicles that haven't been serviced in the last 6 months or 5000 miles
	vehicles, err := getVehiclesForMaintenance(db, time.Now(), 6*time.Month, 5000)
	if err != nil {
		log.Fatalf("Failed to get vehicles for maintenance: %v", err)
	}

	for _, vehicle := range vehicles {
		// Create a maintenance task
		task := createMaintenanceTask(vehicle, time.Now())

		// Insert the maintenance task into the database
		stmt, err := db.Prepare("INSERT INTO maintenance_tasks (vehicle_id, task_name, task_description, due_date) VALUES (?, ?, ?, ?)")
		if err != nil {
			log.Fatalf("Failed to prepare insert statement: %v", err)
		}
		_, err = stmt.Exec(task.VehicleID, task.TaskName, task.TaskDescription, task.DueDate)
		if err != nil {
			log.Fatalf("Failed to insert maintenance task into database: %v", err)
		}

		fmt.Printf("Scheduled maintenance task for vehicle %d: %s\n", vehicle.ID, task.TaskName)
	}
}

