package maintenance

import (
	"database/sql"
	"fmt"

	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Asset interface {
	GetID() int
	GetName() string
}

type Vehicle struct {
	ID       int
	Make     string
	Model    string
	Year     int
	VIN      string
	Mileage  int
	LastServ time.Time
}

func (v *Vehicle) GetID() int {
	return v.ID
}

func (v *Vehicle) GetName() string {
	return fmt.Sprintf("%d %s %s", v.Year, v.Make, v.Model)
}

func GetAssetsForMaintenance(db *sql.DB, currentDate time.Time, maxAge time.Duration, maxMileage int) ([]Asset, error) {
	// Select all vehicles that haven't been serviced in the last maxAge or maxMileage, whichever is lower
	rows, err := db.Query("SELECT id, make, model, year, vin, mileage, last_serv FROM vehicles WHERE last_serv < ? OR mileage > ?",
		currentDate.Add(-maxAge), maxMileage)
	if err != nil {
		return nil, fmt.Errorf("failed to select vehicles for maintenance: %v", err)
	}
	defer rows.Close()

	assets := []Asset{}
	for rows.Next() {
		vehicle := Vehicle{}
		err := rows.Scan(&vehicle.ID, &vehicle.Make, &vehicle.Model, &vehicle.Year, &vehicle.VIN, &vehicle.Mileage, &vehicle.LastServ)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vehicle: %v", err)
		}
		assets = append(assets, &vehicle)
	}

	return assets, nil
}

type MaintenanceTask struct {
	ID              int
	AssetID         int
	AssetName       string
	TaskName        string
	TaskDescription string
	DueDate         time.Time
	Frequency       time.Duration
	LastDone        time.Time
}

func (t *MaintenanceTask) IsOverdue() bool {
	return time.Now().Sub(t.LastDone) > t.Frequency
}

func (t *MaintenanceTask) Complete() {
	t.LastDone = time.Now()
}

func (t *MaintenanceTask) ScheduleNext() {
	t.DueDate = t.LastDone.Add(t.Frequency)
}

func CreateMaintenanceTask(asset Asset, dueDate time.Time, taskName, taskDescription string, frequency time.Duration) MaintenanceTask {
	return MaintenanceTask{
		AssetID:         asset.GetID(),
		AssetName:       asset.GetName(),
		TaskName:        taskName,
		TaskDescription: taskDescription,
		DueDate:         dueDate,
		Frequency:       frequency,
		LastDone:        time.Time{},
	}
}
