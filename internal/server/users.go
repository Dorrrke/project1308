package server

import (
	"net/http"

	"github.com/Dorrrke/project1308/internal/domain/user/models"

	//"github.com/Dorrrke/project1308/internal/service/userservice"

	"github.com/gin-gonic/gin"
)

func (srv *RentAPI) register(ctx *gin.Context) {
	var user models.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// usecase := userservice.NewUserService(srv.db)
	// if err := usecase.SaveUser(user); err != nil {
	// 	if errors.Is(err, userErrors.ErrUserAlreadyExists) {
	// 		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	ctx.JSON(http.StatusOK, user)
}

func (srv *RentAPI) login(ctx *gin.Context) {
	var usReq models.UserRequest
	if err := ctx.BindJSON(&usReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// usecase := userservice.NewUserService(srv.db)
	// user, err := usecase.LoginUser(usReq)
	// if err != nil {
	// 	if errors.Is(err, userErrors.ErrInvalidPassword) || errors.Is(err, userErrors.ErrUserNoExists) {
	// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	ctx.JSON(http.StatusOK, models.User{})
}
