package driverm.s

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Driver struct {
	ID           int
	FirstName    string
	LastName     string
	LicenseNo    string
	LicenseClass string
	ExpiryDate   time.Time
	Status       string
}

func main() {
	// Open a database connection
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/driver_management")
	//db, err := sql.Open("sqlite3", "./driver_management.db")
	// for connecting to sqlite3 database in the directory 
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Get all active drivers
	drivers, err := getActiveDrivers(db)
	if err != nil {
		log.Fatalf("Failed to get active drivers: %v", err)
	}

	// Check the status of each driver's license
	for _, driver := range drivers {
		if driver.ExpiryDate.Before(time.Now()) {
			driver.Status = "Expired"
		} else if driver.ExpiryDate.Sub(time.Now()) <= 30*24*time.Hour {
			driver.Status = "Expiring soon"
		} else {
			driver.Status = "Active"
		}

		// Update the driver's status in the database
		stmt, err := db.Prepare("UPDATE drivers SET status=? WHERE id=?")
		if err != nil {
			log.Fatalf("Failed to prepare update statement: %v", err)
		}
		_, err = stmt.Exec(driver.Status, driver.ID)
		if err != nil {
			log.Fatalf("Failed to update driver status in database: %v", err)
		}

		fmt.Printf("Driver %s %s (License No: %s) is %s\n", driver.FirstName, driver.LastName, driver.LicenseNo, driver.Status)
	}
}

func getActiveDrivers(db *sql.DB) ([]Driver, error) {
	// Select all drivers with a status of "Active"
	rows, err := db.Query("SELECT id, first_name, last_name, license_no, license_class, expiry_date, status FROM drivers WHERE status='Active'")
	if err != nil {
		return nil, fmt.Errorf("Failed to select active drivers: %v", err)
	}
	defer rows.Close()

	drivers := []Driver{}
	for rows.Next() {
		driver := Driver{}
		err := rows.Scan(&driver.ID, &driver.FirstName, &driver.LastName, &driver.LicenseNo, &driver.LicenseClass, &driver.ExpiryDate, &driver.Status)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan driver: %v", err)
		}
		drivers = append(drivers, driver)
	}

	return drivers, nil
}
