package asset

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	maxIdleConnections = 5
	maxOpenConnections = 10
)

var (
	ErrAssetNotFound = errors.New("asset not found")
)

type Asset struct {
	ID               int
	Name             string
	SerialNumber     string
	Description      string
	PurchaseDate     time.Time
	PurchaseCost     float64
	DepreciationRate float64
	Depreciation     float64
}

type DBManager struct {
	db *sql.DB
}

func NewDBManager(dsn string) (*DBManager, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	db.SetMaxIdleConns(maxIdleConnections)
	db.SetMaxOpenConns(maxOpenConnections)

	return &DBManager{db: db}, nil
}

func (m *DBManager) Close() error {
	return m.db.Close()
}

func (m *DBManager) GetAssetByID(id int) (*Asset, error) {
	row := m.db.QueryRow("SELECT id, name, description, serial_number, purchase_date, purchase_cost, depreciation_rate, depreciation FROM assets WHERE id = ?", id)

	var asset Asset
	err := row.Scan(&asset.ID, &asset.Name, &asset.Description, &asset.SerialNumber, &asset.PurchaseDate, &asset.PurchaseCost, &asset.DepreciationRate, &asset.Depreciation)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAssetNotFound
		}

		return nil, fmt.Errorf("failed to get asset: %v", err)
	}

	return &asset, nil
}

func (m *DBManager) UpdateAssetDepreciation(id int) error {
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	asset, err := m.GetAssetByID(id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get asset: %v", err)
	}

	depreciation := calculateDepreciation(*asset)

	_, err = tx.Exec("UPDATE assets SET depreciation = ? WHERE id = ?", depreciation, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update asset record: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Calculated depreciation for asset %d: %.2f", id, depreciation)

	return nil
}

func calculateDepreciation(asset Asset) float64 {
	currentYear := time.Now().Year()
	yearsSincePurchase := currentYear - asset.PurchaseDate.Year()

	depreciation := asset.PurchaseCost * asset.DepreciationRate * float64(yearsSincePurchase)

	return depreciation
}

func (m *DBManager) GenerateMaintenanceReport() {
	// Generate maintenance report logic here
}

func (m *DBManager) GenerateDepreciationReport() {
	// Generate depreciation report logic here
}
