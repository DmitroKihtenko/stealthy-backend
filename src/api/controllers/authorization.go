package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stealthy-backend/api/services"
	"stealthy-backend/base"
)

type AuthorizationController struct {
	AuthService services.BaseAuthorizationService
}

func (controller AuthorizationController) Authorize(context *gin.Context) {
	tokenString := context.GetHeader("Authorization")
	if tokenString == "" {
		context.Error(base.ServiceError{
			Summary: "Authorization token required",
			Status:  http.StatusUnauthorized,
		})
		context.Abort()
		return
	}
	user, err := controller.AuthService.ParseToken(tokenString)
	if err != nil {
		context.Error(err)
		context.Abort()
		return
	}
	context.Set("auth", user)
	context.Next()
}
