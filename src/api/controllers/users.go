package controllers

import (
	"SharingBackend/api"
	"SharingBackend/api/services"
	"SharingBackend/base"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type UserController struct {
	Service         *services.UserService
	SchemaValidator *validator.Validate
}

// SignUpUser Sign-up new user godoc
// @Summary      Sign-up new user
// @Description  This method adds a new user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param   	 request  body  api.SignUpRequest true "User sign-up schema"
// @Success      201  {object}  api.UserResponse
// @Failure      400  {object}  api.ErrorResponse
// @Failure      422  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /v1/users [post]
func (controller UserController) SignUpUser(c *gin.Context) {
	base.Logger.Info("Requested creating user")

	var request api.SignUpRequest
	if err := c.BindJSON(&request); err != nil {
		return
	}

	err := controller.SchemaValidator.Struct(request)
	if err != nil {
		c.Error(base.WrapValidationErrors(err))
		return
	}

	user, err := controller.Service.AddUser(&request)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusCreated, user)
}

// GetUser Get user data
// @Summary      Get authorized user data
// @Description  This method returns authorized user data
// @Tags         Users
// @Security     User
// @Accept       json
// @Produce      json
// @Success      200  {object}  api.UserResponse
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /v1/users/me [get]
func (controller UserController) GetUser(c *gin.Context) {
	base.Logger.Info("Requested authorized user data")

	auth, err := GetAuthenticatedUser(c)
	if err != nil {
		return
	}

	user, err := controller.Service.GetUserPublicData(auth.Username)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}
