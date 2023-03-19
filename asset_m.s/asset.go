package asset_m.s

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)


type Asset struct {
	ID               int
	Name             string
	SerialNumber     string
	Description      string
	PurchaseDate     time.Time
	PurchasePrice    float64
	DepreciationRate float64
}

func assets() {
	
	// Schedule maintenance for all vehicles that haven't been serviced in the last 6 months or 5000 miles
	


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


