package carservice

import (
	carErrors "github.com/Dorrrke/project1308/internal/domain/cars/errors"
	carDomain "github.com/Dorrrke/project1308/internal/domain/cars/models"
	"github.com/google/uuid"
)

type CarsStorage interface {
	GetAllCars() ([]carDomain.Car, error)
	GetCarByID(string) (carDomain.Car, error)
	GetAvailableCars() ([]carDomain.Car, error)
	AddCar(carDomain.Car) error
	UpdateAvailable(string) error
}

type CarsService struct {
	db CarsStorage
}

func NewUserService(db CarsStorage) *CarsService {
	return &CarsService{
		db: db,
	}
}

func (s *CarsService) GetAllCars() ([]carDomain.Car, error) {
	return s.db.GetAllCars()
}

func (s *CarsService) GetCarByID(id string) (carDomain.Car, error) {
	car, err := s.db.GetCarByID(id)
	if err != nil {
		return carDomain.Car{}, err
	}

	if !car.Available {
		return carDomain.Car{}, carErrors.ErrCarNotAvailable
	}

	car.Available = false
	if err := s.db.UpdateAvailable(id); err != nil {
		return carDomain.Car{}, err
	}

	return s.db.GetCarByID(id)
}

func (s *CarsService) GetAvailableCars() ([]carDomain.Car, error) {
	return s.db.GetAvailableCars()
}

func (s *CarsService) AddCar(car carDomain.Car) error {
	cid := uuid.New().String()
	car.CID = cid

	return s.db.AddCar(car)
}
