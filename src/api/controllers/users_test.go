package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"stealthy-backend/api"
	"stealthy-backend/api/services"
	"stealthy-backend/base"
	"stealthy-backend/tests"
	"testing"
)

func getRequestUrl(config *base.BackendConfig, url string) string {
	apiVersionPath := "/v1"
	return config.Server.BasePath + apiVersionPath + url
}

func setupUsersRouter(
	config *base.BackendConfig,
	userService services.BaseUserService,
	authService services.BaseAuthorizationService,
) *gin.Engine {
	schemaValidator := base.CreateValidator()

	authController := AuthorizationController{
		AuthService: authService,
	}
	userController := UserController{
		Service:         userService,
		SchemaValidator: schemaValidator,
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.NoRoute(api.NoRouteHandler)
	router.NoMethod(api.NoMethodHandler)
	router.Use(api.LogsHandler)
	router.Use(api.ErrorHandler)
	router.Use(api.CORSHandler)

	applicationGroup := router.Group(config.Server.BasePath)
	v1 := applicationGroup.Group("/v1")

	usersGroup := v1.Group("/users")
	usersGroup.POST("", userController.SignUpUser)

	withAuthUsersGroup := v1.Group("/users").Use(authController.Authorize)
	withAuthUsersGroup.GET("/me", userController.GetUser)

	return router
}

type UsersApiTestSuite struct {
	suite.Suite
	Config              *base.BackendConfig
	AuthToken           string
	UserResponseFixture *api.UserResponse
	AddUserFixture      *api.SignUpRequest
	UserFixture         *api.User
}

func (s *UsersApiTestSuite) SetupTest() {
	s.Config = &base.BackendConfig{}
	s.Config.SetDefaults()
	s.Config.MongoDB.Database = "sharing-backend-test"
	s.Config.Server.Socket = "test-app-host"
	s.Config.Logs.AppName = "sharing-backend-test"

	s.AuthToken = "authorization_token"
	s.AddUserFixture = &api.SignUpRequest{
		Username: "valid_username",
		Password: "v@l1d_p@ssw0RD",
	}
	s.UserResponseFixture = &api.UserResponse{
		Username: s.AddUserFixture.Username,
	}
	s.UserFixture = &api.User{
		Username: s.AddUserFixture.Username,
	}
}

func (s *UsersApiTestSuite) TestApiGetUser() {
	usersServiceMock := tests.NewBaseUserService(s.T())
	authServiceMock := tests.NewBaseAuthorizationService(s.T())
	authServiceMock.On("ParseToken", s.AuthToken).Return(s.UserFixture, nil)
	usersServiceMock.On("GetUserPublicData", s.UserFixture.Username).Return(
		s.UserResponseFixture, nil,
	)

	router := setupUsersRouter(s.Config, usersServiceMock, authServiceMock)

	recorder := httptest.NewRecorder()
	url := getRequestUrl(s.Config, "/users/me")
	req, err := http.NewRequest("GET", url, nil)

	assert.NoError(s.T(), err)

	headerArray := []string{s.AuthToken}
	req.Header["Authorization"] = headerArray
	router.ServeHTTP(recorder, req)

	assert.Equal(s.T(), http.StatusOK, recorder.Code)
	actualResponse := api.UserResponse{}
	assert.NoError(s.T(), json.Unmarshal(recorder.Body.Bytes(), &actualResponse))
	assert.Equal(s.T(), s.UserResponseFixture, &actualResponse)
}

func (s *UsersApiTestSuite) TestApiAddUser() {
	usersServiceMock := tests.NewBaseUserService(s.T())
	authServiceMock := tests.NewBaseAuthorizationService(s.T())
	usersServiceMock.On("AddUser", s.AddUserFixture).Return(
		s.UserResponseFixture, nil,
	)

	router := setupUsersRouter(s.Config, usersServiceMock, authServiceMock)

	recorder := httptest.NewRecorder()
	url := getRequestUrl(s.Config, "/users")

	requestBody, err := json.Marshal(s.AddUserFixture)
	assert.NoError(s.T(), err)

	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBody))

	assert.NoError(s.T(), err)

	router.ServeHTTP(recorder, req)

	assert.Equal(s.T(), http.StatusCreated, recorder.Code)
	actualResponse := api.UserResponse{}
	assert.NoError(s.T(), json.Unmarshal(recorder.Body.Bytes(), &actualResponse))
	assert.Equal(s.T(), s.UserResponseFixture, &actualResponse)
}

func TestUsersApi(t *testing.T) {
	suite.Run(t, new(UsersApiTestSuite))
}
