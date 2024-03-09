package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"stealthy-backend/api"
	"stealthy-backend/api/services"
	"stealthy-backend/base"
)

type TokenController struct {
	UserService     services.BaseUserService
	AuthService     services.BaseAuthorizationService
	SchemaValidator *validator.Validate
}

// SignIn Sign-in user
// @Summary      Sign-in user
// @Description  This method authenticates user and returns access token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param   	 request  body  api.SignInRequest true "User sign-in schema"
// @Success      200  {object}  api.TokenResponse
// @Failure      400  {object}  api.ErrorResponse
// @Failure      422  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /v1/login [post]
func (controller TokenController) SignIn(c *gin.Context) {
	base.Logger.Info("Requested JWT")

	var request api.SignInRequest
	if err := c.BindJSON(&request); err != nil {
		return
	}

	err := controller.SchemaValidator.Struct(request)
	if err != nil {
		c.Error(base.WrapValidationErrors(err))
		return
	}

	user, err := controller.UserService.GetUserByCredentials(&request)
	if err != nil {
		c.Error(err)
		return
	}

	token, err := controller.AuthService.GenerateToken(user)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, api.TokenResponse{Token: token})
}
