package api

import (
	"SharingBackend/base"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func NoRouteHandler(c *gin.Context) {
	response := ErrorResponse{
		Summary: "Route not found",
	}
	statusCode := http.StatusNotFound

	c.JSON(statusCode, response)
}

func NoMethodHandler(c *gin.Context) {
	response := ErrorResponse{
		Summary: "Method not allowed",
	}
	statusCode := http.StatusMethodNotAllowed

	c.JSON(statusCode, response)
}

func LogsHandler(c *gin.Context) {
	base.Logger.WithFields(logrus.Fields{
		"path":   c.Request.URL.Path,
		"method": c.Request.Method,
	}).Info("Incoming request")

	c.Next()

	base.Logger.WithFields(logrus.Fields{
		"path":   c.Request.URL.Path,
		"method": c.Request.Method,
		"status": c.Writer.Status(),
	}).Info("Outgoing response")
}

func ErrorHandler(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			base.Logger.WithFields(logrus.Fields{
				"detail": fmt.Sprint(r),
			}).Error("Unknown processing request error")

			response := ErrorResponse{
				Summary: "Unexpected server error",
				Detail:  fmt.Sprint(r),
			}
			statusCode := http.StatusInternalServerError

			c.JSON(statusCode, response)
		}
	}()

	c.Next()

	for _, err := range c.Errors {
		var parsedError base.ServiceError
		ok := errors.As(err.Err, &parsedError)
		var response ErrorResponse
		var statusCode = 0
		if ok {
			base.Logger.WithFields(logrus.Fields{
				"detail": parsedError.Detail,
			}).Error("Error processing request: ", parsedError.Summary)

			response = ErrorResponse{
				Summary: parsedError.Summary,
				Detail:  parsedError.Detail,
			}
			statusCode = parsedError.Status
		} else {
			base.Logger.WithFields(logrus.Fields{
				"detail": err.Error(),
			}).Error("Error processing request")

			response = ErrorResponse{
				Summary: err.Error(),
				Detail:  err.Meta,
			}
			statusCode = http.StatusInternalServerError
		}
		if statusCode == 0 {
			statusCode = http.StatusInternalServerError
		}

		c.JSON(statusCode, &response)
	}
}

func CORSHandler(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set(
		"Access-Control-Allow-Headers",
		"Content-Type, Content-Length, Accept-Encoding, Authorization, "+
			"Content-Disposition")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
