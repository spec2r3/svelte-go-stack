package models

import (
	"database/sql"
	"fmt"
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

func GetCarByID(id int64) (*Car, error) {
	query := `SELECT id, brand, model, engine, gearbox FROM cars WHERE id = ?`
	row := db.DB.QueryRow(query, id)

	var car Car
	err := row.Scan(&car.ID, &car.Brand, &car.Model, &car.Engine, &car.Gearbox)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no car found with ID %d", id)
		}
		return nil, fmt.Errorf("failed to scan car: %w", err)
	}

	return &car, nil
}

func GetCarsByBrand(brand string) ([]*Car, error) {
	query := `SELECT id, brand, model, engine, gearbox FROM cars WHERE brand = ?`
	rows, err := db.DB.Query(query, brand)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var cars []*Car
	for rows.Next() {
		var car Car
		if err := rows.Scan(&car.ID, &car.Brand, &car.Model, &car.Engine, &car.Gearbox); err != nil {
			return nil, fmt.Errorf("failed to scan car: %w", err)
		}
		cars = append(cars, &car)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if len(cars) == 0 {
		return nil, fmt.Errorf("no cars found with brand %s", brand)
	}

	return cars, nil
}

func (car *Car) Update() error {
	query := `UPDATE cars 
	SET brand = ?, model = ?, engine = ?, gearbox = ?
	WHERE id = ?`

	// Prepare the statement
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.Exec(car.Brand, car.Model, car.Engine, car.Gearbox, car.ID)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

func DeleteCarById(id int64) error {
	query := `DELETE FROM cars WHERE id = ?`

	// Prepare the statement
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the statement
	result, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	// Check the number of affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no car found with ID %d", id)
	}

	return nil
}
