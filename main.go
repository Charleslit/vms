package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Charleslit/vms/asset"
	"github.com/Charleslit/vms/incidents"
	"github.com/Charleslit/vms/maintenance"
	"github.com/Charleslit/vms/vehicle"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Open a database connection
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/vehicle_management")
	//db, err := sql.Open("sqlite3", "./vehicle_management.db")
	// for connecting to sqlite3 database in the directory
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Schedule maintenance for all vehicles that haven't been serviced in the last 6 months or 5000 miles
	vehicles, err := maintenance.GetVehiclesForMaintenance(db, time.Now(), 6*time.Month, 5000)
	if err != nil {
		log.Fatalf("Failed to get vehicles for maintenance: %v", err)
	}

	for _, vehicle := range vehicles {
		// Create a maintenance task
		task := maintenance.CreateMaintenanceTask(vehicle, time.Now())

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

	// Create an accident report
	report := &incidents.Accident{
		VehicleID:    1,
		Type:         "Collision",
		Description:  "Vehicle collided with another car at an intersection",
		ReportedTime: time.Now(),
	}

	err = incidents.CreateAccident(db, report)
	if err != nil {
		log.Fatalf("Failed to create accident report: %v", err)
	}

	// Get all assets
	assets, err := asset.GetAll(db)
	if err != nil {
		log.Fatalf("Failed to get assets: %v", err)
	}

	fmt.Println("All assets:")
	for _, a := range assets {
		fmt.Printf("ID: %d, Name: %s, Type: %s, Model: %s\n", a.ID, a.Name, a.Type, a.Model)
	}

	// Get all vehicles
	vehicles, err = vehicle.GetAll(db)
	if err != nil {
		log.Fatalf("Failed to get vehicles: %v", err)
	}

	fmt.Println("All vehicles:")
	for _, v := range vehicles {
		fmt.Printf("ID: %d, Make: %s, Model: %s, Year: %d, VIN: %s, Mileage: %d, LastServiced: %s\n", v.ID, v.Make, v.Model, v.Year, v.VIN, v.Mileage, v.LastServiced.Format(time.RFC3339))
	}
}
