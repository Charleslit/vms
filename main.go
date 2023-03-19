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
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/mydb")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Get all vehicles due for mantain
	vehicles, err := mantain.GetVehiclesForMaintenance(db, time.Now(), 6*time.Month, 5000)
	if err != nil {
		log.Fatalf("Failed to get vehicles for mantain: %v", err)
	}

	// Schedule mantain tasks for each vehicle
	for _, v := range vehicles {

		task := mantain.CreateMaintenanceTask(v, time.Now(), "Maintenance", "Regular maintenance check", 6*time.Month)

		if err := mantain.InsertMaintenanceTask(db, task); err != nil {
			log.Fatalf("Failed to insert mantain task into database: %v", err)
		}

		fmt.Printf("Scheduled mantain task for vehicle %d: %s\n", v.ID, task.TaskName)
	}

	// Report an incident
	report := &incidents.Accident{
		VehicleID:    1,
		Type:         "Collision",
		Description:  "Vehicle collided with another car at an intersection",
		ReportedTime: time.Now(),
	}

	if err := incidents.CreateAccident(db, report); err != nil {
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
