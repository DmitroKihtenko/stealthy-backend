package controllers

import (
	"SharingBackend/api/services"
	"SharingBackend/base"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthorizationController struct {
	AuthService *services.AuthorizationService
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
