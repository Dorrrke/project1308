package db

import (
	"context"
	"time"

	carDomain "github.com/Dorrrke/project1308/internal/domain/cars/models"
	"github.com/jackc/pgx/v5"
)

type carsStorage struct {
	db *pgx.Conn
}

func (cs *carsStorage) GetAllCars() ([]carDomain.Car, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := cs.db.Query(ctx, "SELECT * FROM cars")
	if err != nil {
		return nil, err
	}

	var cars []carDomain.Car
	for rows.Next() {
		var car carDomain.Car
		if err := rows.Scan(&car.CID, &car.Lable, &car.Model, &car.Year, &car.Available, &car.Count); err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cars, nil
}

func (us *userStorage) GetCarByID(id string) (carDomain.Car, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var car carDomain.Car
	err := us.db.QueryRow(ctx, "SELECT * FROM cars WHERE cid = $1", id).
		Scan(&car.CID, &car.Lable, &car.Model, &car.Year, &car.Available, &car.Count)
	if err != nil {
		return carDomain.Car{}, err
	}

	return car, nil
}

func (cs *carsStorage) GetAvailableCars() ([]carDomain.Car, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := cs.db.Query(ctx, "SELECT * FROM cars WHERE available = true AND count > 0")
	if err != nil {
		return nil, err
	}

	var cars []carDomain.Car
	for rows.Next() {
		var car carDomain.Car
		if err := rows.Scan(&car.CID, &car.Lable, &car.Model, &car.Year, &car.Available, &car.Count); err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cars, nil
}

func (cs *carsStorage) AddCar(car carDomain.Car) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := cs.db.Exec(ctx, "INSERT INTO cars (cid, lable, model, year, available, count) VALUES ($1, $2, $3, $4, $5, $6)",
		car.CID, car.Lable, car.Model, car.Year, car.Available, car.Count)
	if err != nil {
		return err
	}

	return nil
}

// TODO: добавить в аргументы новй статус Available.
func (cs *carsStorage) UpdateAvailable(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := cs.db.Exec(ctx, "UPDATE cars SET available = false WHERE cid = $1", id)
	if err != nil {
		return err
	}

	return nil
}
