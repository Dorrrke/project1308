package inmemory

import (
	carDomain "github.com/Dorrrke/project1308/internal/domain/cars/models"
	userDomain "github.com/Dorrrke/project1308/internal/domain/user/models"
)

type Storage struct {
	users map[string]userDomain.User
	cars  map[string]carDomain.Car
}

func NewInMemoryStorage() *Storage {
	return &Storage{
		users: make(map[string]userDomain.User),
		cars:  make(map[string]carDomain.Car),
	}
}
