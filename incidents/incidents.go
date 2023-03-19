package incidents

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Incident struct {
	ID           int       `db:"id"`
	Type         string    `db:"incident_type"`
	Description  string    `db:"incident_description"`
	ReportedTime time.Time `db:"reported_time"`
}

type Accident struct {
	Incident
	VehicleID int `db:"vehicle_id"`
}

func CreateAccident(db *sql.DB, report *Accident) error {
	// Insert the accident report into the database
	stmt, err := db.Prepare("INSERT INTO accidents (vehicle_id, incident_type, incident_description, reported_time) VALUES (:vehicle_id, :incident_type, :incident_description, :reported_time)")
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	_, err = stmt.Exec(sql.Named("vehicle_id", report.VehicleID), sql.Named("incident_type", report.Type), sql.Named("incident_description", report.Description), sql.Named("reported_time", report.ReportedTime))
	if err != nil {
		return fmt.Errorf("failed to insert accident report into database: %w", err)
	}

	log.Printf("Accident report created: %v\n", report)

	return nil
}
