package models

import (
	"gooooo/db"
	_ "log"
)

type Car struct {
	ID      int64
	Brand   string `binding:"required"`
	Model   string `binding:"required"`
	Engine  string `binding:"required"`
	Gearbox string `binding:"required"`
}

// Save method to insert the Car instance into the database
func (c *Car) Save() error {
	query := `INSERT INTO cars (brand, model, engine, gearbox) VALUES (?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(c.Brand, c.Model, c.Engine, c.Gearbox)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = id

	return nil
}

// GetAllCars returns all saved car instances from the database
func GetAllCars() ([]Car, error) {
	query := `SELECT id, brand, model, engine, gearbox FROM cars`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []Car
	for rows.Next() {
		var car Car
		if err := rows.Scan(&car.ID, &car.Brand, &car.Model, &car.Engine, &car.Gearbox); err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cars, nil
}
