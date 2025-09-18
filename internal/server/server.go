package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Dorrrke/project1308/internal"
	carDomain "github.com/Dorrrke/project1308/internal/domain/cars/models"
	userDomain "github.com/Dorrrke/project1308/internal/domain/user/models"
	"github.com/Dorrrke/project1308/internal/server/auth"
	"github.com/Dorrrke/project1308/internal/server/middleware"
	"github.com/rs/zerolog"

	"github.com/gin-gonic/gin"
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

type RentAPI struct {
	srv       *http.Server
	db        Storage
	jwtSigner auth.HS256Signer
	log       *zerolog.Logger
}

func NewServer(cfg internal.Config, db Storage, log *zerolog.Logger) *RentAPI {
	sigenr := auth.HS256Signer{
		Secret:     []byte("UltraH@rdSecretKey123"),
		Issuer:     "rent-service",
		Audience:   "rent-client",
		AccessTTL:  15 * time.Minute,
		RefreshTTL: 24 * 7 * time.Hour,
	}

	httpSrv := http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), // localhost:8080
	}

	api := RentAPI{
		srv:       &httpSrv,
		db:        db,
		jwtSigner: sigenr,
		log:       log,
	}

	api.configRouter()

	return &api
}

func (api *RentAPI) Run() error {
	return api.srv.ListenAndServe()
}

func (api *RentAPI) Shutdown() error {
	return nil
}

func (api *RentAPI) configRouter() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(middleware.ZerologMiddleware(api.log))

	users := router.Group("/users")
	users.POST("/login", api.login)
	users.POST("/register", api.register)
	users.GET("/profile", middleware.AuthMiddleware(api.jwtSigner), api.profile)
	// users.GET("/cars")

	cars := router.Group("/cars")
	cars.GET("/list", api.getAllCars)
	cars.POST("/get-rent/:id", middleware.AuthMiddleware(api.jwtSigner), api.getRent)
	cars.GET("/rent-cars", middleware.AuthMiddleware(api.jwtSigner), api.getAvailableCars)
	cars.POST("/add-car", middleware.AuthMiddleware(api.jwtSigner), api.addCar)

	router.POST("/refresh", api.refresh)
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello world")
	})

	api.srv.Handler = router
}

func (api *RentAPI) refresh(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims, err := api.jwtSigner.ParseRefreshToken(refreshToken, auth.ParseOptions{
		ExpectedIssuer:   api.jwtSigner.Issuer,
		ExpectedAudience: api.jwtSigner.Audience,
		AllowedMethods:   []string{"HS256"},
		Leeway:           60 * time.Second,
	})
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	access, err := api.jwtSigner.NewAccessToken(claims.Subject)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newRefresh, err := api.jwtSigner.NewRefreshToken(claims.Subject)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie("refresh_token", newRefresh, 3600*24*7, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"access": access})
}
