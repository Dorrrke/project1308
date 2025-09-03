package server

import (
	"fmt"
	"net/http"

	"github.com/Dorrrke/project1308/internal"
	userDomain "github.com/Dorrrke/project1308/internal/domain/user/models"

	"github.com/gin-gonic/gin"
)

type UserStorage interface {
	SaveUser(user userDomain.User) error
	GetUser(userReq userDomain.UserRequest) (userDomain.User, error)
}

type RentAPI struct {
	srv *http.Server
	db  UserStorage
}

func NewServer(cfg internal.Config, db UserStorage) *RentAPI {
	httpSrv := http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), // localhost:8080
	}

	api := RentAPI{
		srv: &httpSrv,
		db:  db,
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
	router := gin.Default()

	users := router.Group("/users")
	users.POST("/login", api.login)
	users.POST("/register", api.register)
	users.GET("/profile")
	users.GET("/cars")

	cars := router.Group("/cars")
	cars.GET("/list")
	cars.POST("/get-rent")
	cars.GET("/rent-cars")
	cars.POST("/add-car")

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello world")
	})

	api.srv.Handler = router
}
