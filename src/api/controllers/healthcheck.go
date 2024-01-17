package controllers

import (
	"SharingBackend/api"
	"SharingBackend/base"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CheckHealth Service health
// @Summary      Check service health
// @Description  This method returns service health
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200  {object}  api.HealthcheckResponse
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /v1/health [get]
func CheckHealth(c *gin.Context) {
	base.Logger.Info("Requested healthcheck")

	response := api.HealthcheckResponse{
		Status: "ok",
	}

	c.IndentedJSON(http.StatusOK, response)
}
