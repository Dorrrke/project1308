package worker

import (
	"context"
	"time"

	"github.com/Dorrrke/project1308/internal/domain"
	carDomain "github.com/Dorrrke/project1308/internal/domain/cars/models"
	userDomain "github.com/Dorrrke/project1308/internal/domain/user/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type UserStorage interface {
	SaveUser(user userDomain.User) error
	GetUser(userReq userDomain.UserRequest) (userDomain.User, error)
	GetUserByID(uid string) (userDomain.User, error)
}
type CarsStorage interface {
	GetAllCars() ([]carDomain.Car, error)
	GetCarByID(string) (carDomain.Car, error)
	GetAvailableCars() ([]carDomain.Car, error)
	AddCar(carDomain.Car) error
	UpdateAvailable(string) error
}

type Storage interface {
	UserStorage
	CarsStorage
}
type Assigner struct {
	stor       *Storage
	ordersChan chan uuid.UUID
	workerCnt  int
	ctx        context.Context
	cancel     context.CancelFunc
	log        zerolog.Logger
}

func NewAssigner(st *Storage, workersCnt int) *Assigner {
	ctx, cancel := context.WithCancel(context.Background())
	return &Assigner{
		stor:       st,
		ordersChan: make(chan uuid.UUID, workersCnt),
		workerCnt:  workersCnt,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (a *Assigner) Start() {
	for i := 0; i < a.workerCnt; i++ {
		go a.workerLoop(i + 1)
	}

	go a.scanner()
}

func (a *Assigner) workerLoop(id int) {
	a.log.Debug().Msgf("worker %d started", id)
	for {
		select {
		case <-a.ctx.Done():
			a.log.Debug().Msgf("worker %d stopped", id)
			return
		case oid := <-a.ordersChan:
			a.handleOne(oid)
		}
	}
}

func (a *Assigner) handleOne(oid uuid.UUID) {
	_, cancel := context.WithTimeout(a.ctx, domain.ContextTimeout)
	defer cancel()

	// TODO: попытка отдать аренду ( проверка доступности авто )
}

func (a *Assigner) Stop() {
	a.cancel()
	close(a.ordersChan)
}

func (a *Assigner) Submit(orderID uuid.UUID) {
	select {
	case a.ordersChan <- orderID:
	default:
		a.log.Info().Msgf("job channel full, dropping order %s", orderID.String())
	}
}

func (a *Assigner) scanner() {
	t := time.NewTicker(5 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-a.ctx.Done():
			return
		case <-t.C:
			// TODO: проверка на наличие заказов

			a.Submit(uuid.New())
		}
	}
}
