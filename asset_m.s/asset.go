package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Vehicle struct {
	ID        int
	Make      string
	Model     string
	Year      int
	VIN       string
	Mileage   int
	LastServ  time.Time
	Purchase  time.Time
	PurchasePrice float64
}

type Asset struct {
	ID               int
	Name             string
	SerialNumber     string
	Description      string
	PurchaseDate     time.Time
	PurchasePrice    float64
	DepreciationRate float64
}

func main() {
	// Open a database connection
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/vehicle_management")
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

	// Get asset information
	assets, err := getAssets(db)
	if err != nil {
		log.Fatalf("Failed to get assets: %v", err)
	}

	for _, asset := range assets {
		// Calculate depreciation
		depreciation := calculateDepreciation(asset)

		// Update asset record with depreciation
		stmt, err := db.Prepare("UPDATE assets SET depreciation = ? WHERE id = ?")
		if err != nil {
			log.Fatalf("Failed to prepare update statement: %v", err)
		}
		_, err = stmt.Exec(depreciation, asset.ID)
		if err != nil {
			log.Fatalf("Failed to update asset record: %v", err)
		}

		fmt.Printf("Calculated depreciation for asset %d: %.2f\n", asset.ID, depreciation)
	}

	// Generate reports
	generateMaintenanceReport(db)
	generateDepreciationReport(db)
}

func getAssets(db *sql.DB) ([]Asset, error) {
    // Retrieve all assets from the database
    rows, err := db.Query("SELECT id, name, description, serial_number, purchase_date, purchase_cost, depreciation_rate FROM assets")
    if err != nil {
        return nil, fmt.Errorf("Failed to retrieve assets: %v", err)
    }
    defer rows.Close()

    assets := []Asset{}
    for rows.Next() {
        asset := Asset{}
        err := rows.Scan(&asset.ID, &asset.Name, &asset.Description, &asset.SerialNumber, &asset.PurchaseDate, &asset.PurchaseCost, &asset.DepreciationRate)
        if err != nil {
            return nil, fmt.Errorf("Failed to scan asset: %v", err)
        }
        assets = append(assets, asset)
    }

    return assets, nil
}

func calculateDepreciation(asset Asset) float64 {
    // Calculate number of years since purchase
    currentYear := time.Now().Year()
    yearsSincePurchase := currentYear - asset.PurchaseDate.Year()

    // Calculate total depreciation
    depreciation := asset.PurchasePrice * asset.DepreciationRate * float64(yearsSincePurchase)

    return depreciation
}


