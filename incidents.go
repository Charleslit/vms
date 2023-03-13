package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Incident struct {
	ID           int
	Type         string
	Description  string
	ReportedTime time.Time
}

type Accident struct {
	Incident
	VehicleID int
}

func main() {
	// Open a database connection
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/driver_management")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create an accident report
	report := Accident{
		Incident: Incident{
			Type:         "collision",
			Description:  "Driver rear-ended another vehicle at an intersection",
			ReportedTime: time.Now(),
		},
		VehicleID: 123,
	}

	// Insert the accident report into the database
	stmt, err := db.Prepare("INSERT INTO accidents (vehicle_id, incident_type, incident_description, reported_time) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatalf("Failed to prepare insert statement: %v", err)
	}
	_, err = stmt.Exec(report.VehicleID, report.Type, report.Description, report.ReportedTime)
	if err != nil {
		log.Fatalf("Failed to insert accident report into database: %v", err)
	}

	fmt.Printf("Accident report created: %v\n", report)
}
