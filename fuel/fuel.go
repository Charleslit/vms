package fuel

import (
    "fmt"
    "strconv"
    "time"
)

type FuelConsumption struct {
    VehicleID    int
    Date         time.Time
    Odometer     int
    FuelQuantity float64
}

func main() {
    // Example fuel consumption data entered by user
    input := map[string]string{
        "vehicleID":    "1",
        "date":         "2023-03-13",
        "odometer":     "50000",
        "fuelQuantity": "20.5",
    }

    // Validate and process the input data
    fc, err := processFuelConsumption(input)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Display the processed fuel consumption data
    fmt.Printf("Vehicle ID: %d\n", fc.VehicleID)
    fmt.Printf("Date: %s\n", fc.Date.Format("2006-01-02"))
    fmt.Printf("Odometer: %d\n", fc.Odometer)
    fmt.Printf("Fuel Quantity: %.2f\n", fc.FuelQuantity)
}

func processFuelConsumption(input map[string]string) (*FuelConsumption, error) {
    // Validate the vehicle ID
    vehicleID, err := strconv.Atoi(input["vehicleID"])
    if err != nil {
        return nil, fmt.Errorf("Invalid vehicle ID")
    }

    // Parse the date
    date, err := time.Parse("2006-01-02", input["date"])
    if err != nil {
        return nil, fmt.Errorf("Invalid date")
    }

    // Validate the odometer reading
    odometer, err := strconv.Atoi(input["odometer"])
    if err != nil {
        return nil, fmt.Errorf("Invalid odometer reading")
    }

    // Validate the fuel quantity
    fuelQuantity, err := strconv.ParseFloat(input["fuelQuantity"], 64)
    if err != nil {
        return nil, fmt.Errorf("Invalid fuel quantity")
    }

    // Create and return a FuelConsumption object with the processed data
    return &FuelConsumption{
        VehicleID:    vehicleID,
        Date:         date,
        Odometer:     odometer,
        FuelQuantity: fuelQuantity,
    }, nil
}
