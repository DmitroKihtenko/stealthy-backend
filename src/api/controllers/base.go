package controllers

import (
	"SharingBackend/api"
	"SharingBackend/base"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func generateShortUUID() string {
	return base64.RawURLEncoding.EncodeToString([]byte(uuid.New().String()))
}

func GetAuthenticatedUser(c *gin.Context) (*api.User, error) {
	value, exists := c.Get("auth")
	if !exists {
		err := base.ServiceError{
			Summary: "Request not authenticated",
			Status:  http.StatusForbidden,
		}
		c.Error(err)
		return nil, err
	}
	return value.(*api.User), nil
}
