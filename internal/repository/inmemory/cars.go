package inmemory

import (
	carErrors "github.com/Dorrrke/project1308/internal/domain/cars/errors"
	carDomain "github.com/Dorrrke/project1308/internal/domain/cars/models"
)

func (storage *Storage) GetAllCars() ([]carDomain.Car, error) {
	if len(storage.cars) == 0 {
		return nil, carErrors.ErrCarsNotFound
	}

	var cars []carDomain.Car
	for _, car := range storage.cars {
		cars = append(cars, car)
	}
	return cars, nil
}

func (storage *Storage) GetCarByID(id string) (carDomain.Car, error) {
	car, ok := storage.cars[id]
	if !ok {
		return carDomain.Car{}, carErrors.ErrCarNotFound
	}
	return car, nil
}

func (storage *Storage) UpdateAvailable(id string) error {
	car := storage.cars[id]
	car.Available = !car.Available
	storage.cars[id] = car
	return nil
}

func (storage *Storage) GetAvailableCars() ([]carDomain.Car, error) {
	if len(storage.cars) == 0 {
		return nil, carErrors.ErrCarsNotFound
	}

	var cars []carDomain.Car
	for _, car := range storage.cars {
		if car.Available {
			cars = append(cars, car)
		}
	}

	if len(cars) == 0 {
		return nil, carErrors.ErrNotAvailableCars
	}
	return cars, nil
}

func (storage *Storage) AddCar(car carDomain.Car) error {
	for key, c := range storage.cars {
		if c.Lable == car.Lable && c.Model == car.Model && c.Year == car.Year {
			c.Count = c.Count + car.Count
			storage.cars[key] = c
			return nil
		}
	}
	storage.cars[car.CID] = car
	return nil
}
