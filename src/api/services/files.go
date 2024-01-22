package services

import (
	"SharingBackend/api"
	"SharingBackend/base"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

type FilesService struct {
	Context    *context.Context
	Collection *mongo.Collection
}

func (service FilesService) CheckFileDataExists(fileId string) (bool, error) {
	var data api.FileData
	opts := options.FindOne().SetProjection(bson.D{{Key: "identifier", Value: 1}})
	err := service.Collection.FindOne(*service.Context, bson.D{
		primitive.E{Key: "identifier", Value: fileId},
	}, opts).Decode(&data)

	if err == nil {
		return true, nil
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	} else {
		return false, base.NewDatabaseError(err)
	}
}

func (service FilesService) AddFile(request *api.FileData) (*api.AddFileResponse, error) {
	exists, err := service.CheckFileDataExists(request.Identifier)
	if exists && err != nil {
		return nil, err
	} else if exists {
		err := base.ServiceError{
			Summary: fmt.Sprintf("File '%s' already exist", request.Identifier),
			Status:  http.StatusBadRequest,
		}
		return nil, err
	} else {
		if _, err := service.Collection.InsertOne(*service.Context, &request); err != nil {
			return nil, base.NewDatabaseError(err)
		}
		return &api.AddFileResponse{Identifier: request.Identifier}, nil
	}
}

func (service FilesService) GetFile(fileId string) (*api.FileData, error) {
	var fileData api.FileData
	err := service.Collection.FindOne(*service.Context, bson.D{
		primitive.E{Key: "identifier", Value: fileId},
	}).Decode(&fileData)

	if err == nil {
		return &fileData, nil
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, base.ServiceError{
			Summary: fmt.Sprintf("File '%s' not found", fileId),
			Status:  http.StatusNotFound,
		}
	} else {
		return nil, base.NewDatabaseError(err)
	}
}
