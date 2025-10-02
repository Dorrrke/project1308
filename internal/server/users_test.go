package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dorrrke/project1308/internal/domain/user/errors"
	"github.com/Dorrrke/project1308/internal/domain/user/models"
	"github.com/Dorrrke/project1308/internal/server/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	var srv RentAPI
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.POST("/register", srv.register)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close()

	type want struct {
		body       string
		statusCode int
	}

	type test struct {
		name     string
		userJson string
		user     models.User
		req      string
		method   string
		mockFlag bool
		err      error
		want     want
	}

	tests := []test{
		{
			name:     "successful register",
			userJson: `{"name":"John Doe","email":"TmB0S@example.com","password":"password123","age":30,"phone":"+1234567890"}`,
			user:     models.User{Name: "John Doe", Email: "TmB0S@example.com", Password: "password123", Age: 30, Phone: "+1234567890"},
			req:      "/register",
			method:   http.MethodPost,
			mockFlag: true,
			err:      nil,
			want: want{
				body:       `{"uid":"","name":"John Doe","email":"TmB0S@example.com","password":"password123","age":30,"phone":"+1234567890"}`,
				statusCode: http.StatusOK,
			},
		},
		{
			name:     "bad json",
			userJson: `{name=JsonDoe}`,
			req:      "/register",
			method:   http.MethodPost,
			mockFlag: false,
			err:      fmt.Errorf("flag err"),
			want: want{
				statusCode: http.StatusBadRequest,
				body:       `{"error":"invalid character`,
			},
		},
		{
			name:     "user already exists",
			userJson: `{"name":"John Doe","email":"TmB0S@example.com","password":"password123","age":30,"phone":"+1234567890"}`,
			user:     models.User{Name: "John Doe", Email: "TmB0S@example.com", Password: "password123", Age: 30, Phone: "+1234567890"},
			req:      "/register",
			method:   http.MethodPost,
			mockFlag: true,
			err:      errors.ErrUserAlreadyExists,
			want: want{
				statusCode: http.StatusConflict,
				body:       `{"error":"user with this email or phone already exists"}`,
			},
		},

		{
			name:     "internal error",
			userJson: `{"name":"John Doe","email":"TmB0S@example.com","password":"password123","age":30,"phone":"+1234567890"}`,
			user:     models.User{Name: "John Doe", Email: "TmB0S@example.com", Password: "password123", Age: 30, Phone: "+1234567890"},
			req:      "/register",
			method:   http.MethodPost,
			mockFlag: true,
			err:      fmt.Errorf("internal error"),
			want: want{
				statusCode: http.StatusInternalServerError,
				body:       `{"error":"internal error"}`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewStorage(t)
			srv.db = repo
			if tc.mockFlag {
				repo.On("SaveUser", mock.MatchedBy(func(u models.User) bool {
					if u.Name == tc.user.Name &&
						u.Email == tc.user.Email &&
						u.Age == tc.user.Age &&
						u.Phone == tc.user.Phone {
						return true
					}
					return false
				})).Return(tc.err)
			}

			req := resty.New().R()
			req.Method = tc.method
			req.URL = httpSrv.URL + tc.req
			req.Body = tc.userJson

			resp, err := req.Send()
			assert.NoError(t, err)

			assert.Equal(t, tc.want.statusCode, resp.StatusCode())
			respBody := string(resp.Body())
			if tc.err == nil {
				assert.Equal(t, tc.want.body, respBody)
			} else {
				assert.Contains(t, respBody, tc.want.body)
			}

		})
	}
}

func TestLogin(t *testing.T) {
	var srv RentAPI
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.POST("/login", srv.login)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close()

	testPass := "password123"
	testPassHash, err := bcrypt.GenerateFromPassword([]byte(testPass), bcrypt.DefaultCost)
	assert.NoError(t, err)
	type userResp struct {
		User   models.User `json:"user"`
		Access string      `json:"access"`
	}
	type want struct {
		cookie     bool
		body       userResp
		error      string
		statusCode int
	}

	type test struct {
		name     string
		userJson string
		userReq  models.UserRequest
		req      string
		method   string
		mockFlag bool
		err      error
		want     want
	}

	tests := []test{
		{
			name:     "success",
			userReq:  models.UserRequest{Email: "TmB0S@example.com", Password: testPass},
			userJson: `{"email":"TmB0S@example.com","password":"password123"}`,
			req:      "/login",
			method:   http.MethodPost,
			mockFlag: true,
			err:      nil,
			want: want{
				cookie:     true,
				statusCode: http.StatusOK,
				body: userResp{
					User: models.User{
						Name:     "John Doe",
						Email:    "TmB0S@example.com",
						Password: string(testPassHash),
						Age:      30,
						Phone:    "+1234567890",
					},
				},
			},
		},
		{
			name:     "unauthorized",
			userReq:  models.UserRequest{Email: "TmB0S@example.com", Password: "password1"},
			userJson: `{"email":"TmB0S@example.com","password":"password1"}`,
			req:      "/login",
			method:   http.MethodPost,
			mockFlag: true,
			err:      errors.ErrInvalidPassword,
			want: want{
				cookie:     true,
				statusCode: http.StatusUnauthorized,
				error:      `{"error":"invalid password"}`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := mocks.NewStorage(t)
			srv.db = repo
			if tc.mockFlag {
				repo.On("GetUser", tc.userReq).
					Return(models.User{
						Name:     "John Doe",
						Email:    "TmB0S@example.com",
						Password: string(testPassHash),
						Age:      30,
						Phone:    "+1234567890",
					}, nil)
			}

			req := resty.New().R()
			req.Method = tc.method
			req.URL = httpSrv.URL + tc.req
			req.Body = tc.userJson

			resp, err := req.Send()
			assert.NoError(t, err)

			assert.Equal(t, tc.want.statusCode, resp.StatusCode())
			if tc.err == nil {
				assert.NotNil(t, resp.Cookies())
				var respUser userResp
				err = json.Unmarshal(resp.Body(), &respUser)
				assert.NoError(t, err)
				assert.Equal(t, tc.want.body.User, respUser.User)
				assert.NotEmpty(t, respUser.Access)
			} else {
				assert.Contains(t, string(resp.Body()), tc.want.error)
			}
		})
	}
}
