CREATE DATABASE driver_management;

USE driver_management;

CREATE TABLE drivers (
  id INT NOT NULL AUTO_INCREMENT,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  license_no VARCHAR(255) NOT NULL,
  license_class VARCHAR(255) NOT NULL,
  expiry_date DATE NOT NULL,
  status VARCHAR(255) NOT NULL,
  PRIMARY KEY (id)
);


CREATE TABLE fuel_consumption (
  id INT NOT NULL AUTO_INCREMENT,
  vehicle_id INT NOT NULL,
  date DATE NOT NULL,
  miles INT NOT NULL,
  gallons FLOAT NOT NULL,
  price_per_gallon FLOAT NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (vehicle_id) REFERENCES vehicles(id)
);
CREATE TABLE assets (
    id INT PRIMARY KEY,
    name VARCHAR(255),
    description VARCHAR(255),
    acquisition_date DATE,
    status VARCHAR(255)
);
