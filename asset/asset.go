package asset

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
	PurchaseCost     float64
	DepreciationRate float64
}

func assets() {
	db, err := sql.Open("mysql", "vms.sql")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	assets, err := getAssets(db)
	if err != nil {
		log.Fatalf("Failed to get assets: %v", err)
	}

	updateAssets(db, assets)

	generateMaintenanceReport(db)
	generateDepreciationReport(db)
}

func getAssets(db *sql.DB) ([]Asset, error) {
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

func updateAssets(db *sql.DB, assets []Asset) error {
	for _, asset := range assets {
		depreciation := calculateDepreciation(asset)

		stmt, err := db.Prepare("UPDATE assets SET depreciation = ? WHERE id = ?")
		if err != nil {
			return fmt.Errorf("Failed to prepare update statement: %v", err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(depreciation, asset.ID)
		if err != nil {
			return fmt.Errorf("Failed to update asset record: %v", err)
		}

		fmt.Printf("Calculated depreciation for asset %d: %.2f\n", asset.ID, depreciation)
	}

	return nil
}

func calculateDepreciation(asset Asset) float64 {
	currentYear := time.Now().Year()
	yearsSincePurchase := currentYear - asset.PurchaseDate.Year()

	depreciation := asset.PurchaseCost * asset.DepreciationRate * float64(yearsSincePurchase)

	return depreciation
}

func generateMaintenanceReport(db *sql.DB) {
	// Generate maintenance report logic here
}

func generateDepreciationReport(db *sql.DB) {
	// Generate depreciation report logic here
}
