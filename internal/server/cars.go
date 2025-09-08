package server

import (
	"errors"
	"net/http"

	carErrors "github.com/Dorrrke/project1308/internal/domain/cars/errors"
	carDomain "github.com/Dorrrke/project1308/internal/domain/cars/models"
	"github.com/Dorrrke/project1308/internal/service/carservice"
	"github.com/gin-gonic/gin"
)

func (srv *RentAPI) getRent(ctx *gin.Context) {
	carID := ctx.Param("id")
	usecase := carservice.NewUserService(srv.db)
	car, err := usecase.GetCarByID(carID)
	if err != nil {
		if errors.Is(err, carErrors.ErrCarNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, carErrors.ErrCarNotAvailable) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// userID, exists := ctx.Get("userID")
	// if !exists {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
	// 	return
	// }

	ctx.JSON(http.StatusOK, gin.H{"msg": "rent successful", "car": car})
}

func (srv *RentAPI) getAllCars(ctx *gin.Context) {
	usecase := carservice.NewUserService(srv.db)
	cars, err := usecase.GetAllCars()
	if err != nil {
		if errors.Is(err, carErrors.ErrCarsNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"cars": cars})
}

func (srv *RentAPI) getAvailableCars(ctx *gin.Context) {
	usecase := carservice.NewUserService(srv.db)
	cars, err := usecase.GetAvailableCars()
	if err != nil {
		if errors.Is(err, carErrors.ErrCarNotAvailable) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "available to rent", "cars": cars})
}

func (srv *RentAPI) addCar(ctx *gin.Context) {
	var car carDomain.Car
	if err := ctx.BindJSON(&car); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usecase := carservice.NewUserService(srv.db)
	if err := usecase.AddCar(car); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, car)
}
