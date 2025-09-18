package server

import (
	"errors"
	"net/http"

	"github.com/Dorrrke/project1308/internal/domain"
	"github.com/Dorrrke/project1308/internal/domain/user/models"
	"github.com/Dorrrke/project1308/internal/service/userservice"

	userErrors "github.com/Dorrrke/project1308/internal/domain/user/errors"

	"github.com/gin-gonic/gin"
)

func (srv *RentAPI) register(ctx *gin.Context) {
	var user models.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usecase := userservice.NewUserService(srv.db)
	if err := usecase.SaveUser(user); err != nil {
		if errors.Is(err, userErrors.ErrUserAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (srv *RentAPI) login(ctx *gin.Context) {
	var usReq models.UserRequest
	if err := ctx.BindJSON(&usReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usecase := userservice.NewUserService(srv.db)
	user, err := usecase.LoginUser(usReq)
	if err != nil {
		if errors.Is(err, userErrors.ErrInvalidPassword) || errors.Is(err, userErrors.ErrUserNoExists) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Генерация токена
	access, err := srv.jwtSigner.NewAccessToken(user.UID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refresh, err := srv.jwtSigner.NewRefreshToken(user.UID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie("refresh_token", refresh, domain.CookieMaxAge, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"user":   user,
		"access": access,
	})
}

func (srv *RentAPI) profile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
		return
	}

	usecase := userservice.NewUserService(srv.db)
	uid, ok := userID.(string)
	if !ok {
		srv.log.Error().Msg("userID is not a string")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "userID is not a string"})
		return
	}

	user, err := usecase.GetUserByID(uid)
	if err != nil {
		if errors.Is(err, userErrors.ErrUserNoExists) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
