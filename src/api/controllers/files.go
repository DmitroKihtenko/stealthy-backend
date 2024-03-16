package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"stealthy-backend/api"
	"stealthy-backend/api/services"
	"stealthy-backend/base"
	"strconv"
	"strings"
	"time"
)

type FilesController struct {
	FilesService         services.BaseFilesService
	FilesMetadataService services.BaseFilesMetadataService
	FilesExpConfig       *base.FilesExpirationConfig
	SchemaValidator      *validator.Validate
}

// UploadFile Upload file
// @Summary      Upload file for user
// @Description  This method uploads a new file to user's space
// @Tags         Files
// @Security     User
// @Accept       multipart/form-data
// @Produce      json
// @Success      200  {object}  api.AddFileResponse
// @Failure      400  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /v1/files [post]
func (controller FilesController) UploadFile(c *gin.Context) {
	base.Logger.Info("Requested file upload")

	fileMetadata := api.FileMetadata{}
	fileData := api.FileData{}
	auth, err := GetAuthenticatedUser(c)
	if err != nil {
		return
	}
	fileMetadata.Username = auth.Username

	fileMetadata.Identifier = generateShortUUID()

	fileSize, err := strconv.ParseInt(c.GetHeader("Content-Length"), 10, 64)
	if err != nil {
		c.Error(base.NewFilesRequestError(err))
		return
	}
	err = c.Request.ParseMultipartForm(math.MaxInt64)
	if err != nil {
		c.Error(base.NewFilesRequestError(err))
		return
	}

	fileMetadata.Size = fileSize
	fileBytes := make([]byte, fileSize)

	fileMultipart, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.Error(base.NewFilesRequestError(err))
		return
	}

	fileMetadata.Name = fileHeader.Filename

	mimetype := fileHeader.Header.Get("Content-Type")
	if mimetype != "" {
		fileMetadata.Mimetype = mimetype
	} else {
		if fileMetadata.Size > 0 {
			fileMetadata.Mimetype = "application/octet-stream"
		} else {
			fileMetadata.Mimetype = "application/x-empty"
		}
	}

	defer func(fileMultipart multipart.File) {
		err := fileMultipart.Close()
		if err != nil {
			base.Logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Warn("Close file error")
		}
	}(fileMultipart)
	_, err = fileMultipart.Read(fileBytes)
	if err != nil {
		c.Error(base.NewFilesRequestError(err))
		return
	}

	fileMetadata.Creation = time.Now().Unix()
	fileMetadata.Expiration = time.Now().Add(
		time.Minute * time.Duration(
			controller.FilesExpConfig.MinutesLifetimeDefault,
		),
	).Unix()

	fileData.Identifier = fileMetadata.Identifier
	fileData.Data = fileBytes

	err = controller.SchemaValidator.Struct(fileMetadata)
	if err != nil {
		c.Error(base.WrapValidationErrors(err))
		return
	}
	err = controller.SchemaValidator.Struct(fileData)
	if err != nil {
		c.Error(base.WrapValidationErrors(err))
		return
	}

	response, err := controller.FilesMetadataService.AddFileMetadata(&fileMetadata)
	if err != nil {
		c.Error(err)
		return
	}
	_, err = controller.FilesService.AddFile(&fileData)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusCreated, &response)
}

// GetFileMetadataList Get files metadata
// @Summary      Get user's files metadata
// @Description  This method returns a files metadata list for specific user
// @Tags         Files
// @Security     User
// @Accept       json
// @Produce      json
// @Param 		 _ 	  query     api.PaginationQueryParameters false "Pagination parameters"
// @Success      200  {object}  api.FileMetadataListResponse
// @Failure      400  {object}  api.ErrorResponse
// @Failure      422  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /v1/files [get]
func (controller FilesController) GetFileMetadataList(c *gin.Context) {
	base.Logger.Info("Requested files metadata list")

	queryParams := api.PaginationQueryParameters{}

	auth, err := GetAuthenticatedUser(c)
	if err != nil {
		return
	}

	skip, err := strconv.ParseInt(
		c.DefaultQuery(base.SkipQueryParam, strconv.FormatInt(0, 10)), 10, 64)
	if err != nil {
		c.Error(base.NewQueryParamError(base.SkipQueryParam, err))
		return
	}

	val := strconv.FormatInt(20, 10)

	limit, err := strconv.ParseInt(
		c.DefaultQuery(base.LimitQueryParam, val), 10, 64)
	if err != nil {
		c.Error(base.NewPathParamError(base.LimitQueryParam, err))
		return
	}

	queryParams.Skip = skip
	queryParams.Limit = limit
	if err = controller.SchemaValidator.Struct(queryParams); err != nil {
		c.Error(base.WrapValidationErrors(err))
		return
	}

	response, err := controller.FilesMetadataService.GetFileMetadataList(
		&queryParams,
		auth.Username,
	)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, &response)
}

// DownloadFile Download file
// @Summary      Download file
// @Description  This method downloads a specific file
// @Tags         Files
// @Accept       json
// @Produce      multipart/form-data
// @Param 		 identifier path string true "File ID" example(YTE1YzhmMjMtYTEwMi00ZmQ0LTk1ZWUtZmM4ZDAyMjc3MmNm)
// @Success      200
// @Failure      400  {object}  api.ErrorResponse
// @Failure      422  {object}  api.ErrorResponse
// @Failure      500  {object}  api.ErrorResponse
// @Router       /v1/files/{identifier} [get]
func (controller FilesController) DownloadFile(c *gin.Context) {
	base.Logger.Info("Requested file download")

	fileId := c.Param(base.FileIdPathParam)
	if fileId == "" {
		c.Error(base.NewPathParamRequiredError(base.FileIdPathParam))
		return
	}

	fileMetadata, err := controller.FilesMetadataService.GetFileMetadata(fileId)
	if err != nil {
		c.Error(err)
		return
	}
	fileData, err := controller.FilesService.GetFile(fileId)
	if err != nil {
		c.Error(err)
		return
	}

	filename := url.QueryEscape(fileMetadata.Name)
	filename = strings.ReplaceAll(filename, "+", "%20")

	c.Header(
		"Content-Disposition",
		"attachment; filename=\""+filename+"\"",
	)
	c.Header("Content-Type", fileMetadata.Mimetype)
	c.Header("Content-Length", strconv.FormatInt(fileMetadata.Size, 10))

	c.Data(http.StatusOK, fileMetadata.Mimetype, fileData.Data)
}
