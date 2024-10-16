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

func GetAllCars(page, pageSize int) ([]Car, int, error) {
	offset := (page - 1) * pageSize

	var carCount int
	countQuery := `SELECT COUNT(*) FROM cars`
	err := db.DB.QueryRow(countQuery).Scan(&carCount)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, brand, model, engine, gearbox FROM cars LIMIT ? OFFSET ?`
	rows, err := db.DB.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var cars []Car
	for rows.Next() {
		var car Car
		if err := rows.Scan(&car.ID, &car.Brand, &car.Model, &car.Engine, &car.Gearbox); err != nil {
			return nil, 0, err
		}
		cars = append(cars, car)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return cars, carCount, nil
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

func GetCarsByBrand(brand string, page, pageSize int) ([]*Car, int, error) {

	offset := (page - 1) * pageSize

	countQuery := `SELECT COUNT(*) FROM cars WHERE brand = ?`
	var carCount int
	err := db.DB.QueryRow(countQuery, brand).Scan(&carCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute count query: %w", err)
	}

	query := `SELECT id, brand, model, engine, gearbox FROM cars WHERE brand = ? LIMIT ? OFFSET ?`
	rows, err := db.DB.Query(query, brand, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var cars []*Car
	for rows.Next() {
		var car Car
		if err := rows.Scan(&car.ID, &car.Brand, &car.Model, &car.Engine, &car.Gearbox); err != nil {
			return nil, 0, fmt.Errorf("failed to scan car: %w", err)
		}
		cars = append(cars, &car)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return cars, carCount, nil
}

func (c *Car) Update() error {
	query := `UPDATE cars 
	SET brand = ?, model = ?, engine = ?, gearbox = ?
	WHERE id = ?`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.Brand, c.Model, c.Engine, c.Gearbox, c.ID)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

func DeleteCarById(id int64) error {
	query := `DELETE FROM cars WHERE id = ?`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no car found with ID %d", id)
	}

	return nil
}

func (c *Car) Force() error {
	query := `INSERT INTO cars (id, brand, model, engine, gearbox) VALUES (?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(c.ID, c.Brand, c.Model, c.Engine, c.Gearbox)
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

func CarCountChecker(brand, model string) (bool, error) {
	query := `SELECT COUNT(*) FROM cars WHERE brand = ? AND model = ?`
	var count int
	err := db.DB.QueryRow(query, brand, model).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}
	return count > 0, nil
}
