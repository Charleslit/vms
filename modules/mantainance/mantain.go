package mantainance

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Vehicle struct {
	ID       int
	Make     string
	Model    string
	Year     int
	VIN      string
	Mileage  int
	LastServ time.Time
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

func getVehiclesForMaintenance(db *sql.DB, currentDate time.Time, maxAge time.Duration, maxMileage int) ([]Vehicle, error) {
	// Select all vehicles that haven't been serviced in the last maxAge or maxMileage, whichever is lower
	rows, err := db.Query("SELECT id, make, model, year, vin, mileage, last_serv FROM vehicles WHERE last_serv < ? OR mileage > ?",
		currentDate.Add(-maxAge), maxMileage)
	if err != nil {
		return nil, fmt.Errorf("Failed to select vehicles for maintenance: %v", err)
	}
	defer rows.Close()

	vehicles := []Vehicle{}
	for rows.Next() {
		vehicle := Vehicle{}
		err := rows.Scan(&vehicle.ID, &vehicle.Make, &vehicle.Model, &vehicle.Year, &vehicle.VIN, &vehicle.Mileage, &vehicle.LastServ)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan vehicle: %v", err)
		}
		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}

func createMaintenanceTask(vehicle Vehicle, currentDate time.Time) MaintenanceTask {
	// Create a maintenance task for the vehicle
	task := MaintenanceTask{}
	task.VehicleID = vehicle.ID
	task.TaskName = fmt.Sprintf("Oil Change for %d %s %s", vehicle.Year, vehicle.Make, vehicle.Model)
	task.TaskDescription = "Perform an oil change and filter replacement"
	task.DueDate = currentDate.AddDate(0, 0, 30) // Schedule the task 30 days from now

	return task
}

type MaintenanceTask struct {
	ID          int
	Description string
	Frequency   time.Duration
	LastDone    time.Time
}

func (t *MaintenanceTask) IsOverdue() bool {
	return time.Now().Sub(t.LastDone) > t.Frequency
}

func (t *MaintenanceTask) ScheduleNext() {
	t.LastDone = time.Now()
}

